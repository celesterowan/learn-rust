# Week 4: Modules, tests, errors, and analyzer integration

**Prerequisite:** week 3 typed-record gate. **Objectives:** make errors values,
organize modules, and integrate a deterministic byte parser and analyzer.

**Primary reading (60m):** Microsoft's
[ch8 crates/modules](https://microsoft.github.io/RustTraining/c-cpp-book/ch08-crates-and-modules.html)
(visibility, mod/use, Cargo tests),
[testing patterns](https://microsoft.github.io/RustTraining/c-cpp-book/ch08-1-testing-patterns.html)
(built-in assertions, helpers, integration tests), and
[ch9 errors](https://microsoft.github.io/RustTraining/c-cpp-book/ch09-error-handling.html)
(Option/Result, match, `?`). Consult
[error practices](https://microsoft.github.io/RustTraining/c-cpp-book/ch09-1-error-handling-best-practices.html)
only for preserving causes and avoiding inappropriate fallbacks.
Skip workspaces, dependencies, property/snapshot crates, and thiserror/anyhow.

**Six-hour budget:** reading 60m + Starter 30m + Intermediate 150m +
tests 90m + reflection/checks 30m = 360m.

## C++ mental model -> Rust concept

Headers, exceptions, and nullable outputs -> explicit module visibility,
Result, Option, and propagation. Defaulting a failed duration to zero would
violate the log contract even though it avoids a panic.

## Required Starter

Read the [contract](../docs/logscope-contract.md). Implement the supplied W4-1
`parse_line(line: &[u8]) -> Result<Record<'_>, ParseError>`, reusing week 3's
record model. It accepts one normalized nonempty line; empty input is an
error. No LF splitting or CR stripping occurs in this function. Validate ASCII
before converting service to &str; borrow the original bytes.
First test `INFO api 10`, a missing duration, and `INFO api -1`.
Explain which caller counts a parse error and which errors abort analysis.

## Required Intermediate project

Extend the supplied `rust/logscope` crate without replacing the library.
The [week 4 spec and ordered TODOs](../docs/exercises/rust-logscope.md#week-4-parser-checked-summary-sequential-command)
map each function to a file and focused test. Modules, types, executable
wiring, and active acceptance tests are supplied; do not recreate them.
Build and test compilation must pass while TODO acceptance tests fail.
Required library APIs:

- `analyze(input: &[u8]) -> Result<Summary, AnalyzeError>`
- `render(summary: &Summary) -> String`

Implement the supplied Summary helpers and use the typed AnalyzeError.
Retain owned service keys in Summary;
do not make the entire result borrow the input. Follow the shared contract
for byte-line normalization, invalid records, checked counters, sorting, and
output. Record parsing and summary operations already have distinct declared
modules and re-exports from lib.rs. Use straightforward loops for now:
chapter 12's closure/iterator refactor is week 5. Keep main as orchestration,
not a second parser. Preserve the supplied public signatures and error variants.
Use at least one explicit match and one `?` for appropriate error propagation.

Read the file with byte-oriented file I/O. Handle command arguments, errors,
and exit statuses exactly as specified. Build a complete report before
writing stdout and handle write errors.

From the course root, run only your OS's block.

**Windows PowerShell**

```powershell
Push-Location .\rust\logscope
cargo run -- ..\..\fixtures\sample.log
cargo test --test week03 --test week04_parser --test week04_summary --test week04_cli
cargo fmt --check
cargo clippy --all-targets -- -D warnings
Pop-Location
```

**macOS Terminal**

```sh
(
  cd rust/logscope &&
  cargo run -- ../../fixtures/sample.log &&
  cargo test --test week03 --test week04_parser --test week04_summary --test week04_cli &&
  cargo fmt --check &&
  cargo clippy --all-targets -- -D warnings
)
```

## Gate, tests, and deliverables

Complete the shared case matrix, including a truly zero-byte input in memory.
This is when the full malformed-record matrix becomes required, not week 3.
Compare render output bytes against all three expected-output files. The supplied
`tests/week04_cli.rs` uses `std::process::Command` and Cargo's
`env!("CARGO_BIN_EXE_logscope")` integration-test binary path, which supplies
the correct executable name on either OS. Use Path::join for fixture paths.
Assert success stdout/stderr/status, wrong argument count, missing file,
and a temporary filename containing spaces. Clean up temporary files; choose
unique paths instead of a fixed global filename. No shell is needed in tests.

The supplied helper tests construct near-limit summary state to check overflow.
Failed analysis must not become a success report.
No partial-state recovery API is required.

Submit updated crate and `docs/submissions/week-04.md`.
**Gate:** fixtures match byte-for-byte; invalid records are counted, fatal
errors propagate without a report; public integration tests pass. Freeze this
contract for Go and repair failures before the week 5 refactor.
**Rubric:** Starter/parser 20 + aggregation/render/CLI 35 +
test matrix 30 + modules/reflection/formatting 15 = 100.

**C++ questions:** (1) Why own keys when the parser borrows them?
(2) Which orderings does your map guarantee, and which does the renderer add?

**Optional Challenge (ungraded):** add public-API documentation examples using
only std. Do not add source-guide dependencies or a workspace.
