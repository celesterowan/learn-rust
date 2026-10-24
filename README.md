# Rust and Go for a C++ programmer

A 12-week, assignment-driven course for someone who knows C++ but has never
worked in Rust or Go. Budget **5-7 hours per week**. Learn the languages by
building small programs, writing tests, and explaining how your design differs
from the C++ version you would normally write.

The main project is a small log analyzer implemented in both languages.
Rust comes first, Go second, and the final week compares your implementations.
Use the standard libraries only. No frameworks, paid tools, or prior web
development experience are required.

**Current contents:** this README is the course overview and initial assignment
guide. Detailed assignment handouts, shared fixtures, and reference notes are
not included yet. There are no starter implementations or automated graders.

## Learning goals

By the end, you should be able to:

- Explain Rust ownership, borrowing, moves, and lifetimes without treating them
  as C++ pointers with different syntax.
- Use Rust enums, `Option`, `Result`, traits, iterators, and tests.
- Explain Go values, pointers, slice aliasing, interfaces, and explicit errors.
- Manage resources explicitly in both languages.
- Build bounded concurrent work and handle errors without losing results.
- Test a small Go HTTP handler and shut down its local server.
- Compare language trade-offs using examples from your own code.

## Setup

Use a separate checkout of this repository, not a directory inside another
project. Install current stable tools:

