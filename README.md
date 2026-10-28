# Rust and Go for a C++ programmer

A 12-week, assignment-driven course for someone who knows C++ but has never
worked in Rust or Go. Budget **5-7 hours per week**. Learn the languages by
building small programs, writing tests, and explaining how your design differs
from the C++ version you would normally write.

The main project is a small log analyzer implemented in both languages.
Rust comes first, Go second, and the final week compares your implementations.
Use the standard libraries only. No frameworks, paid tools, or prior web
development experience are required.

This is an **independent companion** to Microsoft's
[Rust for C/C++ Programmers](https://microsoft.github.io/RustTraining/c-cpp-book/),
not a complete official Microsoft curriculum or an endorsed course. The Rust
sequence and exercise method follow selected parts of that guide; the
assignments below are original. Go is our separate track using official Go
documentation.

**Included:** all twelve weekly handouts, setup and reference guides, shared
fixtures, and a reflection template. This README also contains the assignments
inline; the chapter-to-week map links to their individual handouts.
Start with [SETUP.md](SETUP.md), keep the
[C++ comparison guide](CPP-TO-RUST-AND-GO.md) nearby, and adapt the
[reflection template](docs/templates/reflection.md) for each submission.
The [shared logscope contract](docs/logscope-contract.md) defines exact analyzer
behavior. All required source projects now include buildable starter
infrastructure, function signatures, learner TODOs, and active acceptance tests.
Completed implementations and learner submissions are not provided.

## Project specs and implementation guides

Start each project with its spec, then follow that week's reading and rubric.
Each guide maps the implementation order to specific files, functions, and tests.

| Weeks | Supplied project | Spec |
| --- | --- | --- |
| 1 | Rust greeting and control-flow lab | [Greeting](docs/exercises/rust-hello.md) |
| 2 and 7 | Matching Rust and Go statistics, plus Go aliasing labs | [Statistics](docs/exercises/stats.md) |
| 3-6 | Rust records, analyzer, I/O traits, and scoped workers | [Rust logscope](docs/exercises/rust-logscope.md) |
| 8-11 | Go analyzer, interfaces, cancellable workers, and local HTTP server | [Go logscope](docs/exercises/go-logscope.md) |
| 12 | Supplied cross-language executable acceptance harness | [Parity capstone](docs/exercises/parity.md) |

Builds and test compilation should pass before you implement anything.
Behavioral acceptance tests are intentionally active and failing at learner
TODOs. Use the week's focused test commands; later-week TODOs are not due early.
By each gate, its acceptance tests must pass. Never skip tests or weaken
expected results to hide an unfinished implementation.

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

Start with the detailed instructions for
**[Windows PowerShell](SETUP.md#windows-powershell)** or
**[macOS Terminal](SETUP.md#macos-terminal)**. Install or verify the following
tools before the relevant week; do not install another copy if yours works.

| Download | Windows | macOS |
| --- | --- | --- |
| Git | [Git for Windows](https://git-scm.com/install/windows) | Apple's Command Line Tools, requested with `xcode-select --install` |
| Rust and Cargo | [Stable rustup installer](https://www.rust-lang.org/tools/install) | [Official rustup instructions](https://doc.rust-lang.org/book/ch01-01-installation.html) |
| Native linker/SDK | MSVC **Desktop development with C++** workload and Windows SDK | Apple's Command Line Tools; full Xcode is not required |
| Go, by week 7 | [Stable Windows MSI](https://go.dev/dl/) for your architecture | [Stable macOS PKG](https://go.dev/dl/): arm64 for Apple silicon, amd64 for Intel |

An editor is required, but IDEs, Homebrew, GitHub CLI, and editor extensions
are optional. Rustup includes Cargo; add rustfmt and clippy through rustup.
The exercises use the standard libraries only. You do not need to download
completed projects, third-party packages, or Microsoft's training repository.
See [SETUP.md](SETUP.md) for installation, PATH, architecture, and race-detector
details. Reopen your terminal after installing tools.

These verification commands work in **both shells**; defer the Go command if
you are not installing Go until week 7:

```text
git --version
rustup show
rustc --version
cargo --version
rustup component add rustfmt clippy
go version
```

### Download the course

Already have this checkout with its handouts and fixtures? Use it; skip cloning.
Otherwise, choose a separate user-writable directory, not one inside Chromium
or another repository. These examples use your home directory's Projects folder.
To submit through your own fork, first fork on GitHub and replace the URL below
with **your fork's Code > HTTPS URL**. Cloning downloads a local checkout; it
does not create a fork or grant permission to push to the original repository.

**Windows PowerShell**

```powershell
New-Item -ItemType Directory -Force "$HOME\Projects" -ErrorAction Stop
Set-Location "$HOME\Projects" -ErrorAction Stop
git clone https://github.com/celesterowan/learn-rust.git learn-rust
if ($LASTEXITCODE -ne 0) { throw 'Clone failed; resolve the error before continuing.' }
Set-Location .\learn-rust -ErrorAction Stop
```

**macOS Terminal**

```sh
mkdir -p "$HOME/Projects" &&
cd "$HOME/Projects" &&
git clone https://github.com/celesterowan/learn-rust.git learn-rust &&
cd learn-rust
```

If the destination exists, use the existing checkout or choose another name,
not a nested clone. You should see README.md, SETUP.md, assignments, and fixtures.
If the selected branch lacks these materials, obtain the published course
branch/version rather than creating missing handouts.

**ZIP alternative:** choose **Code > Download ZIP** on GitHub, then fully
extract it using **Extract All** on Windows or Finder on macOS. Open the
extracted directory containing README.md. A ZIP has no Git history/remotes;
cloning is the recommended route for later commits and pushes.

Continue with [week 1](assignments/week-01.md). Handouts and fixtures are
included, along with all source scaffolds and tests linked above. Implement
the supplied TODOs as each assignment directs; no extra project download or
manual module initialization is needed.

### Course directories and shell commands

Run only your OS's command block. Examples start at the **course root**
(the folder containing this README), unless stated otherwise. Windows uses
PowerShell; macOS uses Terminal's zsh. Parenthesized macOS blocks leave your
shell at the course root. Read blocks before running them: open the supplied
project, implement the assignment, then run its focused checks. Stop and resolve failures
before continuing. [SETUP.md](SETUP.md#shell-conventions-and-directory-plan)
explains the shell conventions.

```text
rust/
  week01-hello/    independent Cargo binary crate
  week02-stats/    independent Cargo binary crate
  logscope/        staged TODOs: records week 3 through workers week 6
go/
  stats/           independent Go module
  logscope/        Go analyzer, weeks 8-12
    cmd/logscope/  sequential CLI
    cmd/logserve/  separate HTTP command, TODOs due week 11
fixtures/          provided shared inputs and expected outputs, week 4
docs/
  submissions/    weekly reflections and self-assessments
```

All source directories above are supplied; create docs/submissions/ for your
own reflections. Do not run cargo new or go mod init over the scaffolds.
Each Rust project
has its own Cargo manifest; each Go project has its own Go module. Run commands
inside the corresponding project, not at the repository root.
There is no root Cargo workspace or Go module. Do not follow the source's
workspace or dependency exercises for this course. Keep Cargo.lock for these
application exercises; go.sum may not exist without external dependencies.

Create the submission directory once from the course root.

**Windows PowerShell**

```powershell
New-Item -ItemType Directory -Force .\docs\submissions
```

**macOS Terminal**

```sh
mkdir -p docs/submissions
```

## Microsoft chapter-to-week map

The source's self-study pacing is not our time budget. Select the named
concepts and stop reading when the allotted time is spent; completing all its
exercises, sidebars, and dependency examples is not part of this course.

| Week | Microsoft chapters | Course checkpoint |
| --- | --- | --- |
| [1](assignments/week-01.md) | 1-4: motivation, getting started, types, control flow | Build/test a greeting; explain bindings and expression results |
| [2](assignments/week-02.md) | 5-6: collections, references, strings, enums/match | Compute stats; model absence and match alternatives |
| [3](assignments/week-03.md) | 7 and selected lifetime subchapter | Distinguish moves/Copy/Clone/Drop; borrow typed records |
| [4](assignments/week-04.md) | 8-9 and selected testing/error subchapters | Modules, Result parser, frozen analyzer contract |
| [5](assignments/week-05.md) | 10-12: traits/generics, conversions, closures/iterators | Generic I/O boundary with unchanged output |
| [6](assignments/week-06.md) | 13: safe concurrency | Bounded scoped workers, checked merge, all joins |
| [7](assignments/week-07.md) | Separate Go track: values, slices, maps | Statistics parity and aliasing |
| [8](assignments/week-08.md) | Separate Go track: errors and files | Deterministic Go analyzer |
| [9](assignments/week-09.md) | Separate Go track: interfaces and resources | Testable I/O and explicit cleanup |
| [10](assignments/week-10.md) | Separate Go track: goroutines and context | Bounded workers and cancellation |
| [11](assignments/week-11.md) | Separate Go track: net/http | Local summary server |
| [12](assignments/week-12.md) | Revisit Rust 7-13; short ch16 discussion | Parity and comparison using existing projects |

Basic tests and CLI containers appear early as small previews. Week 2 previews
Option and tiny Result/checked-add matches for numeric input, not full error
composition. The byte parser and its malformed-input matrix wait until week 4.
Cell/RefCell graphs, type-state/builders, advanced iterators, and additional
source exercise sets are not added to the six Rust weeks.

## Weekly routine and grading

Every week's budget below totals six hours; use a seventh for catch-up.
Rust weeks follow **C++ mental model -> Rust concept -> required Starter ->
required Intermediate project -> optional Challenge -> reflection**. The tiers
divide existing work, not double it. Challenge work is ungraded.
Follow the linked spec's TODO order for each Starter and Intermediate task.
Build and compile tests first, then run the week's focused acceptance checks
while implementing. Run full acceptance only when all stages are due.

Following Microsoft's self-study method, attempt before expanding a source
solution. If stuck for more than 15 minutes, inspect a relevant example or
source solution, close it, and reimplement from memory. These original
assignments have no supplied solutions. Record assistance; do not submit
copied Microsoft exercises or solutions. Our std-only, explicit-error contract
overrides context-specific examples that use crates or default failed values.

Each week, submit your source and tests plus
`docs/submissions/week-NN.md` containing:

1. What you built and the commands you ran.
2. Observed results, including failures or unfinished requirements.
3. One C++ habit that helped and one that got in the way.
4. Answers to the week's reflection question.
5. Your self-score and the evidence supporting it.
6. In Rust weeks, Starter/Intermediate results and gate evidence; note any
   source example consulted and what you reimplemented after closing it.

Use each week's 100-point rubric. Target 80 or higher, but pass the explicit
Rust gate before the next week even if your score is already 80. Repair failed
required checks; a plausible happy-path demo does not replace tests.
Keep previous projects buildable. Read only the selected source sections.

Commit one completed assignment at a time when you are ready. Suggested message:
`Complete week 01: Rust tooling and greeting`. Do not commit build output,
executables, secrets, or machine-specific paths.

## Shared analyzer contract and fixtures

The [shared contract](docs/logscope-contract.md) is required for both logscope
CLIs and frozen in week 4. It is not Microsoft's exercise format.
The files described below are provided under fixtures. Preserve their ASCII/LF
bytes; avoid shell redirection that can change encoding or newlines.
The included .gitattributes preserves LF across Git checkouts. Verify exact
expected bytes in tests.

Input is bytes: split on LF and remove at most one trailing CR from each
resulting line, including the final unterminated line. Ignore lines empty
after that removal. A final LF adds no record; an empty file has zero records.
Every nonempty line must be ASCII and exactly `LEVEL service duration`:

- Exactly one ASCII space between fields; no trim, tabs, leading/trailing
  spaces, extra fields, or comments.
- LEVEL: INFO, WARN, or ERROR, case-sensitive.
- service: `[a-z][a-z0-9-]{0,15}` (1-16 ASCII bytes).
- duration: 1-5 ASCII digits, value 0..60000 inclusive; leading zeros allowed.
  No signs, units, or decimal points.

Every malformed nonempty line increments invalid once and contributes nothing
else, including non-ASCII, invalid UTF-8, NUL, or whitespace-only lines.
Do not decode the whole input as UTF-8 or split on Unicode whitespace.

Aggregate valid/invalid counts, per-level counts, and per-service count,
duration total, and maximum. Use u64/uint64 counts and sums; detect arithmetic
overflow and return an error. Sort unique service keys by ascending ASCII name,
never map iteration order. Include all levels even when zero. Output is ASCII
with LF separators and exactly one final LF, no extra blank line.

`fixtures/sample.log` contains these lines (including the one blank line):

```text
INFO api 10
WARN api 20
ERROR worker 100

INFO worker 0
ERROR auth 7
INFO api 15
DEBUG api 4
INFO Api 2
INFO api -1
```

Its complete `fixtures/sample.expected.txt` is:

```text
valid=6 invalid=3
level INFO=3 WARN=1 ERROR=2
service api count=3 total_ms=45 max_ms=20
service auth count=1 total_ms=7 max_ms=7
service worker count=2 total_ms=100 max_ms=100
```

The provided `fixtures/empty.log` is zero bytes. Its
`fixtures/empty.expected.txt` is:

```text
valid=0 invalid=0
level INFO=0 WARN=0 ERROR=0
```

`fixtures/malformed.log` contains just the final three lines of sample.log.
Its `fixtures/malformed.expected.txt` is:

```text
valid=0 invalid=3
level INFO=0 WARN=0 ERROR=0
```

Each CLI accepts exactly one positional file path, no flags or stdin mode.
Quote paths with spaces; a name beginning with `-` is still a path.
Wrong argument count: exit 2, empty stdout, stderr starting `error:`.
File read, analysis overflow, or stdout write failure: exit 1 and `error:`.
Success, even with malformed records: exit 0, exact report stdout, empty stderr.
Diagnostic wording after the prefix may differ by language/OS.
Build the report before writing it: read/analysis failures produce no report.
Write failure may leave partial stdout; stderr itself is best-effort.
Ordinary invalid input must never panic.

Required case matrix: all three fixtures; CRLF; no final LF; blank lines;
single-CR removal; service lengths 1/16/17; `00001`, `0`, `60000`, rejected
`60001` and six digits; negative/decimal durations; missing/extra fields;
tabs/doubled/leading/trailing spaces; mixed-case levels; uppercase services;
non-ASCII/invalid UTF-8/NUL; reverse service insertion order; and direct
near-limit overflow tests. No huge fixture is needed to trigger overflow.
Compare bytes, not strings that normalize newlines. Invoke compiled binaries
for CLI tests, not cargo/go wrappers that may add build output or alter status.

## The 12 assignments

### Week 1: Motivation, tooling, types, and control flow

**Primary reading (45m):** Microsoft's
[ch1 motivation](https://microsoft.github.io/RustTraining/c-cpp-book/ch01-introduction-and-motivation.html)
(intro and safety motivation),
[ch2 getting started](https://microsoft.github.io/RustTraining/c-cpp-book/ch02-getting-started.html)
(main/Cargo),
[ch3 types and variables](https://microsoft.github.io/RustTraining/c-cpp-book/ch03-built-in-types.html),
and [ch4 control flow](https://microsoft.github.io/RustTraining/c-cpp-book/ch04-control-flow.html).
Preview just basic [test assertions](https://doc.rust-lang.org/book/ch11-01-writing-tests.html).
No source exercise set is due.

**Budget:** setup 60m + reading 45m + Starter 30m + Intermediate 90m +
tests 75m + reflection/checks 60m = 360m.

**C++ mental model -> Rust concept:** declarations/statements -> immutable
bindings, shadowing, and value-producing blocks.

**Required Starter:** in the greeting crate, test a function using an if
expression to classify u32 values 0 and 7 as zero/nonzero. Predict and observe
a mutability compiler error; correct it and explain shadowing versus mutation.

**Required Intermediate:** open the supplied `rust/week01-hello` crate and
follow the [greeting spec's TODO order](docs/exercises/rust-hello.md).
The manifest, executable wiring, and tests are provided; implement the learner
functions rather than creating a new crate.

No arguments prints `Hello, learner!` plus LF; one nonempty name prints
`Hello, NAME!` plus LF. Preserve names, including quoted spaces and Unicode.
An empty name or extra arguments means exit 2, empty stdout, and stderr
starting `error:`. Success exits 0 with empty stderr. Factor and test greeting
construction; test default, Ada, Ada Lovelace, and empty-name rejection.
After `cargo build`, manually check extra arguments using the executable
inside the crate: `.\target\debug\week01-hello.exe Ada extra` on Windows or
`./target/debug/week01-hello Ada extra` on macOS. Immediately display
`$LASTEXITCODE` in PowerShell or `printf 'exit=%s\n' "$?"` in Terminal;
the expected code is 2.
No panic for validation. CLI String/Vec and basic tests are small previews;
no custom Result architecture or lifetime annotations are required.

After implementing, run from the course root.

**Windows PowerShell**

```powershell
Push-Location .\rust\week01-hello
cargo run -- Ada
cargo test
cargo fmt --check
Pop-Location
```

**macOS Terminal**

```sh
(
  cd rust/week01-hello &&
  cargo run -- Ada &&
  cargo test &&
  cargo fmt --check
)
```

**Gate/deliverable:** submit crate and week-01 reflection with versions,
test evidence, and one corrected compiler diagnostic. Explain the
semicolon-free final expression before week 2. No broken examples remain.
**Rubric:** Starter 15 + greeting behavior 30 + tests 25 +
tooling/reflection/formatting 30 = 100.
**Optional Challenge:** test a loop over 0..=3 without changing CLI output.
**Reflect:** how does expression-valued control flow differ from your C++
habits, and what does Cargo replace in your build/test workflow?

### Week 2: Collections, structs, enums, and statistics

**Prerequisite:** week 1 gate.
**Primary reading (60m):** Microsoft's
[ch5 structures](https://microsoft.github.io/RustTraining/c-cpp-book/ch05-data-structures.html)
(arrays, slices, String/&str, structs, Vec; skim HashMap) and
[ch6 enums/match](https://microsoft.github.io/RustTraining/c-cpp-book/ch06-enums-and-pattern-matching.html)
(variants, exhaustiveness, methods). Skip long reference deep dives.
**Budget:** reading 60m + Starter 30m + Intermediate 120m + tests 90m +
reflection/checks 60m = 360m.

**C++ mental model -> Rust concept:** vector/span and variant -> Vec/slices
and explicit data-carrying alternatives.
**Required Starter:** classify an empty versus nonempty slice with the supplied
enum carrying the count for the nonempty case. Test both alternatives with
exhaustive match; identify owners versus views for an array, Vec, and slice.

**Required Intermediate:** implement borrowed-slice stats for `&[u32]` with
u64 count/sum and `Option<u32>` minimum/maximum. Empty is valid and has no extrema.
Pure calculation accepts any u32, including u32::MAX. Check addition and return
an error on overflow.

The CLI accepts zero or more arguments, each 1-4 ASCII digits with value
0..1000; leading zeros allowed. Reject signs, decimals, whitespace inside an
argument, empty arguments, and out-of-range values: exit 2, empty stdout,
`error:` stderr. Success exits 0, empty stderr, one LF-terminated line.
`4 0 9 4` prints `count=4 sum=17 min=0 max=9`; no arguments prints
`count=0 sum=0 min=none max=none`. Arithmetic failure exits 1 with `error:`.
Validate all input before printing.

Use the supplied crate in `rust/week02-stats`; do not run `cargo new` over it.
Start with the [shared spec and ordered TODOs](docs/exercises/stats.md).
The scaffold is designed to build, but its active acceptance tests intentionally
fail until you implement the learner functions.

After implementing, run from the course root.

**Windows PowerShell**

```powershell
Push-Location .\rust\week02-stats
cargo run -- 4 0 9 4
cargo test
cargo fmt --check
Pop-Location
```

**macOS Terminal**

```sh
(
  cd rust/week02-stats &&
  cargo run -- 4 0 9 4 &&
  cargo test &&
  cargo fmt --check
)
```

**Small ch9 previews:** Option models absent extrema/checked_add; numeric parse
returns Result. Use tiny explicit matches, not `?` or composed errors yet.
The log parser and full malformed-input matrix are not due.

**Gate/tests/deliverable:** submit crate and week-02 reflection. Test empty,
singleton, duplicates, zero, ascending/descending, u32::MAX, and near-limit
checked-add state. Reject `-1`, `1001`, `1.0`, `""`, `" 1"`, and five digits;
accept `0001`. Show the caller retains its vector. Before week 3, explain
empty versus zero and match every enum variant.
**Rubric:** Starter 15 + statistics/CLI 40 + tests 25 +
reflection/formatting 20 = 100.
**Optional Challenge:** mean as `Option<f64>`, no output changes.
**Reflect:** why take a slice rather than Vec? Compare absence and the dataset
enum with std::optional and std::variant.

### Week 3: Ownership and a typed in-memory record

**Prerequisite:** week 2 gate.
**Primary reading (60m):** Microsoft's
[ch7 ownership](https://microsoft.github.io/RustTraining/c-cpp-book/ch07-ownership-and-borrowing.html)
(moves, borrows, Copy/Clone/Drop) and
[lifetimes](https://microsoft.github.io/RustTraining/c-cpp-book/ch07-1-lifetimes-and-borrowing-deep-dive.html)
(one-input elision, references in structs). Skim only Box versus Rc in
[smart pointers](https://microsoft.github.io/RustTraining/c-cpp-book/ch07-2-smart-pointers-and-interior-mutability.html);
defer graph/interior-mutability exercises.
**Budget:** reading 60m + Starter 30m + Intermediate 120m + tests 90m +
reflection/checks 60m = 360m.

**C++ mental model -> Rust concept:** string_view and moves -> checked
reference lifetimes and unavailable moved-from non-Copy bindings.
**Required Starter:** move a String, copy a u32, clone a String; record a
use-after-move diagnostic and remove the broken code. Show shared then mutable
borrowing after last use. Use a separate test-only Drop marker to set a
borrowed boolean on destruction; check it after dropping the marker.
Do not combine Copy and Drop.

**Required Intermediate:** use the supplied `rust/logscope` library.
Follow the [week 3 TODOs and focused tests](docs/exercises/rust-logscope.md);
later parser, CLI, reader, and worker TODOs are supplied but not due yet.

Define Level (Info/Warn/Error), `Record<'a>` with a borrowed service &str and
u32 duration, and OwnedRecord with String service. Required APIs:
`record(level: Level, service: &str, duration: u32) -> Record<'_>`,
`service_name<'a>(record: &Record<'a>) -> &'a str`, and
`Record::into_owned(self) -> OwnedRecord`. Make fields/accessors testable.
The inputs are already typed; no parser, validation, or Result API is due.
Allocate the service String only when making the owned representation.

Use the spec's week 3 acceptance command and formatting check. Compiling all
test targets should work now; running later-stage tests will fail at their
unfinished TODOs.

**Gate/tests/deliverable:** submit library and week-03 reflection. Construct
Info/api/10, Warn/api/20, Error/worker/100, and zero-duration records.
Test fields, a service pointer identical to the original slice, an owned
conversion surviving the original String's scope, and a borrowed service
surviving a short-lived Record wrapper but not its underlying String.
Explain why lifetime annotations cannot repair a returned borrow of a local.
All ownership tests must pass before week 4.
**Rubric:** Starter experiments 20 + record/lifetime design 35 +
ownership tests 25 + reflection/formatting 20 = 100.
**Optional Challenge:** repeat a small experiment with Box; do not redesign
records around smart pointers.
**Reflect:** what obligation would string_view leave to you? Contrast Rust
moves with C++ standard-library values in a valid-but-unspecified moved state.

### Week 4: Modules, tests, errors, and analyzer integration

**Prerequisite:** week 3 typed-record gate.
**Primary reading (60m):** Microsoft's
[ch8 modules](https://microsoft.github.io/RustTraining/c-cpp-book/ch08-crates-and-modules.html)
(visibility, mod/use, Cargo tests),
[testing patterns](https://microsoft.github.io/RustTraining/c-cpp-book/ch08-1-testing-patterns.html)
(assertions, helpers, integration tests),
[ch9 errors](https://microsoft.github.io/RustTraining/c-cpp-book/ch09-error-handling.html)
(Option/Result, match, `?`), and selected
[error practices](https://microsoft.github.io/RustTraining/c-cpp-book/ch09-1-error-handling-best-practices.html)
(preserve causes; do not default failures). Skip external crates, workspaces,
property/snapshot testing, thiserror, and anyhow.
**Budget:** reading 60m + Starter 30m + Intermediate 150m + tests 90m +
reflection/checks 30m = 360m.

**C++ mental model -> Rust concept:** headers/exceptions/nullable outputs ->
modules, Result, Option, and explicit propagation.
**Required Starter:** add
`parse_line(line: &[u8]) -> Result<Record<'_>, ParseError>` using week 3's
model. Accept a normalized nonempty line, not a complete file; no LF splitting
or CR stripping inside parse_line. Empty is a parse error. Validate bytes
before converting the borrowed service. First test `INFO api 10`, a missing
duration, and `INFO api -1`; explain which caller counts these failures.

**Required Intermediate:** implement the supplied parser/summary/CLI TODOs in
the existing crate, following the [week 4 spec](docs/exercises/rust-logscope.md).
The module declarations and main entry point are already wired, re-exporting
`analyze(input: &[u8]) -> Result<Summary, AnalyzeError>` and
`render(summary: &Summary) -> String` from lib.rs. Own summary service keys.
Use the provided fixtures, implement the complete shared contract, and keep
main as orchestration. Use at least one match and one appropriate `?`.
Plain loops suffice; the closure/iterator refactor waits for week 5.
Whole-file byte reads are acceptable. Handle stdout failures explicitly.

Run the spec's week 3-4 focused acceptance commands, fmt, and clippy. Reader and
worker tests are due in later weeks; a full unfiltered run is not this gate.

**Gate/tests/deliverable:** complete the shared matrix now, including invalid
UTF-8 and overflow. Compare rendered bytes to all three expected files.
Read and run the supplied executable tests using std::process::Command and Cargo's
`env!("CARGO_BIN_EXE_logscope")`, which includes the correct OS executable
suffix. Construct fixture paths with Path::join; test success, wrong count, missing file,
and a unique temporary path with spaces, cleaned up afterward.
Submit code and week-04 reflection. Fixtures, public integration tests,
and fatal error behavior must pass before refactoring; freeze the contract
for Go. Test overflow through near-limit state, not huge input.
**Rubric:** Starter/parser 20 + aggregation/render/CLI 35 +
test matrix 30 + modules/reflection/formatting 15 = 100.
**Optional Challenge:** public API documentation examples, std only.
**Reflect:** why own keys when parsing borrows, and why sort map keys?

### Week 5: Traits, conversions, closures, and a reader boundary

**Prerequisite:** week 4 gate.
**Primary reading (60m):** Microsoft's
[ch10 traits](https://microsoft.github.io/RustTraining/c-cpp-book/ch10-traits.html)
(impl, static/dynamic dispatch),
[generics](https://microsoft.github.io/RustTraining/c-cpp-book/ch10-1-generics.html)
(bounds),
[ch11 conversions](https://microsoft.github.io/RustTraining/c-cpp-book/ch11-from-and-into-traits.html)
(From/Into/Default, fallible versus infallible), and
[ch12 closures/iterators](https://microsoft.github.io/RustTraining/c-cpp-book/ch12-closures.html)
(capture basics, map/collect). Use
[iterator tools](https://microsoft.github.io/RustTraining/c-cpp-book/ch12-1-iterator-power-tools.html)
only as a lookup. Skip operator overloads, type-state/builders, custom
iterators, and boxed closure collections. Look up std::io::Read and Cursor.
**Budget:** reading 60m + Starter 30m + Intermediate 105m + tests 105m +
reflection/checks 60m = 360m.

**C++ mental model -> Rust concept:** constrained templates, conversions, and
lambdas -> trait bounds, From/Into, and borrow-checked closure capture.
Use the [week 5 TODO and test map](docs/exercises/rust-logscope.md).
**Required Starter:** Default for Summary means the empty report. Implement
`From<io::Error>` for AnalyzeError, retaining its cause; test both From and
Into. A failing input does not become valid through Default.
**Required Intermediate:** implement the supplied
`analyze_reader<R: std::io::Read>(reader: R) -> Result<Summary, AnalyzeError>`.
Read all bytes then reuse analyze; use your conversion and `?`. The CLI opens
File and uses this function. Short reads must not truncate input. Preserve
all public APIs and output; this is testability, not a streaming-memory bound.
Refactor one renderer collection loop to a closure with map/collect, keeping
sorting explicit and checked arithmetic intact.

**Gate/tests/deliverable:** test Cursor on empty/sample input, immediate read
failure, partial valid bytes followed by an error, and one-byte reads.
No partial report on read failure. Inspect the original I/O error without
platform-specific wording assertions. Keep every earlier regression test.
Submit updated crate and week-05 reflection with a before/after call graph.
Run fmt, the spec's week 3-5 acceptance commands, and clippy inside rust/logscope, then
`cargo run -- ..\..\fixtures\malformed.log` in Windows PowerShell or
`cargo run -- ../../fixtures/malformed.log` in macOS Terminal.
Gate: File and Cursor both work, failures propagate, fixture bytes are unchanged.
**Rubric:** Starter/conversions 20 + generic reader/iterator refactor 30 +
regression/failure tests 30 + explanation/formatting 20 = 100.
**Optional Challenge:** use `&mut dyn Read` with the API and explain dispatch.
**Reflect:** compare the bound to templates/virtual interfaces; explain File
ownership, iter versus into_iter, and closure captures.

### Week 6: Safe, bounded Rust concurrency

**Prerequisite:** week 5 gate.
**Primary reading (45m):** Microsoft's
[ch13 concurrency](https://microsoft.github.io/RustTraining/c-cpp-book/ch13-concurrency.html)
(spawn, scope, move, joins, Send/Sync; skim Arc/Mutex), then
[`thread::scope`](https://doc.rust-lang.org/std/thread/fn.scope.html)
and [scoped joins](https://doc.rust-lang.org/std/thread/struct.ScopedJoinHandle.html).
Skip the source's one-thread-per-line challenge and extra lock/channel tasks.
**Budget:** reading 45m + Starter 30m + design 30m + Intermediate 105m +
tests 105m + reflection/checks 45m = 360m.

**C++ mental model -> Rust concept:** thread lifetime obligations -> scoped
borrowing, move captures, and Send/Sync constraints.
Use the [week 6 TODO and test map](docs/exercises/rust-logscope.md).
**Required Starter:** two scoped threads process fixed complete in-memory
records with private counts. Join both and assert the combined count.
Explain why borrowing the local input is valid.
**Required Intermediate:** implement the supplied
`analyze_parallel(input: &[u8], workers: usize) -> Result<Summary, AnalyzeError>`.

Add a library-only parallel analyzer accepting 1-4 workers. Partition complete
records, give each worker its own summary, and merge with checked arithmetic.
Reject invalid worker counts. Join every worker and surface worker failures.
Compare results with the sequential implementation, including empty input and
fewer records than workers. Keep the existing CLI unchanged.

At most workers scoped threads, no thread per line or shared mutable map.
Join all handles even after Err or panic; convert a joined panic into an error.
Do not promise recovery with panic=abort. Never silently omit worker output.
**Gate/tests/deliverable:** fixtures, counts 1..4, zero-byte input, CRLF,
unterminated last line, mixed malformed lines, and fewer lines than workers
must equal sequential bytes. Test invalid counts, checked-merge overflow,
and a test-only join helper with Err/panic workers. Submit code and week-06
reflection with an ownership/join diagram. Run all Rust checks.
Gate: bounded creation, every handle joined, errors surfaced, equivalence
tests pass; explain Send versus Sync before Go.
**Rubric:** Starter 15 + bounded ownership/partitioning 25 +
checked merge/failure handling 25 + tests 20 + reflection/formatting 15 = 100.
**Optional Challenge:** time both paths without a speed target.
**Reflect:** what keeps the borrowed buffer alive, and how does this differ
from std::thread lifetime and exception obligations?

### Week 7: Go values, slices, and maps

**Prerequisite:** Rust gates completed; Go installed.
**Read (60m):** [Tour basics](https://go.dev/tour/basics/1) (first 8 pages),
[pointers/structs/slices](https://go.dev/tour/moretypes/1) (pages 1-16),
and [maps](https://go.dev/tour/moretypes/19) (pages 19-22).
**Budget:** reading 60m + setup 30m + implementation 105m +
tests/aliasing lab 105m + reflection/checks 60m = 360m.

Use the supplied module, main, and tests in `go/stats`; do not run `go mod init`
over it. Follow the [shared spec and ordered TODOs](docs/exercises/stats.md),
including the aliasing lab. Tests intentionally fail until implementation.
Go remains a week 7 exercise after the Rust gates. Then run the checks from
the course root.

**Windows PowerShell**

```powershell
Push-Location .\go\stats
go run . 4 0 9 4
go test ./...
gofmt -l .
go vet ./...
Pop-Location
```

**macOS Terminal**

```sh
(
  cd go/stats &&
  go run . 4 0 9 4 &&
  go test ./... &&
  gofmt -l . &&
  go vet ./...
)
```

Reimplement week 2's exact numeric grammar, output, and exit behavior using
`[]uint32` and uint64 count/sum. This unsigned specification matches the
revised Rust stats assignment; signs remain invalid at the CLI. The pure
function returns `(Stats, error)`, with explicit presence for absent extrema.
Use table-driven tests, including all week 2 cases and uint32's maximum.
No arguments is an empty dataset, not a usage failure.

Add tests for independent scalar struct copies, pointer mutation, shared slice
elements, an explicitly independent slice copy, shared map entries, and
missing versus present-zero map keys. Test append with deliberately sufficient
capacity; do not assume a growth policy. Never intentionally panic in normal
input handling.

**Acceptance/deliverable:** stats parity and aliasing tests pass, formatting
listing is empty, vet passes; submit module and week-07 reflection.
**Rubric:** stats/CLI parity 35 + aliasing lab 25 + edge-case tests 25 +
reflection/formatting 15 = 100.
**Optional Challenge:** a copied struct containing a slice; explain aliasing.

**Reflect:** Why are a Go slice, a C++ vector, and a Rust vector not equivalent?

### Week 8: The same analyzer in Go

**Prerequisite:** week 7 plus Rust's week 4 frozen analyzer contract.
**Read (45m):** [module layout](https://go.dev/doc/modules/layout), Go's
[error-handling guidance](https://go.dev/blog/error-handling-and-go)
and [os.ReadFile](https://pkg.go.dev/os#ReadFile).
**Budget:** reading 45m + design 30m + implementation 150m +
tests 90m + reflection/checks 45m = 360m.

Use the supplied `go/logscope` and implement your frozen week 4 contract in Go.
Use explicit errors for file and analysis failures, and count malformed records
according to the contract. Match deterministic output exactly. Reuse the same
inputs and expected results, not the Rust implementation's source structure.

The module, command entry points, and tests are supplied. Start with the
[Go week 8 TODO map and focused checks](docs/exercises/go-logscope.md);
do not rerun module initialization. Later interface/worker/HTTP TODOs are
present but not part of this gate.

Library package logscope lives at the module root, main at
cmd/logscope/main.go, importing by module path. Required APIs:
`ParseLine(line []byte) (Record, error)`,
`Analyze(input []byte) (Summary, error)`, and `Render(summary Summary) string`.
ParseLine takes normalized nonempty lines. Use a Level type/constants and
owned service strings. Use uint64 counters/sums, checked arithmetic, and
byte-based parsing. Whole-file ReadFile is acceptable. Do not use Unicode
Fields/TrimSpace or silently truncate through Scanner's default token limit.
Keep os.Exit at outer main after work returns.

**Acceptance/tests/deliverable:** complete the shared matrix, all fixture byte
outputs, CRLF, zero bytes, missing final LF, invalid UTF-8, and reverse service
insertion order; test near-overflow updates. Manually build/invoke the CLI
to check exits using [these OS-specific commands](SETUP.md#executables-exit-status-and-fixtures);
week 9 automates it. Submit module and week-08 reflection.
Inside go/logscope run the spec's week 8 checks and
`go run .\cmd\logscope ..\..\fixtures\sample.log` in Windows PowerShell or
`go run ./cmd/logscope ../../fixtures/sample.log` in macOS Terminal.
**Rubric:** grammar/aggregation 35 + CLI/render parity 30 +
test coverage 25 + explanation/formatting 10 = 100.
**Optional Challenge:** wrap an error with `%w` while preserving its identity.

**Reflect:** Which errors should carry context, and where should you print them?

### Week 9: Go interfaces and resource cleanup

**Prerequisite:** week 8.
**Read (45m):** [interfaces](https://go.dev/tour/methods/9) (pages 9-12),
[defer](https://go.dev/tour/flowcontrol/12) (pages 12-13), and
[io.ReadAll](https://pkg.go.dev/io#ReadAll).
**Budget:** reading 45m + design 30m + implementation 105m +
tests 135m + reflection/checks 45m = 360m.

Refactor around standard-library reader/writer interfaces where useful.
Use the supplied [week 9 TODOs and tests](docs/exercises/go-logscope.md).
Test short reads, read failures, and write failures with small test doubles.
When opening files explicitly, ensure they close on all return paths.
Preserve CLI behavior and keep interfaces small and consumer-driven.

Required APIs: `AnalyzeReader(r io.Reader) (Summary, error)` reads all bytes
then reuses Analyze; partial bytes plus a read error must not yield a report.
`AnalyzeFile(path string) (Summary, error)` opens, defers Close, and is used
by the CLI. For read-only handles, return Close error if analysis succeeded;
otherwise retain the earlier error. Use an unexported io.ReadCloser helper
for injected close tests. Do not let os.Exit bypass defers.

**Acceptance/tests/deliverable:** table-test immediate/partial read failures,
one-byte reads, Close success/failure, and simultaneous read/close failure.
Verify Close once per opened handle. Keep writer failure tests at your command
boundary. Automate CLI tests with os/exec: build once into t.TempDir, fail
on build failure, and test success, zero/multiple arguments, a missing file,
and a temporary path with spaces. No shell or exact OS error-text assertion.
Use filepath.Join for filesystem paths and add the executable's .exe suffix
only when runtime.GOOS is windows; macOS uses an unsuffixed name.
Submit module and week-09 reflection; run formatting/vet and the spec's week 8-9 tests.
**Rubric:** reader/interface boundary 25 + resource/error policy 25 +
table/process tests 35 + reflection/formatting 15 = 100.
**Optional Challenge:** demonstrate a typed-nil interface in a focused test.

**Reflect:** Which resource lifetimes still require explicit action despite GC?

### Week 10: Go workers and cancellation

**Prerequisite:** week 9.
**Read (60m):** [Tour concurrency](https://go.dev/tour/concurrency/1) (pages 1-6),
[`context`](https://pkg.go.dev/context), and the
[race detector guide](https://go.dev/doc/articles/race_detector.html).
**Budget:** reading 60m + design 45m + implementation 120m +
tests 90m + reflection/checks 45m = 360m.

Implement the supplied library-only analyzer with 1-4 workers. Follow the
[week 10 TODO/test map](docs/exercises/go-logscope.md). Give workers private summaries,
bound queued work, and use context cancellation. On failure or cancellation,
stop sending work, wait for every worker, and return an error rather than a
partial success. Test invalid counts, cancellation, worker errors, and
equivalence with sequential results. Use synchronization, not sleeps, in tests.

Required API:
`AnalyzeParallel(ctx context.Context, input []byte, workers int) (Summary, error)`.
Validate worker count before checking a non-nil context. Use complete normalized
lines. The caller-side producer alone closes jobs; buffer jobs to at most
workers. Start exactly workers goroutines, each sending one final result into
a result channel buffered to workers. No shared mutable map.
Check cancellation on job send/receive and between records. On worker error,
cancel a derived context, collect every worker result, then return the original
failure, not internal context.Canceled. Checked merging follows joining.
Absent another error, return the caller's ctx.Err if canceled, checking once
more before success. The sequential CLI stays unchanged.

**Acceptance/tests/deliverable:** fixtures and 1..4 workers equal sequential
bytes; test zero-byte input, CRLF, fewer lines than workers, invalid counts,
pre-canceled context, and checked-merge overflow. Use an unexported testable
worker helper/callback to signal entry, block on ctx.Done, and demonstrate
in-flight cancellation and every worker exiting before return. Inject a worker
error through the same seam. Timeouts guard deadlock; sleeps do not coordinate.
Submit module and week-10 reflection; run the spec's week 8-10 checks and race
tests. HTTP tests become required in week 11.
**Rubric:** bounded ownership 30 + cancellation/error/join correctness 30 +
equivalence/lifecycle tests 25 + reflection/race evidence 15 = 100.
**Optional Challenge:** measure overhead without a speed requirement.

Run race tests on a supported toolchain. On Windows they require a compatible
C compiler in addition to Go. See the
[Windows/macOS race prerequisites](SETUP.md#race-detector-prerequisites).
Ordinary tests passing is not a race-clean result.

**Reflect:** Who closes each channel, and what proves every worker has stopped?

### Week 11: A small local Go HTTP server

**Prerequisite:** weeks 8-10.
**Read (45m):** [net/http](https://pkg.go.dev/net/http) (Handler, Server, Shutdown),
[httptest](https://pkg.go.dev/net/http/httptest), and
[signal.NotifyContext](https://pkg.go.dev/os/signal#NotifyContext).
**Budget:** reading 45m + design 30m + implementation 135m +
tests 105m + reflection/checks 45m = 360m.

Implement the supplied separate server command using the
[week 11 TODO/test map](docs/exercises/go-logscope.md).
Read and analyze one file at startup, then serve
the immutable report at GET `/summary`. Bind only to `127.0.0.1:8080`.
Define and test unknown-path and unsupported-method behavior. Use `httptest`
without depending on an available port. Handle Ctrl+C, use a bounded shutdown
timeout, and wait for the serving goroutine. No uploads or deployment.

Command: cmd/logserve/main.go, accepting exactly one file path. Wrong count
exits 2; file/analysis/listen failures exit 1 with `error:` stderr.
Do not listen on a bad file. GET `/summary` returns 200, Content-Type
`text/plain; charset=utf-8`, and exact report bytes. Other methods there,
including HEAD, return 405, `Allow: GET`, empty body. Other paths, including
`/summary/`, return 404 with empty body, taking precedence over method checks.
Query strings do not change path behavior. Do not mutate captured report data.

ReadHeaderTimeout is five seconds. Use a dedicated handler/mux. Ctrl+C uses
os.Interrupt and NotifyContext; Shutdown uses a fresh two-second timeout,
then Close if needed, and joins the serving goroutine. ErrServerClosed is
normal; successful shutdown exits 0.

**Acceptance/tests/deliverable:** httptest covers every endpoint case,
headers, empty bodies, and exact sample/empty report bytes. Run concurrent
handler tests without mutations, normal/race checks, and record manual startup
and shutdown. Submit module and week-11 reflection.

**Windows PowerShell**

```powershell
Push-Location .\go\logscope
New-Item -ItemType Directory -Force .\bin
go build -o .\bin\logserve.exe .\cmd\logserve
if ($LASTEXITCODE -ne 0) { Pop-Location; throw 'Build failed.' }
.\bin\logserve.exe ..\..\fixtures\sample.log
$LASTEXITCODE
Pop-Location
```

**macOS Terminal**

```sh
(
  cd go/logscope &&
  mkdir -p bin &&
  go build -o bin/logserve ./cmd/logserve || exit
  ./bin/logserve ../../fixtures/sample.log
  printf 'exit=%s\n' "$?"
)
```

It blocks. In another terminal, from any directory, use
`curl.exe -i http://127.0.0.1:8080/summary` on Windows or
`curl -i http://127.0.0.1:8080/summary` on macOS. Then press Ctrl+C in the first
terminal; its next command displays the server's exit status, expected 0
on successful shutdown.
If the port is occupied, do not kill unrelated processes; retry when free.
Tests must not need that port.
**Rubric:** endpoint contract 30 + startup/shutdown/resource handling 25 +
handler tests 30 + reflection/formatting 15 = 100.
**Optional Challenge:** test on an ephemeral loopback port.

**Reflect:** Why is immutable captured data easier to share between requests?

### Week 12: Cross-language capstone

**Prerequisite:** weeks 1-11.
**Read (30m):** revisit your Microsoft ch7-13 notes and the shared contract.
Consult [testing](https://pkg.go.dev/testing) for gaps only.
**Budget:** review/reading 30m + parity harness 90m + small fixes 90m +
final checks 60m + comparison/submission 90m = 360m.

No new product. Trace and run the supplied `go/logscope/parity_test.go` using
the [capstone spec](docs/exercises/parity.md). It builds both sequential CLIs
into temporary locations and compares bytes and exit status against expected
results, not only against each other. Repair implementation gaps and rerun
checks in every project.

The harness uses filepath.Join, native executable suffixes, argument arrays,
timeouts, and an explicit temporary Cargo target directory. Build failures
fail tests. It covers every fixture, no/two arguments, missing files, filenames
with spaces or leading dashes, CRLF, and invalid bytes with independent
expectations. Both toolchains must be on PATH. The explicit `parity` build tag
keeps this dependency out of earlier Go gates; run
`go test -tags parity -count=1 -run '^TestParity$' .` inside `go/logscope` in
week 12. Untagged tests alone are not capstone acceptance.

Write `docs/comparison.md` (700-1000 words) using concrete examples from your
work: ownership and aliasing, error propagation, resource cleanup, interfaces,
testing, and concurrency. Explain which language you would choose for this
project and why. Separate measured evidence from preference.

Within that existing writing budget, briefly connect one idea from
[Microsoft ch16 case studies](https://microsoft.github.io/RustTraining/c-cpp-book/ch16-case-studies.html)
to your design, such as a closed Level enum rather than a class hierarchy.
No additional case-study implementation. Explain the sequence: typed borrowed
records in week 3, Result parsing/modules in week 4, and trait/conversion
refactoring in week 5.

**Acceptance/deliverable:** source/tests, week-12 reflection, comparison,
and full regression evidence. Preserve both sequential CLIs, parallel
extensions, and the HTTP command. Mark outstanding gates/checks honestly.
**Rubric:** byte/status parity harness 35 + full regression evidence 25 +
comparison with concrete examples 30 + clean navigable portfolio 10 = 100.
**Optional Challenge:** peer demonstration and one documentation improvement.

**Reflect:** Which C++ habits did you stop applying, and which remain valuable?

## Checks before submitting

For an intermediate logscope week, use the cumulative test selectors in its
spec so later-week TODOs do not block the current gate. The commands below are
the full-project checks once all stages are due.

Inside each Rust crate, in either Windows PowerShell or macOS Terminal:

```text
cargo fmt --check
cargo test
cargo clippy --all-targets -- -D warnings
```

Install the `rustfmt` and `clippy` components through `rustup` if missing.
Inside each Go module, in either shell:

```text
gofmt -l .
go vet ./...
go test ./...
```

For week 12, also run
`go test -tags parity -count=1 -run '^TestParity$' .` inside `go/logscope`.
This explicitly enables the cross-language harness and requires both toolchains.

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

When reading simplified comparisons, preserve precise semantics: using a
moved-from C++ standard-library object is not automatically undefined behavior;
operations retain their preconditions. Rust's unavailable moved binding is a
different rule. Lifetime annotations do not extend storage duration; Clone is
not always a deep copy; Rc/Arc cycles can leak, and Arc alone does not synchronize
mutation. Ignoring Result normally warns rather than unconditionally failing
compilation. `chunks_exact` exposes a remainder rather than panicking merely
because one exists. Use stable std documentation to resolve API details.

## Optional post-course roadmap

Outside the 72-hour required plan, choose only what serves a concrete need:

- [ch14 unsafe/FFI](https://microsoft.github.io/RustTraining/c-cpp-book/ch14-unsafe-rust-and-ffi.html):
  raw pointers and safe-wrapper obligations; no required unsafe implementation.
- [ch15 no_std](https://microsoft.github.io/RustTraining/c-cpp-book/ch15-no_std-rust-without-the-standard-library.html):
  core/alloc and embedded constraints.
- [ch16 case studies](https://microsoft.github.io/RustTraining/c-cpp-book/ch16-case-studies.html):
  enum dispatch and arena-index ownership; only a brief week 12 discussion
  belongs to the graded course.
- [ch17 practices](https://microsoft.github.io/RustTraining/c-cpp-book/ch17-best-practices.html):
  a checklist, not a mandate to add dependencies.
- [ch18 C++ semantic deep dives](https://microsoft.github.io/RustTraining/c-cpp-book/ch18-cpp-rust-semantic-deep-dives.html):
  casts, preprocessing, and concepts.
- [ch19 macros](https://microsoft.github.io/RustTraining/c-cpp-book/ch19-macros.html):
  authoring macros is deferred; using println/assert_eq/derive earlier is fine.

The source is a reading resource, not an instruction to install its example
crates or replace this course's explicit failure policies with default values.
