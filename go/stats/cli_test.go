package main

import (
	"bytes"
	"errors"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

var invalidArguments = []string{
	"-1", "1001", "1.0", "", " 1", "1 ", "+1", "\t1", "1\t", "1\n",
	"10000", "00000", "\u0661", "\uff11", "1a",
}

type cliCase struct {
	name   string
	args   []string
	report string
}

func validCLICases() []cliCase {
	return []cliCase{
		{"empty", nil, "count=0 sum=0 min=none max=none\n"},
		{"mixed", []string{"4", "0", "9", "4"}, "count=4 sum=17 min=0 max=9\n"},
		{"zero", []string{"0"}, "count=1 sum=0 min=0 max=0\n"},
		{"limit", []string{"1000"}, "count=1 sum=1000 min=1000 max=1000\n"},
		{"leading zeros", []string{"0001"}, "count=1 sum=1 min=1 max=1\n"},
		{"padded", []string{"0000", "0010", "0999"}, "count=3 sum=1009 min=0 max=999\n"},
	}
}

func TestParseArg(t *testing.T) {
	for _, tc := range []struct {
		arg  string
		want uint32
	}{
		{"0", 0}, {"1000", 1000}, {"0001", 1}, {"0000", 0},
		{"0010", 10}, {"999", 999}, {"0999", 999},
	} {
		t.Run(tc.arg, func(t *testing.T) {
			got, err := ParseArg(tc.arg)
			if err != nil || got != tc.want {
				t.Fatalf("got (%d, %v), want (%d, nil)", got, err, tc.want)
			}
		})
	}
	for _, arg := range invalidArguments {
		t.Run("invalid/"+arg, func(t *testing.T) {
			got, err := ParseArg(arg)
			if got != 0 || !errors.Is(err, ErrInvalidArgument) {
				t.Fatalf("got (%d, %v), want (0, ErrInvalidArgument)", got, err)
			}
		})
	}
}

func TestFormatReport(t *testing.T) {
	for _, tc := range []struct {
		name  string
		stats Stats
		want  string
	}{
		{"empty", Stats{}, "count=0 sum=0 min=none max=none\n"},
		{"mixed", Stats{4, 17, 0, 9, true, true}, "count=4 sum=17 min=0 max=9\n"},
		{"wide", Stats{2, 8589934590, 4294967295, 4294967295, true, true},
			"count=2 sum=8589934590 min=4294967295 max=4294967295\n"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			got, err := FormatReport(tc.stats)
			if err != nil || got != tc.want {
				t.Fatalf("got (%q, %v), want (%q, nil)", got, err, tc.want)
			}
		})
	}
}

func TestRun(t *testing.T) {
	for _, tc := range validCLICases() {
		t.Run(tc.name, func(t *testing.T) {
			got := commandResult(Run(tc.args))
			if got.status != 0 || got.stdout != tc.report || got.stderr != "" {
				t.Fatalf("unexpected result: %+v; want report %q", got, tc.report)
			}
		})
	}
	for _, arg := range invalidArguments {
		t.Run("invalid/"+arg, func(t *testing.T) {
			report, err := Run([]string{"4", arg, "9"})
			if report != "" || !errors.Is(err, ErrInvalidArgument) {
				t.Fatalf("got (%q, %v), want empty report and ErrInvalidArgument", report, err)
			}
			got := commandResult(report, err)
			if got.status != 2 || got.stdout != "" || !strings.HasPrefix(got.stderr, "error:") {
				t.Fatalf("unexpected rejection: %+v", got)
			}
		})
	}
}

func TestOverflowAdapter(t *testing.T) {
	// Tests real status wiring, not a fabricated oversized CLI invocation.
	got := commandResult("", ErrOverflow)
	if got.status != 1 || got.stdout != "" || !strings.HasPrefix(got.stderr, "error:") {
		t.Fatalf("unexpected overflow result: %+v", got)
	}
}

func TestCLIExecutable(t *testing.T) {
	name := "stats"
	if runtime.GOOS == "windows" {
		name += ".exe"
	}
	binary := filepath.Join(t.TempDir(), name)
	if output, err := exec.Command("go", "build", "-o", binary, ".").CombinedOutput(); err != nil {
		t.Fatalf("build executable: %v\n%s", err, output)
	}
	for _, tc := range validCLICases() {
		t.Run(tc.name, func(t *testing.T) {
			status, stdout, stderr := runExecutable(t, binary, tc.args)
			if status != 0 || stdout != tc.report || stderr != "" {
				t.Fatalf("status=%d stdout=%q stderr=%q; want %q", status, stdout, stderr, tc.report)
			}
		})
	}
	for _, arg := range invalidArguments {
		t.Run("invalid/"+arg, func(t *testing.T) {
			status, stdout, stderr := runExecutable(t, binary, []string{"4", arg, "9"})
			if status != 2 || stdout != "" || !strings.HasPrefix(stderr, "error:") {
				t.Fatalf("status=%d stdout=%q stderr=%q", status, stdout, stderr)
			}
		})
	}
}

func runExecutable(t *testing.T, binary string, arguments []string) (int, string, string) {
	t.Helper()
	command := exec.Command(binary, arguments...)
	var stdout, stderr bytes.Buffer
	command.Stdout = &stdout
	command.Stderr = &stderr
	err := command.Run()
	if err == nil {
		return 0, stdout.String(), stderr.String()
	}
	var exitError *exec.ExitError
	if !errors.As(err, &exitError) {
		t.Fatalf("launch executable: %v", err)
	}
	return exitError.ExitCode(), stdout.String(), stderr.String()
}
