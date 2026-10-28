# Go logscope: weeks 8-11

This is a buildable learner scaffold, **not an implementation**. Start week 8
after week 7 and the Rust analyzer. Follow the reading budgets, reflections,
submissions, and rubrics in [week 8](../../assignments/week-08.md),
[week 9](../../assignments/week-09.md), [week 10](../../assignments/week-10.md),
and [week 11](../../assignments/week-11.md). The steps below divide those
assignments; they do not add another course or require third-party packages.

## Supplied infrastructure and learner work

`go/logscope/go.mod` is an independent Go 1.22 module:
`example.com/rust-go-course/logscope`. Its root package is `logscope`.
Both `cmd/logscope` and `cmd/logserve` use package `main`. There is no root
module or workspace. **Do not run `go mod init` over this project.**
No `go.sum` or downloaded dependencies are needed.

| File | Ownership and purpose |
| --- | --- |
| `record.go` | Supplied record/error types; learner byte parser |
| `analyze.go` | Supplied aggregate types; learner arithmetic, normalization, aggregation, rendering |
| `cli.go` | Supplied sequential process adapter; one week-9 wiring replacement |
| `reader.go` | Learner reader, close-policy, and file functions |
| `parallel.go` | Learner checked merge and bounded pool; supplied public callback adapter |
| `server.go` | Learner handler, server configuration/lifetime, and startup; supplied production wiring |
| `cmd/logscope/main.go` | Supplied `os.Exit` boundary |
| `cmd/logserve/main.go` | Supplied interrupt context and `os.Exit` boundary |
| `test_helpers_test.go`, `week08_test.go` through `week11_test.go` | Supplied active acceptance tests and test-only resource fakes |

Keep signatures/types stable. Implement the marked TODO bodies; add stdlib
imports as needed. Do not copy Rust syntax or introduce a filesystem framework,
CLI flags, dependencies, extra commands, or a root module. Preserve relevant
comments. Trace the supplied process and signal adapters in your reflection.
Learner algorithms, concurrency coordination, and HTTP lifecycle remain TODOs.

Every error-returning learner stub returns `ErrUnimplemented`; `Render`,
`normalizedLines`, `NewSummaryHandler`, and `newHTTPServer` explicitly panic
because their signatures have no error return. The integer-returning startup
stub prints `error: TODO: ...` and returns 1. Those are informative placeholders,
not legitimate completed behavior or fake empty successes. Test-only panic
guards turn synchronous TODO panics into failures so other tests still run.
They do not skip tests or make production panic acceptable.

**Build and test compilation must pass before implementation.** Acceptance
tests deliberately fail until their steps are implemented. A missing compiler,
syntax error, failed binary build, timeout, or race-tool failure is not an
intentional TODO acceptance failure. Do not weaken assertions or disable tests.

## Frozen bytes, types, and errors

The [shared logscope contract](../logscope-contract.md) is authoritative.
Read its entire required case matrix and the actual files under `fixtures/`.
Expected files are byte oracles, not text to normalize through the shell.

- Split raw bytes on LF. Strip at most one final CR from each resulting line,
  including the unterminated final line. Omit only lines now of length zero.
- For each remaining line, require exactly `LEVEL service duration`, two
  single ASCII spaces, no other whitespace or fields. Validate ASCII before
  converting to strings. Never use Unicode `Fields` or `TrimSpace`.
- Levels are exactly `INFO`, `WARN`, `ERROR`. Service is
  `[a-z][a-z0-9-]{0,15}`. Duration is 1-5 ASCII decimal digits, numeric value
  0..60000 inclusive. Leading zeros are allowed. Non-ASCII, invalid UTF-8,
  NUL, extra CR, and whitespace-only records are malformed, not fatal I/O.
- `Level` is a named string type; constants are `LevelInfo`, `LevelWarn`,
  `LevelError`. `Record` has `Level Level`, `Service string`, and
  `DurationMS uint64`. Its service must not alias mutable input storage.
- `ServiceSummary` has `Count`, `TotalMS`, `MaxMS`, all `uint64`.
  `Summary` has `Valid`, `Invalid` (`uint64`), `Levels map[Level]uint64`,
  and `Services map[string]ServiceSummary`. A zero value is empty; nil maps
  and empty maps mean the same thing. Missing level entries mean zero.
