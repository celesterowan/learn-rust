package logscope

import "io"

func AnalyzeReader(reader io.Reader) (Summary, error) {
	// TODO G9.1: Read all bytes before Analyze. Wrap read errors with %w;
	// bytes returned alongside an error must not become a partial success.
	return Summary{}, ErrUnimplemented
}

// analyzeReadCloser takes ownership of an already-open handle and closes once.
func analyzeReadCloser(reader io.ReadCloser) (Summary, error) {
	// TODO G9.2: Defer Close. A read/analysis error wins over a close error;
	// otherwise return a wrapped close error and Summary{} if Close fails.
	return Summary{}, ErrUnimplemented
}

func AnalyzeFile(path string) (Summary, error) {
	// TODO G9.3: Open, wrap open errors with %w, transfer the successful handle
	// to analyzeReadCloser, and wire analyzeCLIFile to this function.
	return Summary{}, ErrUnimplemented
}
