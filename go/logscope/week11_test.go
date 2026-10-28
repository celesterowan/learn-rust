package logscope

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/http/httptest"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

func TestWeek11Handler(t *testing.T) {
	for _, report := range []string{
		string(fixture(t, "sample.expected.txt")),
		string(fixture(t, "empty.expected.txt")),
	} {
		t.Run(report[:7], func(t *testing.T) {
			defer catchTODO(t)
			handler := NewSummaryHandler(report)
			for _, test := range []struct {
				method, target string
				status         int
			}{
				{"GET", "/summary", 200},
				{"GET", "/summary?x=1", 200},
				{"HEAD", "/summary", 405},
				{"POST", "/summary", 405},
				{"PUT", "/summary", 405},
				{"DELETE", "/summary", 405},
				{"OPTIONS", "/summary", 405},
				{"PATCH", "/summary", 405},
				{"GET", "/", 404},
				{"GET", "/summary/", 404},
				{"POST", "/other", 404},
				{"HEAD", "/summary/", 404},
			} {
				t.Run(test.method+test.target, func(t *testing.T) {
					defer catchTODO(t)
					response := httptest.NewRecorder()
					handler.ServeHTTP(response, httptest.NewRequest(test.method, test.target, nil))
					if response.Code != test.status {
						t.Errorf("status=%d; want %d", response.Code, test.status)
					}
					wantBody := ""
					if test.status == 200 {
						wantBody = report
						if response.Header().Get("Content-Type") != "text/plain; charset=utf-8" {
							t.Errorf("Content-Type=%q", response.Header().Get("Content-Type"))
						}
					}
					if response.Body.String() != wantBody {
						t.Errorf("body=%q; want %q", response.Body.String(), wantBody)
					}
					wantAllow := ""
					if test.status == 405 {
						wantAllow = "GET"
					}
					if response.Header().Get("Allow") != wantAllow {
						t.Errorf("Allow=%q; want %q", response.Header().Get("Allow"), wantAllow)
					}
				})
			}
		})
	}
}

func TestWeek11ConcurrentHandler(t *testing.T) {
	defer catchTODO(t)
	report := string(fixture(t, "sample.expected.txt"))
	handler := NewSummaryHandler(report)
	var wait sync.WaitGroup
	for index := 0; index < 24; index++ {
		wait.Add(1)
		go func() {
			defer wait.Done()
			defer catchTODO(t)
			recorder := httptest.NewRecorder()
			handler.ServeHTTP(recorder, httptest.NewRequest("GET", "/summary", nil))
			if recorder.Code != 200 || recorder.Body.String() != report {
				t.Error("concurrent handler changed immutable report")
			}
		}()
	}
	wait.Wait()
}

func TestWeek11HTTPConfiguration(t *testing.T) {
	defer catchTODO(t)
	handler := http.HandlerFunc(func(http.ResponseWriter, *http.Request) {})
	server := newHTTPServer(handler)
	if server == nil || server.Handler == nil {
		t.Fatal("server needs a dedicated handler")
	}
	if server.ReadHeaderTimeout != 5*time.Second {
		t.Errorf("ReadHeaderTimeout=%s; want 5s", server.ReadHeaderTimeout)
	}
}

type observedListener struct {
	net.Listener
	accepting chan struct{}
	once      sync.Once
	closes    atomic.Int32
}

func (listener *observedListener) Accept() (net.Conn, error) {
	listener.once.Do(func() { close(listener.accepting) })
	return listener.Listener.Accept()
}

func (listener *observedListener) Close() error {
	listener.closes.Add(1)
	return listener.Listener.Close()
}

func startHTTP(t *testing.T, handler http.Handler) (context.CancelFunc, *observedListener, <-chan error) {
	t.Helper()
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	observed := &observedListener{Listener: listener, accepting: make(chan struct{})}
	ctx, cancel := context.WithCancel(context.Background())
	done := make(chan error, 1)
	t.Cleanup(func() {
		cancel()
		listener.Close()
	})
	go func() {
		var serveError error
		defer func() {
			if value := recover(); value != nil {
				serveError = fmt.Errorf("serveHTTP panicked: %v", value)
			}
			done <- serveError
		}()
		serveError = serveHTTP(ctx, observed, handler)
	}()
	select {
	case <-observed.accepting:
	case err := <-done:
		t.Fatalf("server returned before accepting: %v", err)
	case <-time.After(10 * time.Second):
		t.Fatal("deadlock guard: server never began accepting")
	}
	return cancel, observed, done
}