- `ParseLine` returns `Record{}` and `*ParseError{Reason: ...}` for malformed
  input. `Reason` is informative, not exact graded wording. Supplied `Error`
  and `Unwrap` make both `errors.As` and `errors.Is(err, ErrInvalidRecord)`
  work. The empty-line direct parser call is invalid too; only normalization
  omits empty records.
- `ErrOverflow` means a count or sum cannot fit `uint64`. `ErrWorkers`
  means worker count is outside 1..4. Preserve causes with `%w` at I/O
  boundaries. Ordinary malformed input must never panic.
- `Analyze`, reader/file functions, and parallel functions return exactly
  `Summary{}` on error, not a partial aggregate. Helpers updating a summary
  may leave partial state on error; callers must discard it. No rollback
  implementation is required.
- Check before **every** count increment/addition and duration sum addition.
  Maximum duration is selected, not added. Overflow is never counted invalid.
- Parsing, normalization, and analysis do not mutate caller bytes. Merge does
  not mutate or retain the source's maps. No map is shared between workers.

`Render` accepts a valid aggregate (inconsistent hand-built aggregates are
outside its precondition). It always includes all three levels, then services
sorted by ascending ASCII name; it does not mutate the aggregate. Use decimal
numbers, exactly the fields below, LF separators, and exactly one final LF:

```text
valid=6 invalid=3
level INFO=3 WARN=1 ERROR=2
service api count=3 total_ms=45 max_ms=20
service auth count=1 total_ms=7 max_ms=7
service worker count=2 total_ms=100 max_ms=100
```

Those are the complete `sample.log` report bytes. `empty.log` has no records,
and `malformed.log` has invalid=3 and all other counts zero; compare their
corresponding `.expected.txt` files exactly. An empty aggregate has only:

```text
valid=0 invalid=0
level INFO=0 WARN=0 ERROR=0
```

## Ordered implementation map

Commands in this table run **inside `go/logscope`**, not at the course root.
The test selector, not the filename, determines the week: the supplied
`TestWeek08CLIAdapter` lives beside process/resource fakes in `week09_test.go`
but has no dependency on week-9 TODOs.

