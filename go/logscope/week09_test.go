package logscope

import (
	"bytes"
	"context"
	"errors"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"runtime"
	"strings"
	"testing"
	"time"
)

var errReadTest = errors.New("injected read failure")
var errCloseTest = errors.New("injected close failure")
var errWriteTest = errors.New("injected write failure")

type failingReader struct {
	data []byte
	err  error
}

func (reader *failingReader) Read(buffer []byte) (int, error) {
	n := copy(buffer, reader.data)
	reader.data = reader.data[n:]
	return n, reader.err
}

type oneByteReader struct {
	data []byte
}

func (reader *oneByteReader) Read(buffer []byte) (int, error) {
	if len(reader.data) == 0 {
		return 0, io.EOF
	}
	if len(buffer) == 0 {
		return 0, nil
	}
	buffer[0] = reader.data[0]
	reader.data = reader.data[1:]
	return 1, nil
}

type trackedReadCloser struct {
	io.Reader
	closeError error
	closes     int
}

func (reader *trackedReadCloser) Close() error {
	reader.closes++
	return reader.closeError
}

func TestWeek09AnalyzeReader(t *testing.T) {
	for _, test := range []struct {
		name   string
		reader io.Reader
		cause  error
	}{
		{"bytes_reader", bytes.NewReader([]byte("INFO api 7\n")), nil},
		{"string_reader", strings.NewReader("INFO api 7\n"), nil},
		{"one_byte_reads", &oneByteReader{[]byte("INFO api 7\n")}, nil},
		{"immediate_failure", &failingReader{nil, errReadTest}, errReadTest},
		{"partial_failure", &failingReader{[]byte("INFO api 7\n"), errReadTest}, errReadTest},
	} {
		t.Run(test.name, func(t *testing.T) {
			defer catchTODO(t)
			summary, err := AnalyzeReader(test.reader)
			if test.cause != nil {
				if !errors.Is(err, test.cause) || err == test.cause {
					t.Errorf("read error = %v; want wrapped cause", err)
				}
				requireZeroSummary(t, summary)
				return
			}
			if err != nil {
				t.Fatal(err)
			}
			if summary.Valid != 1 || summary.Services["api"] != (ServiceSummary{1, 7, 7}) {
				t.Errorf("reader summary = %#v", summary)
			}
		})
	}
}

func TestWeek09ReadCloser(t *testing.T) {
	for _, test := range []struct {
		name                        string
		readError, closeError, want error
	}{
		{"success", nil, nil, nil},
		{"close_failure", nil, errCloseTest, errCloseTest},
		{"read_failure", errReadTest, nil, errReadTest},
		{"both_fail", errReadTest, errCloseTest, errReadTest},
	} {
		t.Run(test.name, func(t *testing.T) {
			defer catchTODO(t)
			var source io.Reader = strings.NewReader("INFO api 7\n")
			if test.readError != nil {
				source = &failingReader{[]byte("INFO api 7\n"), test.readError}
			}
			reader := &trackedReadCloser{Reader: source, closeError: test.closeError}
			summary, err := analyzeReadCloser(reader)
			if reader.closes != 1 {
				t.Errorf("Close called %d times; want exactly one", reader.closes)
			}
			if test.want != nil {
				if !errors.Is(err, test.want) || err == test.want {
					t.Errorf("error = %v; want wrapped %v", err, test.want)
				}
				if test.readError != nil && errors.Is(err, errCloseTest) {
					t.Error("close error replaced/combined with earlier read error")
				}
				requireZeroSummary(t, summary)
			} else if err != nil || summary.Valid != 1 {
				t.Errorf("analysis = %#v, %v", summary, err)
			}
		})
	}
}

func TestWeek09AnalyzeFile(t *testing.T) {
	defer catchTODO(t)
	path := filepath.Join(t.TempDir(), "sample with spaces.log")
	if err := os.WriteFile(path, fixture(t, "sample.log"), 0600); err != nil {
		t.Fatal(err)
	}
	got, err := AnalyzeFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if Render(got) != string(fixture(t, "sample.expected.txt")) {
		t.Error("AnalyzeFile changed the fixture report")
	}
	// On Windows a leaked open file handle also prevents this removal.
	if err := os.Remove(path); err != nil {
		t.Fatalf("file remains open after analysis: %v", err)
	}
	missing, err := AnalyzeFile(path)
	if !errors.Is(err, os.ErrNotExist) {
		t.Errorf("missing file error = %v; want wrapped os.ErrNotExist", err)
	}
	requireZeroSummary(t, missing)
}

type failedWriter struct {
	short  bool
	writes int
}

func (writer *failedWriter) Write(data []byte) (int, error) {
	writer.writes++
	if writer.short {
		return len(data) / 2, nil
	}
	return 0, errWriteTest
}