func waitHTTP(t *testing.T, done <-chan error, wantSuccess bool) {
	t.Helper()
	select {
	case err := <-done:
		if wantSuccess && err != nil {
			t.Errorf("normal shutdown reported failure: %v", err)
		}
		if !wantSuccess && err == nil {
			t.Error("failed graceful shutdown must report an error")
		}
	case <-time.After(10 * time.Second):
		t.Fatal("deadlock guard: server did not join")
	}
}

func TestWeek11ServeAndCancel(t *testing.T) {
	report := "valid=0 invalid=0\nlevel INFO=0 WARN=0 ERROR=0\n"
	handler := http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {
		io.WriteString(response, report)
	})
	cancel, listener, done := startHTTP(t, handler)
	client := &http.Client{Timeout: 5 * time.Second}
	t.Cleanup(client.CloseIdleConnections)
	response, err := client.Get("http://" + listener.Addr().String() + "/summary")
	if err != nil {
		t.Fatal(err)
	}
	body, readError := io.ReadAll(response.Body)
	closeError := response.Body.Close()
	if readError != nil || closeError != nil || string(body) != report {
		t.Fatalf("HTTP report=%q read=%v close=%v", body, readError, closeError)
	}
	cancel()
	waitHTTP(t, done, true)
	if listener.closes.Load() == 0 {
		t.Error("listener was not closed before serveHTTP returned")
	}
	connection, err := net.DialTimeout("tcp", listener.Addr().String(), time.Second)
	if err == nil {
		connection.Close()
		t.Error("server still accepts connections after shutdown")
	}
}

func TestWeek11GracefulActiveRequest(t *testing.T) {
	entered, release, exited := make(chan struct{}), make(chan struct{}), make(chan struct{})
	var releaseOnce sync.Once
	releaseRequest := func() { releaseOnce.Do(func() { close(release) }) }
	t.Cleanup(releaseRequest)
	handler := http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {
		close(entered)
		<-release
		io.WriteString(response, "complete\n")
		close(exited)
	})
	cancel, listener, done := startHTTP(t, handler)
	client := &http.Client{Timeout: 5 * time.Second}
	t.Cleanup(client.CloseIdleConnections)
	requestDone := make(chan error, 1)
	go func() {
		response, err := client.Get("http://" + listener.Addr().String() + "/summary")
		if err == nil {
			body, readError := io.ReadAll(response.Body)
			response.Body.Close()
			err = readError
			if err == nil && string(body) != "complete\n" {
				err = fmt.Errorf("incomplete body: %q", body)
			}
		}
		requestDone <- err
	}()
	select {
	case <-entered:
	case err := <-requestDone:
		t.Fatalf("request failed before handler: %v", err)
	case <-time.After(10 * time.Second):
		t.Fatal("deadlock guard: request never entered")
	}
	cancel()
	releaseRequest()
	waitHTTP(t, done, true)
	select {
	case <-exited:
	default:
		t.Error("server returned before the active handler finished")
	}
	if err := <-requestDone; err != nil {
		t.Errorf("graceful shutdown interrupted request: %v", err)
	}
}

func TestWeek11ForcedClose(t *testing.T) {
	entered, exited := make(chan struct{}), make(chan struct{})
	handler := http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {
		close(entered)
		<-request.Context().Done()
		close(exited)
	})
	cancel, listener, done := startHTTP(t, handler)
	client := &http.Client{Timeout: 8 * time.Second}
	t.Cleanup(client.CloseIdleConnections)
	requestDone := make(chan struct{})
	go func() {
		defer close(requestDone)
		response, err := client.Get("http://" + listener.Addr().String() + "/summary")
		if err == nil {
			response.Body.Close()
		}
	}()
	select {
	case <-entered:
	case <-requestDone:
		t.Fatal("request ended before entering handler")
	case <-time.After(10 * time.Second):
		t.Fatal("deadlock guard: request never entered")
	}
	cancel()
	// This request finishes only after Shutdown's fresh 2s deadline and Close.
	waitHTTP(t, done, false)
	select {
	case <-exited:
	case <-time.After(5 * time.Second):
		t.Fatal("Close did not cancel the active request")
	}
	<-requestDone
}

type failedListener struct {
	err    error
	closes atomic.Int32
}

func (listener *failedListener) Accept() (net.Conn, error) { return nil, listener.err }
func (listener *failedListener) Close() error              { listener.closes.Add(1); return nil }
func (listener *failedListener) Addr() net.Addr            { return &net.TCPAddr{} }