| Step | File and exact function | Work | Focused acceptance command |
| --- | --- | --- | --- |
| G8.1 | `record.go`: `ParseLine([]byte) (Record, error)` | Exact grammar, typed error, owned service | `go test -run '^TestWeek08ParseLine$' .` |
| G8.2 | `analyze.go`: `checkedAdd(uint64, uint64) (uint64, error)` | Exact sum or zero/ErrOverflow | `go test -run '^TestWeek08CheckedAdd$' .` |
| G8.3 | `normalizedLines([]byte) [][]byte` | Byte splitting, one CR removal, empty omission | `go test -run '^TestWeek08NormalizedLines$' .` |
| G8.4 | `addRecord(*Summary, Record) error` | Checked valid/level/service updates; initialize maps | `go test -run '^TestWeek08AddRecord$' .` |
| G8.5 | `accumulateLine(*Summary, []byte) error` | Count only malformed errors invalid; check increment | `go test -run '^TestWeek08(AccumulateLine|Overflow)$' .` |
| G8.6 | `Analyze([]byte) (Summary, error)` | Compose normalization and accumulation | `go test -run '^TestWeek08Analyze$' .` |
| G8.7 | `Render(Summary) string` | Exact deterministic nonmutating bytes | `go test -run '^TestWeek08Render$' .` |
| G8 gate | `cli.go`, `cmd/logscope/main.go` supplied | Trace statuses/output; fixture parity and writer failures | `go test -run '^TestWeek08' ./...` |
| G9.1 | `reader.go`: `AnalyzeReader(io.Reader) (Summary, error)` | Read all, wrap failures, reject partial bytes plus error | `go test -run '^TestWeek09AnalyzeReader$' .` |
| G9.2 | `analyzeReadCloser(io.ReadCloser) (Summary, error)` | Own handle, defer one Close, preserve earlier failure | `go test -run '^TestWeek09ReadCloser$' .` |
| G9.3 | `AnalyzeFile(string) (Summary, error)`; `cli.go`: `analyzeCLIFile(string) (Summary, error)` | Open with wrapped error; transfer handle; replace adapter body with AnalyzeFile call | `go test -run '^TestWeek09(AnalyzeFile|ReaderFixtureParity|CLIExecutable)$' .` |
| G9 gate | Existing sequential CLI unchanged externally | Recheck weeks 8-9 | `go test -run '^TestWeek(08|09)' ./...` |
| G10.1 | `parallel.go`: `mergeSummaries(*Summary, Summary) error` | Checked independent-map merge | `go test -run '^TestWeek10(Merge|MergeOverflow|NoSourceMutation)$' .` |
| G10.2 | `analyzeParallel(context.Context, []byte, int, lineProcessor) (Summary, error)` | Bounded jobs, local state, cancel, collect every result, merge after join | `go test -run '^TestWeek10' .` |
| G10 gate | `AnalyzeParallel(context.Context, []byte, int) (Summary, error)` supplied adapter | Sequential equivalence; CLI remains sequential | `go test -run '^TestWeek(08|09|10)' ./...` |
| G11.1 | `server.go`: `NewSummaryHandler(string) http.Handler` | Immutable report; endpoint table below | `go test -run '^TestWeek11(Handler|ConcurrentHandler)$' .` |
| G11.2 | `newHTTPServer(http.Handler) *http.Server` | Dedicated handler, 5s header timeout | `go test -run '^TestWeek11HTTPConfiguration$' .` |
| G11.3 | `serveHTTP(context.Context, net.Listener, http.Handler) error` | Own listener, serve/join/shutdown/close | `go test -run '^TestWeek11(ServeAndCancel|GracefulActiveRequest|ForcedClose|ServeFailure)$' .` |
| G11.4 | `runServer` dependency seam below | Validate, analyze once, render, listen, serve; diagnostics/statuses | `go test -run '^TestWeek11(ServerStartup|ServerExecutableFailures)$' .` |
| G11 gate | `RunServer`, `cmd/logserve/main.go` supplied | Trace signal lifecycle and manual Ctrl+C evidence | `go test ./...` |

### Sequential command and resources

`RunCLI(arguments []string, stdout, stderr io.Writer) int` is supplied.
Arguments exclude the executable name. Exactly one positional path, including
a name starting `-`, is accepted. No flags/stdin. Argument count errors exit 2;
file, analysis, or stdout write errors exit 1; success exits 0. Failure stderr
starts `error:`; wording after it is not graded. Stderr failures are best effort.
On successful analysis, build the full report before writing; stdout may be
partial only if the write itself fails. A short write without an error is
still failure. On success stdout is exactly the report and stderr is empty.

The supplied week-8 `analyzeCLIFile` uses `os.ReadFile` and `Analyze`.
This deliberately avoids the future reader TODO. In G9.3 replace its entire
body with the call to `AnalyzeFile` and remove its now-unused `os` import.
Do not add a separate CLI mode or keep the read-all shortcut after week 9.

`AnalyzeReader` does not own/close its `io.Reader`.
`analyzeReadCloser` owns the successfully opened handle and closes exactly
once, including read/analysis failure paths. A close error is returned only
if analysis succeeded; otherwise preserve the earlier read/analysis cause,
not an error joined with the later close error. All error results are zero
summaries. `AnalyzeFile` wraps open errors and transfers a successful handle
to the close helper; it must not also close it. Explain named returns or
your small-helper choice in the week-9 reflection.

Tests use immediate/partial-error readers, one-byte reads, close counters,
failing/short writers, temp files, paths with spaces and starting `-`,
and real executables built once per command test into `t.TempDir`.
Build failure is fatal, not an accepted TODO failure. Process tests use
`os/exec`, never shells or `go run`, select `.exe` only on Windows, and
compare captured bytes/statuses without matching OS error text.
Practical files cannot force uint64 overflow; near-limit helper tests cover
arithmetic instead of a giant input or test-only CLI switch.

### Pool lifecycle and deterministic tests

