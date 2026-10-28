package logscope

import (
	"context"
	"fmt"
	"io"
	"net"
	"net/http"
)

// NewSummaryHandler captures only immutable report text, not Summary maps.
func NewSummaryHandler(report string) http.Handler {
	// TODO G11.1: Check path before method; implement the exact endpoint table.
	panic(ErrUnimplemented)
}

func newHTTPServer(handler http.Handler) *http.Server {
	// TODO G11.2: Use handler and a five-second ReadHeaderTimeout.
	panic(ErrUnimplemented)
}

// serveHTTP owns listener, joins serving, and closes resources before returning.
func serveHTTP(ctx context.Context, listener net.Listener, handler http.Handler) error {
	// TODO G11.3: Serve in one goroutine. On cancellation, Shutdown with a fresh
	// two-second timeout context, fall back to Close, and join. Normalize
	// http.ErrServerClosed to success. Preserve real serving/shutdown errors.
	return ErrUnimplemented
}

// runServer keeps startup testable without occupying the fixed production port.
func runServer(ctx context.Context, arguments []string, stderr io.Writer,
	analyze func(string) (Summary, error),
	listen func(string, string) (net.Listener, error),
	serve func(context.Context, net.Listener, http.Handler) error) int {
	// TODO G11.4: Validate arguments, analyze once, render before listening,
	// bind tcp 127.0.0.1:8080, then serve. Diagnose errors with "error:".
	// A listener returned successfully by listen transfers to serve.
	fmt.Fprintf(stderr, "error: %s\n", ErrUnimplemented)
	return 1
}

// RunServer is supplied production dependency wiring; no global HTTP handlers.
func RunServer(ctx context.Context, arguments []string, stderr io.Writer) int {
	return runServer(ctx, arguments, stderr, AnalyzeFile, net.Listen, serveHTTP)
}