func TestWeek08CLIAdapter(t *testing.T) {
	for _, test := range []struct {
		name      string
		arguments []string
		status    int
	}{
		{"no_arguments", nil, 2},
		{"too_many_arguments", []string{"a", "b"}, 2},
		{"missing_file", []string{filepath.Join(t.TempDir(), "missing")}, 1},
	} {
		t.Run(test.name, func(t *testing.T) {
			defer catchTODO(t)
			var stdout, stderr bytes.Buffer
			status := RunCLI(test.arguments, &stdout, &stderr)
			if status != test.status || stdout.Len() != 0 || !strings.HasPrefix(stderr.String(), "error:") {
				t.Errorf("status=%d stdout=%q stderr=%q", status, stdout.String(), stderr.String())
			}
		})
	}
	for _, short := range []bool{false, true} {
		t.Run(map[bool]string{false: "writer_failure", true: "short_write"}[short], func(t *testing.T) {
			defer catchTODO(t)
			writer := &failedWriter{short: short}
			var stderr bytes.Buffer
			status := RunCLI([]string{filepath.Join("..", "..", "fixtures", "sample.log")}, writer, &stderr)
			if status != 1 || writer.writes == 0 || !strings.HasPrefix(stderr.String(), "error:") {
				t.Errorf("status=%d writes=%d stderr=%q", status, writer.writes, stderr.String())
			}
		})
	}
	t.Run("broken_stderr_is_best_effort", func(t *testing.T) {
		defer catchTODO(t)
		var stdout bytes.Buffer
		if status := RunCLI(nil, &stdout, &failedWriter{}); status != 2 || stdout.Len() != 0 {
			t.Errorf("broken stderr changed argument status: %d", status)
		}
	})
}

func buildCommand(t *testing.T, name string) string {
	t.Helper()
	filename := name
	if runtime.GOOS == "windows" {
		filename += ".exe"
	}
	path := filepath.Join(t.TempDir(), filename)
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()
	command := exec.CommandContext(ctx, "go", "build", "-o", path, "./cmd/"+name)
	if output, err := command.CombinedOutput(); err != nil {
		t.Fatalf("executable build failed (not an intentional TODO test failure): %v\n%s", err, output)
	}
	return path
}

func TestWeek09CLIExecutable(t *testing.T) {
	executable := buildCommand(t, "logscope")
	directory := t.TempDir()
	path := filepath.Join(directory, "path with spaces.log")
	if err := os.WriteFile(path, fixture(t, "sample.log"), 0600); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(directory, "-sample.log"), fixture(t, "sample.log"), 0600); err != nil {
		t.Fatal(err)
	}
	for _, test := range []struct {
		name   string
		args   []string
		status int
		want   []byte
	}{
		{"spaces", []string{path}, 0, fixture(t, "sample.expected.txt")},
		{"hyphen_is_path", []string{"-sample.log"}, 0, fixture(t, "sample.expected.txt")},
		{"no_arguments", nil, 2, nil},
		{"multiple_arguments", []string{path, path}, 2, nil},
		{"missing", []string{filepath.Join(directory, "missing")}, 1, nil},
	} {
		t.Run(test.name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()
			command := exec.CommandContext(ctx, executable, test.args...)
			command.Dir = directory
			var stdout, stderr bytes.Buffer
			command.Stdout, command.Stderr = &stdout, &stderr
			err := command.Run()
			status := 0
			if err != nil {
				var exitError *exec.ExitError
				if !errors.As(err, &exitError) {
					t.Fatal(err)
				}
				status = exitError.ExitCode()
			}
			if status != test.status || !bytes.Equal(stdout.Bytes(), test.want) {
				t.Errorf("status=%d stdout=%q; want status=%d stdout=%q", status, stdout.Bytes(), test.status, test.want)
			}
			if test.status == 0 && stderr.Len() != 0 {
				t.Errorf("success stderr = %q", stderr.Bytes())
			}
			if test.status != 0 && !strings.HasPrefix(stderr.String(), "error:") {
				t.Errorf("failure stderr = %q; want error: prefix", stderr.Bytes())
			}
		})
	}
}

func TestWeek09ReaderFixtureParity(t *testing.T) {
	for _, name := range []string{"sample", "empty", "malformed"} {
		t.Run(name, func(t *testing.T) {
			defer catchTODO(t)
			input := fixture(t, name+".log")
			fromBytes, err := Analyze(input)
			if err != nil {
				t.Fatal(err)
			}
			fromReader, err := AnalyzeReader(bytes.NewReader(input))
			if err != nil || !reflect.DeepEqual(fromBytes, fromReader) {
				t.Errorf("reader=%#v, %v; bytes=%#v", fromReader, err, fromBytes)
			}
		})
	}
}