The caller supplies a non-nil context. Validate workers in 1..4 **before**
checking `ctx.Err`; return `ErrWorkers` even for a canceled invalid-count call.
Check pre-cancellation before spawning or processing.

`lineProcessor` is exactly:

```go
type lineProcessor func(context.Context, *Summary, []byte) error
```

It is a private test seam, not a public alternate analysis API. The supplied
`AnalyzeParallel` uses `accumulateLine` through this callback. The callback
runs synchronously in its worker and receives the derived context, that
worker's exclusive local summary, and one normalized nonempty complete line.
No callback goroutines, shared summary, or panic recovery is required.

Implement the handout's algorithm, not goroutine-per-line:

- Start exactly `workers` worker goroutines. The calling goroutine is the
  producer; do not start an extra producer.
- Jobs channel capacity is at most `workers`; final-result channel capacity
  is exactly `workers`. The producer alone closes jobs.
- Each worker owns a local Summary and sends exactly one final result,
  including its error, before exiting. Final-result sends must not be skipped
  merely because cancellation is already signaled.
- Observe derived cancellation when sending/receiving jobs and between
  records. A worker processing error cancels the derived context, unblocking
  producer and peers. Do not cancel the caller's context.
- Collect all worker results and join all workers before returning, even on
  error or empty input. Merge only after joining, with checked arithmetic.
- A substantive worker error or checked-merge overflow wins over cancellation.
  Internal `context.Canceled` observations must not hide the initiating error.
  If there is no such failure, return the caller's `ctx.Err`, checking again
  just before success. Never return silently partial success.

Lifecycle tests use channel handshakes, separate local-summary identities,
and active/exit counters, not sleeps or assumptions about job distribution.
They wait for all four callbacks to enter before canceling or releasing the
injected worker failure, then require all to exit before the pool returns.
Ten-second deadlines are deadlock guards, not throughput thresholds.
Direct merge tests seed near-limit fields; do not allocate enormous fixtures.
Some requirements (channel capacities, exact goroutine creation/ownership,
and use of defer) also need source inspection and reflection; black-box tests
cannot prove every structural property.

### HTTP startup, immutable bytes, and shutdown

`NewSummaryHandler(report string)` accepts already rendered immutable text.
It must not parse uploads, read files, or capture mutable summary maps.

| Request | Status | Headers | Body |
| --- | --- | --- | --- |
| GET `/summary` (query string ignored) | 200 | `Content-Type: text/plain; charset=utf-8` | Exact report |
| Any other method at `/summary`, including HEAD | 405 | `Allow: GET` | Empty |
| Any other path, including `/summary/`, with any method | 404 | No `Allow` | Empty |

Unknown path checking takes precedence over method checking. Do not use
`http.Error` if it adds a body, a redirecting `/summary/` registration, or
global handler registrations. Concurrent requests only read captured text.

The exact private startup seam is:

```go
func runServer(ctx context.Context, arguments []string, stderr io.Writer,
    analyze func(string) (Summary, error),
    listen func(string, string) (net.Listener, error),
    serve func(context.Context, net.Listener, http.Handler) error) int
```

It accepts exactly one path, returning 2 for wrong count and 1 for
file/analysis/listen/serving/shutdown errors; diagnostics start `error:`.
Success is 0 with empty stderr. Analyze exactly once before listening, render
once, create the immutable handler, then call `listen("tcp", "127.0.0.1:8080")`.
Do not listen on bad input or refresh on requests. A successfully returned
listener transfers to `serve`, which owns closing it; no hidden second owner.
The injected listener function follows net.Listen's convention: no owned
listener on error.

`RunServer(context.Context, []string, io.Writer) int` supplies `AnalyzeFile`,
`net.Listen`, and `serveHTTP`. `cmd/logserve` supplies
`signal.NotifyContext(context.Background(), os.Interrupt)` and stops signal
notification before outer `os.Exit`. No tests send platform-dependent signals
or occupy production port 8080.