func TestWeek11ServeFailure(t *testing.T) {
	defer catchTODO(t)
	cause := errors.New("injected accept failure")
	listener := &failedListener{err: cause}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	err := serveHTTP(ctx, listener, http.NotFoundHandler())
	if !errors.Is(err, cause) || listener.closes.Load() == 0 {
		t.Errorf("serve error=%v closes=%d; want cause preserved and listener closed", err, listener.closes.Load())
	}
}

func TestWeek11ServerStartup(t *testing.T) {
	for _, test := range []struct {
		name                                           string
		args                                           []string
		analyzeError, listenError, serveError          error
		wantStatus, wantAnalyze, wantListen, wantServe int
	}{
		{name: "missing_argument", wantStatus: 2},
		{name: "extra_arguments", args: []string{"a", "b"}, wantStatus: 2},
		{name: "analysis_error", args: []string{"file"}, analyzeError: errReadTest, wantStatus: 1, wantAnalyze: 1},
		{name: "listen_error", args: []string{"file"}, listenError: errReadTest, wantStatus: 1, wantAnalyze: 1, wantListen: 1},
		{name: "serve_error", args: []string{"file"}, serveError: errReadTest, wantStatus: 1, wantAnalyze: 1, wantListen: 1, wantServe: 1},
		{name: "success", args: []string{"file"}, wantAnalyze: 1, wantListen: 1, wantServe: 1},
	} {
		t.Run(test.name, func(t *testing.T) {
			defer catchTODO(t)
			analyzed, listened, served := 0, 0, 0
			var stderr bytes.Buffer
			summary := Summary{Valid: 1, Levels: map[Level]uint64{LevelInfo: 1}, Services: map[string]ServiceSummary{"api": {1, 7, 7}}}
			expected := "valid=1 invalid=0\nlevel INFO=1 WARN=0 ERROR=0\nservice api count=1 total_ms=7 max_ms=7\n"
			listener := &failedListener{}
			status := runServer(context.Background(), test.args, &stderr,
				func(path string) (Summary, error) {
					analyzed++
					if path != "file" {
						t.Errorf("analysis path=%q", path)
					}
					return summary, test.analyzeError
				},
				func(network, address string) (net.Listener, error) {
					listened++
					if network != "tcp" || address != "127.0.0.1:8080" {
						t.Errorf("listen(%q,%q); want tcp on loopback port 8080", network, address)
					}
					if test.listenError != nil {
						return nil, test.listenError
					}
					return listener, nil
				},
				func(ctx context.Context, received net.Listener, handler http.Handler) error {
					served++
					defer received.Close()
					if received != listener {
						t.Error("startup did not transfer its listener")
					}
					// Mutating the original map cannot refresh an immutable report.
					summary.Services["api"] = ServiceSummary{Count: 99}
					for repeat := 0; repeat < 2; repeat++ {
						response := httptest.NewRecorder()
						handler.ServeHTTP(response, httptest.NewRequest("GET", "/summary", nil))
						if response.Body.String() != expected {
							t.Errorf("startup report changed: %q", response.Body.String())
						}
					}
					return test.serveError
				})
			if status != test.wantStatus || analyzed != test.wantAnalyze || listened != test.wantListen || served != test.wantServe {
				t.Errorf("status=%d analyze=%d listen=%d serve=%d; want %d,%d,%d,%d", status, analyzed, listened, served, test.wantStatus, test.wantAnalyze, test.wantListen, test.wantServe)
			}
			if status == 0 && stderr.Len() != 0 || status != 0 && !strings.HasPrefix(stderr.String(), "error:") {
				t.Errorf("stderr=%q", stderr.String())
			}
		})
	}
}

func TestWeek11ServerExecutableFailures(t *testing.T) {
	executable := buildCommand(t, "logserve")
	for _, test := range []struct {
		args   []string
		status int
	}{
		{nil, 2},
		{[]string{"a", "b"}, 2},
		{[]string{filepath.Join(t.TempDir(), "missing")}, 1},
	} {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		command := exec.CommandContext(ctx, executable, test.args...)
		var stdout, stderr bytes.Buffer
		command.Stdout, command.Stderr = &stdout, &stderr
		err := command.Run()
		cancel()
		var exitError *exec.ExitError
		if !errors.As(err, &exitError) || exitError.ExitCode() != test.status || stdout.Len() != 0 || !strings.HasPrefix(stderr.String(), "error:") {
			t.Errorf("args=%q error=%v stdout=%q stderr=%q; want status %d and error: diagnostic", test.args, err, stdout.Bytes(), stderr.Bytes(), test.status)
		}
	}
}
