# One paired exercise: statistics in Rust and Go

This is a supplied source scaffold, not a completed solution. Do the
Rust half in **week 2 after the week 1 gate**. Do the Go half in **week 7 after
weeks 1-6**, not immediately after Rust statistics. The assignments'
reading budgets, prerequisites, reflections, and rubrics still apply:
[week 2](../../assignments/week-02.md), [week 7](../../assignments/week-07.md).
Other supplied projects have their own
[specs](../../README.md#project-specs-and-implementation-guides).

## What is supplied and what you implement

Each project is independent, standard-library-only, and has its own manifest.
There is no root Cargo workspace or Go module. Do not run `cargo new` or
`go mod init` over these supplied directories.

| Location | Purpose |
| --- | --- |
| `rust/week02-stats/Cargo.toml` | Rust 2021 binary + library crate |
| `rust/week02-stats/src/lib.rs` | Types, R1-R6 learner functions, supplied error/process adapter |
| `rust/week02-stats/src/main.rs` | Supplied executable entry point |
| `rust/week02-stats/tests/stats.rs` | Starter, arithmetic, statistics, caller ownership tests |
| `rust/week02-stats/tests/cli.rs` | Grammar, bytes, runner, and real executable tests |
| `go/stats/go.mod` | Go module, minimum language version 1.22 |
| `go/stats/stats.go` | Types and G1-G2 learner functions |
| `go/stats/cli.go` | G3-G5 learner functions and supplied process adapter |
| `go/stats/main.go` | Supplied executable entry point |
| `go/stats/aliasing.go` | G6a-G6g learner lab functions |
| `go/stats/*_test.go` | Statistics, CLI/executable, and aliasing acceptance tests |

Implement only the marked function bodies; keep public signatures, error
variants/sentinels, and test contracts stable. The `main` TODOs ask you to trace
the supplied wiring, not replace it. Types, adapters, tests, and manifests are
infrastructure, not exercises to invent. Report formatting returns a Result/error
only so the starter can signal "unimplemented" without pretending to succeed;
completed formatting always succeeds for a valid Stats.

The starter is designed to **compile**, but acceptance tests are **active and
intentionally fail**. Result/error-returning stubs report Unimplemented;
Rust `classify` and Go lab stubs explicitly panic. The Go lab test turns each
panic into a failed subtest so all seven tasks remain visible. These are starter
placeholders, not acceptable completed behavior. Do not skip tests or change
expected results to get green. Build errors, missing tools, and linker errors
are not intentional TODO failures.

The generated `Cargo.lock` is included; retain it for this application.
No lockfile is hand-authored in this scaffold.
No `go.sum` is needed without dependencies.

## Shared contract

### Pure statistics API (larger range than the CLI)

Rust `calculate(&[u32]) -> Result<Stats, StatsError>` and Go
`Calculate([]uint32) (Stats, error)`:

- Accept the entire unsigned 32-bit range, including 4294967295.
- Do not print, mutate input, or take ownership of caller storage.
- Count and sum are unsigned 64-bit. Check both increments and additions;
  wrapping, saturation, and panic are not substitutes for overflow errors.
- Empty input: count/sum zero, no minimum or maximum. Go nil and empty slices
  have identical statistics.
- Nonempty input: both extrema present, even when both are zero.
- Rust uses `Option<u32>`. Go uses `HasMin` and `HasMax`; both are false for
  empty input and true otherwise. Absent Go numeric fields are zero.
- On arithmetic overflow return `StatsError::Overflow` or `ErrOverflow`.
  Go returns `Stats{}` on error, not a partial aggregate.

`checked_add` / `CheckedAdd` accepts two u64/uint64 operands. It returns their
exact sum when representable (including MAX + 0), otherwise the overflow error.
Go's numeric result on error is zero. Use this helper for count and sum.
The tests exercise the helper near the 64-bit limit without giant allocations.
They do not claim to force aggregate overflow with a practical input slice.

### CLI grammar, output, and statuses

Pass zero or more arguments, excluding the program name. Each token must:

1. Contain exactly 1-4 **ASCII** digits (`0` through `9`).
2. Have numeric value in the inclusive range 0..1000.

Leading zeros are allowed: `0001` means 1 and `0000` means 0. Do not trim.
Reject signs, whitespace anywhere within a token, decimals, empty strings,
non-ASCII digits, alphabetic characters, five digits (even `00000`), and 1001.
The process tests pass arguments directly, not through shell parsing.

Validate every argument before producing output. `Run`/`run` returns a complete
report or an error, never a partial report. A valid number before an invalid
one must not leak a report. Parsing errors are exactly
`StatsError::InvalidArgument` / `ErrInvalidArgument`; Go returns numeric zero.

| Condition | Status | stdout | stderr |
| --- | --- | --- | --- |
| Success | 0 | Exactly one report with final LF | Empty |
| Invalid argument | 2 | Empty | Starts with `error:` |
| Arithmetic failure | 1 | Empty | Starts with `error:` |

The supplied adapters add an LF after diagnostics. Diagnostic wording after
`error:` is not graded. Completed behavior must not panic. Broken pipes, invalid
OS argument encodings, and other host I/O failures are outside this exercise;
the thin supplied entry points assume writable streams and text arguments.

Exact successful bytes for `4 0 9 4`:

```text
count=4 sum=17 min=0 max=9
```

Exact successful bytes for no arguments:

```text
count=0 sum=0 min=none max=none
```

Each line above ends in **one LF byte**, not CRLF, and has no leading/trailing
spaces. Field order, lowercase `none`, and decimal notation are fixed.
The formatter also accepts full-range valid Stats, not just CLI-sized values.
Inconsistent hand-constructed Stats are outside its precondition.

Practical valid CLI arguments cannot reach a 64-bit arithmetic overflow.
Tests check real process statuses 0/2, checked-add overflow, and status-1
adapter wiring directly. They do not fabricate an executable overflow test or
add a test-only CLI switch.

## Rust implementation order

Use explicit `match` for Option and Result this week. The `?` operator,
composed errors, and a new validation framework are not needed.

| Step | File/function | Work | Focused test (inside Rust project) |
| --- | --- | --- | --- |
| R1 | `src/lib.rs`: `classify` | Empty/nonempty enum with count | `cargo test --test stats starter_classifies_slices` |
| R2 | `checked_add` | Match checked addition | `cargo test --test stats checked_add_boundaries` |
| R3 | `calculate` | Borrow, count, sum, extrema | `cargo test --test stats statistics_cases` |
| R3 check | `calculate` | Caller vector remains readable | `cargo test --test stats caller_keeps_vector` |
| R4 | `parse_arg` | Validate token bytes, length, range | `cargo test --test cli parse_grammar` |
| R5 | `format_report` | Exact field order/absence/LF | `cargo test --test cli report_bytes` |
| R6 | `run`; trace `src/main.rs` | Parse all, calculate, render | `cargo test --test cli` |

For the Starter reflection, identify the array, Vec, and slice in
`starter_classifies_slices`: which owns storage, and which is a view?
A borrowed slice resembles `std::span<const uint32_t>`, not a copied
`std::vector`. `Option` resembles `std::optional`; the Dataset enum resembles
a `std::variant` with explicit alternatives. The tests exhaustively match both
variants, but you still explain why empty differs from a zero-valued element.

## Go implementation order

| Step | File/function | Work | Focused test (inside Go project) |
| --- | --- | --- | --- |
| G1 | `stats.go`: `CheckedAdd` | Check before unsigned addition | `go test -run '^TestCheckedAdd$'` |
| G2 | `Calculate` | Full-range, non-mutating statistics | `go test -run '^TestCalculate$'` |
| G2 check | `Calculate` | Caller slice remains unchanged | `go test -run '^TestCallerKeepsSlice$'` |
| G3 | `cli.go`: `ParseArg` | Match the Rust CLI grammar | `go test -run '^TestParseArg$'` |
| G4 | `FormatReport` | Match the exact report bytes | `go test -run '^TestFormatReport$'` |
| G5 | `Run`; trace `main.go` | Full validation, no partial output | `go test -run '^TestRun$'` |
| G5 check | Supplied executable | Real streams and statuses | `go test -run '^TestCLIExecutable$'` |
| G6 | `aliasing.go`: seven functions below | Observe copies and aliases | `go test -run '^TestAliasing$'` |

Go lab inputs are initialized maps, non-nil pointers, and ordinary slices.
No nil-map assignment panic is required. The tests do the mutations on either
side, so returning equal-looking data without the right sharing is insufficient.

| Step/function | Required observation | Subtest |
| --- | --- | --- |
| G6a `CopyScalar` | Returned scalar changes; original remains independent | `value_copy` |
| G6b `MutateScalar` | Pointer mutation changes the original | `pointer_mutation` |
| G6c `ShareSlice` | Same-length view shares element mutations both ways | `shared_slice` |
| G6d `CopySlice` | Equal elements in a separate backing array | `copied_slice` |
| G6e `ShareMap` | Assignments through either variable share entries | `shared_map` |
| G6f `Lookup` | Missing key differs from present zero | `map_presence` |
| G6g `AppendWithinCapacity` | Returned header grows; old length stays; storage shared | `append_capacity` |

For a single lab use, for example,
`go test -run '^TestAliasing$/^append_capacity$'`.
The append test deliberately uses `make([]uint32, 1, 3)` and a full-capacity
view, so its result does not depend on Go's allocation growth policy.
Explain why copying a slice header resembles copying a span, not copying a
vector, and why garbage collection does not prevent aliasing.

## Build, compile tests, and run acceptance

Use the stable toolchains from [SETUP](../../SETUP.md). No installs or external
packages are part of these commands. Go must be available on PATH because its
executable test invokes `go build`. Both projects are independent; Go can wait
until week 7.

Run each line separately and inspect its exit status. **Build and test
compilation should pass even before implementation.** Full acceptance should
fail until the TODOs are completed. Do not confuse `--no-run` / `-c` with
passing acceptance. Rust executable tests use Cargo's built binary; Go builds
its executable once into `t.TempDir()` and cleans it automatically.

### Windows PowerShell (from course root)

```powershell
Push-Location .\rust\week02-stats
cargo build --offline
cargo test --offline --no-run
cargo test --offline
cargo fmt --check
cargo run --offline -- 4 0 9 4
Pop-Location

# Do this half in week 7.
Push-Location .\go\stats
$env:GOTOOLCHAIN = 'local' # Do not auto-download a newer Go toolchain.
New-Item -ItemType Directory -Force .\bin
go build -o bin/stats.exe .
go test -c -o bin/stats.test.exe .
go test -count=1 ./...
gofmt -l .
go vet ./...
.\bin\stats.exe 4 0 9 4
$LASTEXITCODE
Pop-Location
```

### macOS Terminal (from course root)

```sh
(
  cd rust/week02-stats || exit
  cargo build --offline
  cargo test --offline --no-run
  cargo test --offline
  cargo fmt --check
  cargo run --offline -- 4 0 9 4
)

# Do this half in week 7.
(
  cd go/stats || exit
  export GOTOOLCHAIN=local
  mkdir -p bin
  go build -o bin/stats .
  go test -c -o bin/stats.test .
  go test -count=1 ./...
  gofmt -l .
  go vet ./...
  ./bin/stats 4 0 9 4
  echo $?
)
```

`gofmt -l .` must print no paths; use `gofmt -w .` to apply formatting.
Use `cargo fmt` to apply Rust formatting. Direct binaries, not `go run`, are
the source of truth for program exit codes (`go run` may translate the status).
The process tests also avoid PowerShell's version-specific handling of empty
native arguments.

The supplied cases cover empty/singleton/duplicate/zero/ordered/mixed datasets,
maximum 32-bit inputs and wide sums, near-limit checked addition, grammar,
exact report bytes, stream isolation, actual executable statuses, Rust borrowing
and enum matching, and Go aliasing. Add your own cases as you learn; keep these
acceptance cases active.

## Completion and reflection

Before calling a half complete: build it, compile all tests, pass the full test
suite, pass formatting (and Go vet), and inspect the sample report. No TODO
stub may remain on a required path. Do not ship a core that only returns errors
or a lab that still panics.

Submit the project and `docs/submissions/week-02.md` or
`docs/submissions/week-07.md` as the respective assignment requires. Record
actual tool versions/results and explain ownership/views, absence versus zero,
checked arithmetic, validation-before-output, and (in Go) aliasing. These
reflections are learner deliverables, not prewritten solutions.

This paired scaffold stops here. Further example projects require approval.