`newHTTPServer` installs the supplied handler and a five-second
`ReadHeaderTimeout`. `serveHTTP` owns the listener, serves in one goroutine,
and waits for that goroutine to finish. On context cancellation call Shutdown
using a **fresh** background context with a two-second timeout (not the
already-canceled serving context). Cancel the timeout context afterward.
If Shutdown fails, call Close to release active connections. Preserve the
shutdown failure as an error; forced close is not successful graceful shutdown.
Normalize the serving goroutine's `http.ErrServerClosed` to success, not
arbitrary other Serve errors. Do not wait for cancellation forever if Serve
has already failed. Release the listener on every path. No TLS or persistence.

Handler tests use httptest; lifecycle tests use `127.0.0.1:0`, real clients
with timeouts, synchronized active requests, and listener-close observations.
One forced-close test intentionally exercises the actual two-second deadline;
no assertions depend on a narrow elapsed-time window. Tests shut down their
servers and close response bodies/connections. Startup injection checks the
fixed production address, analysis-before-listen, single immutable startup
report, resource transfer, error statuses, and real executable failures.
Manual Ctrl+C/exit-0 evidence remains required by week 11.

## Build, compile-only, staged, full, and race commands

Run only your OS's block, from the course root. Use your existing stable Go
installation described in [SETUP](../../SETUP.md). These commands do not
install tools. If any build/compile command fails, stop and fix that before
interpreting acceptance failures.

### Windows PowerShell

```powershell
Push-Location .\go\logscope
go build ./...
if ($LASTEXITCODE -ne 0) { Pop-Location; throw 'Build failed.' }
go test -run '^$' ./...
if ($LASTEXITCODE -ne 0) { Pop-Location; throw 'Test compilation failed.' }
go test -tags parity -run '^$' ./...
if ($LASTEXITCODE -ne 0) { Pop-Location; throw 'Tagged test compilation failed.' }
gofmt -l .
go vet ./...
Pop-Location
```

Run the selector for your current week; do not run future gates as if due now:

```powershell
Push-Location .\go\logscope
go test -count=1 -run '^TestWeek08' ./...             # Week 8
# go test -count=1 -run '^TestWeek(08|09)' ./...      # Week 9
# go test -count=1 -run '^TestWeek(08|09|10)' ./...   # Week 10
# go test -race -count=1 -run '^TestWeek(08|09|10)' ./...
Pop-Location
```

After week 11, run the full untagged gate:

```powershell
Push-Location .\go\logscope
go test -count=1 ./...
go test -race -count=1 ./...
gofmt -l .
go vet ./...
Pop-Location
```

### macOS Terminal

```sh
(
  cd go/logscope &&
  go build ./... &&
  go test -run '^$' ./... &&
  go test -tags parity -run '^$' ./... &&
  gofmt -l . &&
  go vet ./...
)
```

Select your current week (uncomment only its appropriate cumulative gate):

```sh
(
  cd go/logscope &&
  go test -count=1 -run '^TestWeek08' ./...
  # go test -count=1 -run '^TestWeek(08|09)' ./...
  # go test -count=1 -run '^TestWeek(08|09|10)' ./...
  # go test -race -count=1 -run '^TestWeek(08|09|10)' ./...
)
```

After week 11:

```sh
(
  cd go/logscope &&
  go test -count=1 ./... &&
  go test -race -count=1 ./... &&
  gofmt -l . &&
  go vet ./...
)
```

`gofmt -l .` must list no files; it does not itself fail on differences.
`gofmt -w .` formats source. Check each PowerShell exit status before
proceeding; do not mistake a later success for earlier acceptance.
The week-10 race command excludes week-11 TODOs. Race requires supported
cgo/compiler setup; see [prerequisites](../../SETUP.md#race-detector-prerequisites).
If unavailable, record that validation is incomplete, not race-clean.

The supplied capstone parity test is opt-in under `-tags parity`; ordinary
week-8 through week-11 commands must not require Cargo. Tagged compile-only
checks do not execute the parity test. Week 12's actual cross-language command
and additional toolchain requirements are in its handout and parity spec.

Before each submission, record commands and actual results in
`docs/submissions/week-NN.md`, answer that week's C++ questions, and self-score
against its unchanged rubric. Completed weeks must pass their cumulative
gate. At week 11 no required TODO may remain on a required path. Do not claim
a green scaffold solely because build/test-compilation passed.
