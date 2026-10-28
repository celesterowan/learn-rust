# Rust week 1: greeting scaffold

Start after [setup](../../SETUP.md), alongside [week 1](../../assignments/week-01.md).
The handout retains the reading budget, grading, reflection, and gate.
This is a learner project, not a solution. Do not run `cargo new` over it.
There is no root workspace and no dependency to download.

## Files and supplied boundaries

Package and binary: `week01-hello`; Rust library import: `week01_hello`.
All paths below are relative to `rust/week01-hello`.

| File | Responsibility |
| --- | --- |
| `Cargo.toml` | Independent stdlib-only Cargo package |
| `src/lib.rs` | H1 classifier and H2 greeting TODOs; supplied error type/display |
| `src/main.rs` | Supplied argument, stream, and process-exit adapter |
| `tests/hello.rs` | Active classifier/function/process tests and H3 experiment |

`todo!` in a non-Result exercise is an explicit unfinished panic, not a model
for validating user input. H2 instead returns `HelloError::Unimplemented`.
Remove these failures by implementing the functions, not by deleting tests,
ignoring tests, catching TODO panics, or changing the expected outputs.

## Ordered work

1. **H1, Starter:** implement `classify(value: u32) -> &'static str` using
   an if expression. Return exactly `zero` for 0 and `nonzero` otherwise.
   Run `cargo test --test hello starter_classification`.
2. **H2, Intermediate:** implement
   `greeting(arguments: &[String]) -> Result<String, HelloError>`.
   Arguments exclude the executable name. Zero arguments returns
   `Hello, learner!\n`. Exactly one nonempty name returns `Hello, NAME!\n`.
   Preserve the name, including spaces, Unicode, and whitespace-only names;
   never trim it. An empty name or extra arguments returns
   `HelloError::InvalidArguments`. Do not print or panic in this function.
   Run `cargo test --test hello greeting_contract` and
   `cargo test --test hello invalid_arguments_are_not_panics`.
3. **H3, Starter experiment:** replace only the TODO body of
   `starter_mutability_and_shadowing` with a small passing test that mutates
   a `mut` binding, shadows a binding, and asserts both results. First predict
   and temporarily compile an attempted immutable reassignment. Save the
   diagnostic and correction in `docs/submissions/week-01.md`, then remove the
   invalid code. Run `cargo test --test hello starter_mutability_and_shadowing`.
4. Run all checks below and complete the handout reflection. Explain the final
   expression without a semicolon. Do not redesign the supplied adapter.

The supplied process maps invalid arguments to exit 2 with empty stdout and
`error:` stderr, success to exit 0 with exact report bytes and empty stderr,
and unfinished/internal output errors to exit 1. It writes bytes, so LF is
unchanged on Windows. CLI tests use the compiled Cargo binary directly.
`String`, `Vec`, `Result`, and the tiny supplied error type are previews:
use explicit matches; `?`, lifetime annotations, and error architecture are
not this week's work. Ordinary input validation must not panic.

## Commands and expected results

Start at the course root. Choose one OS block. Initially **build and test
compilation must pass**, but **acceptance tests intentionally fail** at TODOs.
Do not interpret the last command's success as proof of earlier success.
After implementation every command must pass.

**Windows PowerShell**

```powershell
Push-Location .\rust\week01-hello
cargo build
cargo test --no-run
cargo test --test hello
cargo fmt --check
cargo run -- Ada
Pop-Location
```

**macOS Terminal**

```sh
(
  cd rust/week01-hello &&
  cargo build &&
  cargo test --no-run &&
  cargo test --test hello &&
  cargo fmt --check &&
  cargo run -- Ada
)
```

The all-target final acceptance command inside the crate is `cargo test`.
Tests cover 0/7/u32::MAX classification; default, Ada, quoted spaces,
leading/trailing/whitespace-only names, Unicode; empty/extra argument errors;
and exact compiled-process stdout/stderr/status.
Also perform the handout's direct executable exit-code check. No fixture is
needed for greetings. Keep your reflection outside the source crate and do
not leave the deliberately uncompilable experiment checked in.
