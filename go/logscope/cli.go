package logscope

import (
	"fmt"
	"io"
	"os"
)

// analyzeCLIFile is supplied week-8 wiring. Week 9 replaces its body.
func analyzeCLIFile(path string) (Summary, error) {
	// TODO G9.3: Replace this whole-file adapter with a call to AnalyzeFile.
	input, err := os.ReadFile(path)
	if err != nil {
		return Summary{}, err
	}
	return Analyze(input)
}

// RunCLI is supplied process infrastructure. Arguments exclude argv[0].
func RunCLI(arguments []string, stdout, stderr io.Writer) int {
	if len(arguments) != 1 {
		fmt.Fprintln(stderr, "error: expected exactly one file path")
		return 2
	}
	summary, err := analyzeCLIFile(arguments[0])
	if err != nil {
		fmt.Fprintf(stderr, "error: %s\n", err)
		return 1
	}
	report := Render(summary)
	written, err := io.WriteString(stdout, report)
	if err == nil && written != len(report) {
		err = io.ErrShortWrite
	}
	if err != nil {
		fmt.Fprintf(stderr, "error: %s\n", err)
		return 1
	}
	return 0
}
