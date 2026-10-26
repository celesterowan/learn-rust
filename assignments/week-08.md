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

```powershell
New-Item -ItemType Directory -Force .\go\logscope\cmd\logscope
Push-Location .\go\logscope
go mod init example.com/rust-go-course/logscope
Pop-Location
```

Use package logscope at the module root, and package main in
`cmd/logscope/main.go`. Import your library by its module path, not a
relative import. Required library APIs:

- `ParseLine(line []byte) (Record, error)` for normalized, nonempty lines.
- `Analyze(input []byte) (Summary, error)`
- `Render(summary Summary) string`

Use a named Level type and constants; Record owns its service string.
Design Summary with uint64 counters and sums. A failed parse returns an error;
Analyze counts that line invalid and continues. Overflow is an analysis error,
not a malformed line. Check before incrementing or adding.

Implement the entire [shared contract](../docs/logscope-contract.md).
Whole-file os.ReadFile is acceptable this week; week 9 adds an explicit reader
boundary. Do not use bufio.Scanner's default token limit without accounting
for it; simply splitting bytes avoids a hidden long-line truncation policy.
No Unicode Fields or TrimSpace.

```powershell
Push-Location .\go\logscope
go run .\cmd\logscope ..\..\fixtures\sample.log
go test ./...
gofmt -l .
go vet ./...
Pop-Location
```

## Acceptance and submission

Test all record grammar cases and all three fixture byte outputs. Add byte
inputs for CRLF, zero bytes, missing final LF, invalid UTF-8, and reversed
service insertion order. Unit-test a near-overflow update helper.
Manually build and invoke the CLI to check exits; week 9 automates this.
Keep os.Exit at the outer main boundary, after work returns.

Submit module and `docs/submissions/week-08.md`.
**Rubric:** grammar/aggregation 35 + CLI/render parity 30 +
test coverage 25 + explanation/formatting 10 = 100.

**C++ questions:** (1) Why not range over runes for this input?
(2) How do explicit error returns change the control flow compared with
exceptions and Rust's `?`?

**Optional stretch (ungraded):** add error context with `%w` and demonstrate
that the original error remains discoverable.
