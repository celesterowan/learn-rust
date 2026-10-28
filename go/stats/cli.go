package main

import (
	"errors"
	"fmt"
)

// ParseArg accepts 1-4 ASCII digits with value 0..1000, including leading zeros.
// All other input returns zero and ErrInvalidArgument.
func ParseArg(argument string) (uint32, error) {
	// TODO G3: Validate bytes/length before converting and range-checking.
	return 0, ErrUnimplemented
}

// FormatReport renders a valid Stats as exactly one LF-terminated report.
func FormatReport(stats Stats) (string, error) {
	// TODO G4: Render absent extrema as "none", preserving field order and LF.
	return "", ErrUnimplemented
}

// Run validates ALL arguments before calculating/rendering; arguments omit argv[0].
// It returns an empty report on error and never prints or panics.
func Run(arguments []string) (string, error) {
	// TODO G5: Compose ParseArg, Calculate, and FormatReport; stop on errors.
	return "", ErrUnimplemented
}

// commandOutput and commandResult are supplied process wiring, not learner tasks.
type commandOutput struct {
	status int
	stdout string
	stderr string
}

func commandResult(report string, err error) commandOutput {
	if err == nil {
		return commandOutput{stdout: report}
	}
	status := 1
	if errors.Is(err, ErrInvalidArgument) {
		status = 2
	}
	return commandOutput{status: status, stderr: fmt.Sprintf("error: %s\n", err)}
}
