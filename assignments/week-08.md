# Week 8: Port the analyzer contract, not Rust syntax

**Prerequisite:** week 7 and the Rust analyzer. **Objectives:** explicit errors,
byte parsing, deterministic maps, and library/command separation.

**Reading (45m):** [Go module layout](https://go.dev/doc/modules/layout)
(basic package and command examples),
[error handling](https://go.dev/blog/error-handling-and-go)
(error values and propagation), and
[os.ReadFile](https://pkg.go.dev/os#ReadFile).

**Six-hour budget:** reading 45m + port design 30m + implementation 150m +
tests 90m + reflection/checks 45m = 360m.

## Assignment

Use the supplied independent `go/logscope` module; do not run `go mod init`
over it. Follow G8.1-G8.7 in the
[Go logscope exercise spec](../docs/exercises/go-logscope.md) for exact types,
function-to-test steps, byte contracts, and build/test-compilation commands.
The code compiles before implementation; active acceptance tests deliberately
fail on learner TODOs. Later weeks' TODOs are excluded by this week's selector.

Use package logscope at the module root, and package main in
`cmd/logscope/main.go`. Import your library by its module path, not a
relative import. Required library APIs:

- `ParseLine(line []byte) (Record, error)` for normalized, nonempty lines.
- `Analyze(input []byte) (Summary, error)`
- `Render(summary Summary) string`

Use the supplied named Level type/constants and Summary types; Record owns its
service string. Summary uses uint64 counters and sums. A failed parse returns
a typed error; Analyze counts that line invalid and continues.
Overflow is an analysis error,
not a malformed line. Check before incrementing or adding.

Implement the entire [shared contract](../docs/logscope-contract.md).
Whole-file os.ReadFile is acceptable this week; week 9 adds an explicit reader
boundary. Do not use bufio.Scanner's default token limit without accounting
for it; simply splitting bytes avoids a hidden long-line truncation policy.
No Unicode Fields or TrimSpace.

After implementing, run from the course root.

**Windows PowerShell**

```powershell
Push-Location .\go\logscope
go run .\cmd\logscope ..\..\fixtures\sample.log
go test -run '^TestWeek08' ./...
gofmt -l .
go vet ./...
Pop-Location
```

**macOS Terminal**

```sh
(
  cd go/logscope &&
  go run ./cmd/logscope ../../fixtures/sample.log &&
  go test -run '^TestWeek08' ./... &&
  gofmt -l . &&
  go vet ./...
)
```

## Acceptance and submission

Run the supplied tests for all record grammar cases and fixture byte outputs.
Inspect and extend the cases with your own examples. The active tests include
byte inputs for CRLF, zero bytes, missing final LF, invalid UTF-8, and reversed
service insertion order. Unit-test a near-overflow update helper.
Manually build and invoke the CLI to check exits using the
[Windows/macOS commands](../SETUP.md#executables-exit-status-and-fixtures);
week 9 automates this. Do not use go run to assert the program's own exit code.
Keep os.Exit at the outer main boundary, after work returns.

Submit module and `docs/submissions/week-08.md`.
**Rubric:** grammar/aggregation 35 + CLI/render parity 30 +
test coverage 25 + explanation/formatting 10 = 100.

**C++ questions:** (1) Why not range over runes for this input?
(2) How do explicit error returns change the control flow compared with
exceptions and Rust's `?`?

**Optional stretch (ungraded):** add error context with `%w` and demonstrate
that the original error remains discoverable.