- [Rust installation](https://www.rust-lang.org/tools/install) and the
  [Rust Book](https://doc.rust-lang.org/book/).
- [Go installation](https://go.dev/doc/install) and
  [A Tour of Go](https://go.dev/tour/).
- Git and a text editor of your choice.

On Windows, Rust's MSVC toolchain requires Visual Studio C++ build tools and a
Windows SDK. Your existing C++ installation may already provide them.
Reopen your terminal after installing tools, then check:

```powershell
rustc --version
cargo --version
go version
git --version
```

Create implementation directories as the assignments introduce them:

```text
rust/
  week01-hello/    independent Cargo binary crate
  week02-stats/    independent Cargo binary crate
  logscope/        Rust analyzer, weeks 3-6 and 12
go/
  stats/           independent Go module
  logscope/        Go analyzer, weeks 8-12
docs/
  submissions/    weekly reflections and self-assessments
```

These are planned directories, not included starter projects. Each Rust project
has its own Cargo manifest; each Go project has its own Go module. Run commands
inside the corresponding project, not at the repository root.

## Weekly routine and grading

Aim for six hours: **45 minutes reading, 30 minutes planning, two hours
implementing, 90 minutes testing, and 75 minutes reviewing and reflecting**.
Use a seventh hour for catch-up. Stretch work is optional and ungraded.

Each week, submit your source and tests plus
`docs/submissions/week-NN.md` containing:

1. What you built and the commands you ran.
2. Observed results, including failures or unfinished requirements.
3. One C++ habit that helped and one that got in the way.
4. Answers to the week's reflection question.
5. Your self-score and the evidence supporting it.

Unless a detailed handout specifies otherwise, grade each week out of 100:
**required behavior 40, tests and edge cases 30, idiomatic design 15,
reflection and reproducibility 15**. Target 80 or higher. A plausible happy-path
demo does not replace tests. Keep previous projects buildable.

Commit one completed assignment at a time when you are ready. Suggested message:
`Complete week 01: Rust tooling and greeting`. Do not commit build output,
executables, secrets, or machine-specific paths.

## The 12 assignments

### Week 1: Rust tooling and a greeting

**Read:** Rust Book chapters 1-3. Skim the Go installation guide, but focus your
coding time on Rust.

Create `rust/week01-hello` with Cargo. Write a CLI that accepts exactly one
name and prints `Hello, NAME!` followed by a newline. Accept a quoted name with
spaces. Missing or extra arguments must produce a diagnostic on stderr and a
nonzero exit status, without panicking. Test both formatting and argument
validation. Record the tool versions and commands needed to reproduce your run.

**Reflect:** How do Cargo, modules, and tests differ from your C++ build setup?

### Week 2: Ownership through number statistics

**Read:** Rust Book chapter 4 and the sections on vectors and slices.

Create `rust/week02-stats`. Implement a function taking a borrowed slice of
`i32` values and returning count, minimum, maximum, and an `i64` sum. Empty
input has no minimum or maximum. Detect sum overflow instead of wrapping.
Keep parsing separate from calculation. Test empty input, a single value,
negative values, and duplicates. Test an internal checked-add helper with a
near-limit accumulator rather than allocating an enormous slice.

**Reflect:** Where would you normally copy a container or pass a reference?
Which choices does Rust make explicit?

### Week 3: Rust enums, errors, and parsing

**Read:** Rust Book chapters 6 and 9, plus the lifetime syntax in chapter 10.

Start `rust/logscope` as a library and a small binary. Parse one record in the
form `INFO api 20`: level, service, and duration in milliseconds. Represent
levels with an enum and parsing failures with an explicit error type.
Return `Result`, not a panic, for malformed input. First try borrowing the
service name from the input, then explain when an owned string is useful.
Test missing fields, extra fields, invalid levels, and invalid durations.

**Reflect:** How does `Result` compare with exceptions and output parameters?

### Week 4: A deterministic Rust log analyzer

**Read:** Rust Book sections on hash maps, iterators, and test organization.

Extend the parser into a CLI accepting one file path. Count valid and invalid
records, totals per level, and per-service count, total duration, and maximum
duration. Sort service names before rendering. Test empty files, malformed
lines, CRLF, a missing final newline, and file errors. Write a short input/output
contract before coding; freeze it for your future Go implementation.

**Reflect:** Why is hash-map iteration order unsuitable for reproducible output?

### Week 5: Rust traits and testable boundaries

**Read:** Rust Book chapter 10's traits section and chapter 11.

Separate reading, parsing, aggregation, and rendering. Use `std::io::Write`
as a boundary for report output and test it with an in-memory buffer and a
deliberately failing writer. Propagate write errors. Add checked arithmetic to
summary updates and test near-limit state directly. Preserve the week 4 CLI
contract. Avoid creating a custom trait where a standard-library trait suffices.

**Reflect:** Compare trait bounds with C++ templates and virtual interfaces.

### Week 6: Safe, bounded Rust concurrency

**Read:** Standard-library documentation for
[`std::thread::scope`](https://doc.rust-lang.org/std/thread/fn.scope.html).

Add a library-only parallel analyzer accepting 1-4 workers. Partition complete
records, give each worker its own summary, and merge with checked arithmetic.
Reject invalid worker counts. Join every worker and surface worker failures.
Compare results with the sequential implementation, including empty input and
fewer records than workers. Keep the existing CLI unchanged.

**Reflect:** What prevents a borrowed input buffer from outliving its owner?

### Week 7: Go values, slices, and maps

**Read:** A Tour of Go's basics, methods, and interfaces sections.

Create `go/stats` as a Go module and reimplement week 2 using `[]int32` and an
`int64` sum, with equivalent empty-input and overflow behavior. Write
table-driven tests. Add a test demonstrating shared backing storage between
slices and another demonstrating an independent copy.

**Reflect:** Why are a Go slice, a C++ vector, and a Rust vector not equivalent?

### Week 8: The same analyzer in Go

**Read:** Go's
[error-handling guidance](https://go.dev/blog/error-handling-and-go)
and the documentation for `os.ReadFile`.

Create `go/logscope` and implement your frozen week 4 contract in Go.
Use explicit errors for file and analysis failures, and count malformed records
according to the contract. Match deterministic output exactly. Reuse the same
inputs and expected results, not the Rust implementation's source structure.

**Reflect:** Which errors should carry context, and where should you print them?

### Week 9: Go interfaces and resource cleanup

**Read:** The Tour's interfaces and `defer` sections; documentation for
`io.Reader`, `io.Writer`, and `testing`.

Refactor around standard-library reader/writer interfaces where useful.
Test short reads, read failures, and write failures with small test doubles.
When opening files explicitly, ensure they close on all return paths.
Preserve CLI behavior and keep interfaces small and consumer-driven.

**Reflect:** Which resource lifetimes still require explicit action despite GC?

### Week 10: Go workers and cancellation

**Read:** The Tour's concurrency section,
[`context`](https://pkg.go.dev/context), and the
[race detector guide](https://go.dev/doc/articles/race_detector.html).

Add a library-only analyzer with 1-4 workers. Give workers private summaries,
bound queued work, and use context cancellation. On failure or cancellation,
stop sending work, wait for every worker, and return an error rather than a
partial success. Test invalid counts, cancellation, worker errors, and
equivalence with sequential results. Use synchronization, not sleeps, in tests.

Run race tests on a supported toolchain. On Windows they require a compatible
C compiler in addition to Go; ordinary tests passing is not a race-clean result.

**Reflect:** Who closes each channel, and what proves every worker has stopped?

### Week 11: A small local Go HTTP server

**Read:** Documentation for `net/http`, `net/http/httptest`,
`os/signal.NotifyContext`, and `http.Server.Shutdown`.

Add a separate server command. Read and analyze one file at startup, then serve
the immutable report at GET `/summary`. Bind only to `127.0.0.1:8080`.
Define and test unknown-path and unsupported-method behavior. Use `httptest`
without depending on an available port. Handle Ctrl+C, use a bounded shutdown
timeout, and wait for the serving goroutine. No uploads or deployment.

**Reflect:** Why is immutable captured data easier to share between requests?

### Week 12: Cross-language capstone

No new product. Build a parity test harness that runs both compiled CLIs on the
same fixtures. Compare bytes and exit status against expected results, not
only against each other. Cover empty input, malformed input, CRLF, missing
files, and paths with spaces. Rerun checks in every project.

Write `docs/comparison.md` (700-1000 words) using concrete examples from your
work: ownership and aliasing, error propagation, resource cleanup, interfaces,
testing, and concurrency. Explain which language you would choose for this
project and why. Separate measured evidence from preference.

**Reflect:** Which C++ habits did you stop applying, and which remain valuable?

## Checks before submitting

Inside each Rust crate:

```powershell
cargo fmt --check
cargo test
cargo clippy --all-targets -- -D warnings
```

Install the `rustfmt` and `clippy` components through `rustup` if missing.
Inside each Go module:

```powershell
gofmt -l .
go vet ./...
go test ./...
```

`gofmt -l .` must produce an empty listing; it does not fail merely because
formatting differs. Use `gofmt -w .` to format. Add `go test -race ./...` for
concurrent exercises on a supported environment. Record any unavailable check
as incomplete rather than claiming success.

## Working rules

Do not translate C++ line by line. Prefer each language's ordinary idioms,
and explain unfamiliar compiler errors in your reflections. Documentation,
discussion, and help debugging are encouraged; acknowledge assistance and
write your own implementations.

Keep scope small: no Rust async, unsafe code, FFI, third-party packages, Go
frameworks, or production infrastructure. Whole-file reads are acceptable
for these small inputs. Parallel code need not be faster than sequential code.
Optional benchmarks must describe their setup and limitations.
