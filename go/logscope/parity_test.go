//go:build parity

package logscope_test

import (
	"bytes"
	"context"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"testing"
	"time"
)

// TODO (week 12): run this supplied acceptance harness after both implementations
// pass their language-specific gates. Diagnose differences in learner code,
// not by weakening the frozen byte and exit-status expectations below.
func TestParity(t *testing.T) {
	temp := t.TempDir()
	suffix := ""
	if runtime.GOOS == "windows" {
		suffix = ".exe"
	}
	rustTarget := filepath.Join(temp, "rust-target")
	rustBinary := filepath.Join(rustTarget, "debug", "logscope"+suffix)
	goBinary := filepath.Join(temp, "logscope"+suffix)
	manifest := filepath.Join("..", "..", "rust", "logscope", "Cargo.toml")
	parityBuild(t, "cargo", "build", "--offline", "--locked", "--manifest-path", manifest,
		"--target-dir", rustTarget)
	parityBuild(t, "go", "build", "-o", goBinary, "./cmd/logscope")
	binaries := []struct {
		name string
		path string
	}{
		{"rust", rustBinary},
		{"go", goBinary},
	}

	type testCase struct {
		name string
		args []string
		want []byte
		code int
		dir  string
	}
	cases := make([]testCase, 0)
	for _, name := range []string{"sample", "empty", "malformed"} {
		input := filepath.Join("..", "..", "fixtures", name+".log")
		expected := filepath.Join("..", "..", "fixtures", name+".expected.txt")
		cases = append(cases, testCase{
			name: name, args: []string{input}, want: parityRead(t, expected),
		})
	}

	writeInput := func(name string, input []byte) string {
		t.Helper()
		path := filepath.Join(temp, name)
		if err := os.WriteFile(path, input, 0o600); err != nil {
			t.Fatalf("write %s: %v", name, err)
		}
		return path
	}
	sample := parityRead(t, filepath.Join("..", "..", "fixtures", "sample.log"))
	sampleReport := parityRead(t, filepath.Join("..", "..", "fixtures", "sample.expected.txt"))
	spaced := writeInput("input with spaces.log", sample)
	crlf := writeInput("crlf.log", []byte("INFO api 10\r\n\r\nERROR api 0\r"))
	invalid := writeInput("invalid-bytes.log", []byte{
		'I', 'N', 'F', 'O', ' ', 'a', 'p', 'i', ' ', '1', '\n',
		0xff, '\n', 0, '\n',
	})
	writeInput("-input.log", []byte("WARN z 2"))
	cases = append(cases,
		testCase{name: "path_with_spaces", args: []string{spaced}, want: sampleReport},
		testCase{
			name: "crlf_and_unterminated_cr", args: []string{crlf},
			want: []byte("valid=2 invalid=0\nlevel INFO=1 WARN=0 ERROR=1\n" +
				"service api count=2 total_ms=10 max_ms=10\n"),
		},
		testCase{
			name: "invalid_bytes", args: []string{invalid},
			want: []byte("valid=1 invalid=2\nlevel INFO=1 WARN=0 ERROR=0\n" +
				"service api count=1 total_ms=1 max_ms=1\n"),
		},
		testCase{
			name: "leading_dash_filename", args: []string{"-input.log"}, dir: temp,
			want: []byte("valid=1 invalid=0\nlevel INFO=0 WARN=1 ERROR=0\n" +
				"service z count=1 total_ms=2 max_ms=2\n"),
		},
		testCase{name: "no_arguments", code: 2},
		testCase{name: "two_arguments", args: []string{spaced, spaced}, code: 2},
		testCase{name: "missing_file", args: []string{filepath.Join(temp, "missing.log")}, code: 1},
	)

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			for _, binary := range binaries {
				t.Run(binary.name, func(t *testing.T) {
					ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
					defer cancel()
					command := exec.CommandContext(ctx, binary.path, tc.args...)
					command.Dir = tc.dir
					var stdout, stderr bytes.Buffer
					command.Stdout = &stdout
					command.Stderr = &stderr
					err := command.Run()
					if ctx.Err() != nil {
						t.Fatalf("command timed out: %v", ctx.Err())
					}
					code := 0
					if err != nil {
						exit, ok := err.(*exec.ExitError)
						if !ok {
							t.Fatalf("start executable: %v", err)
						}
						code = exit.ExitCode()
					}
					if code != tc.code {
						t.Errorf("exit = %d, want %d; stderr = %q", code, tc.code, stderr.Bytes())
					}
					if !bytes.Equal(stdout.Bytes(), tc.want) {
						t.Errorf("stdout = %q, want %q", stdout.Bytes(), tc.want)
					}
					if tc.code == 0 {
						if stderr.Len() != 0 {
							t.Errorf("unexpected stderr: %q", stderr.Bytes())
						}
					} else if !bytes.HasPrefix(stderr.Bytes(), []byte("error:")) {
						t.Errorf("stderr must start with error:, got %q", stderr.Bytes())
					}
				})
			}
		})
	}
}

func parityBuild(t *testing.T, name string, args ...string) {
	t.Helper()
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Minute)
	defer cancel()
	command := exec.CommandContext(ctx, name, args...)
	if output, err := command.CombinedOutput(); err != nil {
		t.Fatalf("build %q: %v\n%s", command.Args, err, output)
	}
}

func parityRead(t *testing.T, path string) []byte {
	t.Helper()
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read %s: %v", path, err)
	}
	return data
}
