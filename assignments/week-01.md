# Week 1: Motivation, tooling, types, and control flow

**Prerequisite:** C++ experience and [setup](../SETUP.md).
**Objectives:** navigate Cargo, distinguish mutability from shadowing, and
use expressions to produce values.

**Primary reading (45m):** Microsoft's
[ch1 motivation](https://microsoft.github.io/RustTraining/c-cpp-book/ch01-introduction-and-motivation.html)
(intro and safety motivation),
[ch2 getting started](https://microsoft.github.io/RustTraining/c-cpp-book/ch02-getting-started.html)
(main and Cargo),
[ch3 types/variables](https://microsoft.github.io/RustTraining/c-cpp-book/ch03-built-in-types.html),
and [ch4 control flow](https://microsoft.github.io/RustTraining/c-cpp-book/ch04-control-flow.html).
Read selected concepts, not every diagram. Preview only the basic assertion
example in [Rust testing](https://doc.rust-lang.org/book/ch11-01-writing-tests.html).

**Six-hour budget:** setup 60m + reading 45m + Starter 30m + Intermediate 90m +
tests 75m + reflection/submission 60m = 360m.

## C++ mental model -> Rust concept

C++ declarations and statements -> inferred or annotated types, immutable
bindings, shadowing, and value-producing blocks. Cargo handles the small
project's build and test runner; it is not a request to recreate CMake.

## Required Starter

In the supplied greeting crate, implement the H1 classifier using an if
expression; its tests assert zero/nonzero results for 0 and 7.
Predict a mutability error, compile to observe it, then fix it. Record the
difference between changing a mutable binding and shadowing one.

## Required Intermediate project

Use the supplied `rust/week01-hello` crate; do not run `cargo new` over it.
Follow the [file map and ordered H1-H3 TODOs](../docs/exercises/rust-hello.md).
Build and test compilation must pass before implementation; active acceptance
tests intentionally fail until the learner functions and experiment are done.

Required behavior: with no argument,
print `Hello, learner!` plus LF; with one nonempty name argument, print
`Hello, NAME!` plus LF. Preserve the supplied name, including spaces; do not
trim it. Reject an empty name or more than one argument with exit 2, empty
stdout, and stderr starting `error:`. Successful runs exit 0 with empty stderr.
Names can be ordinary Unicode text; this is not the later ASCII log grammar.

Factor greeting construction into a function and test it. Required cases:
default greeting, `Ada`, `"Ada Lovelace"`, and empty-name rejection.
Do not use panic for input validation. Function signatures, a small error
type, and executable wiring are supplied; implement the greeting body.

After implementing, run the checks from the course root.

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

Manually verify the wrong-argument-count exit after building with `cargo build`.
From inside the crate, run `.\target\debug\week01-hello.exe Ada extra` on
Windows or `./target/debug/week01-hello Ada extra` on macOS. Immediately inspect
`$LASTEXITCODE` in PowerShell or `printf 'exit=%s\n' "$?"` in Terminal.
The expected code is 2. See [exit-status guidance](../SETUP.md#executables-exit-status-and-fixtures).

**Small previews:** collecting CLI arguments introduces String/Vec before
chapter 5. An argument lookup may introduce Some/None; a plain branch is
enough for validation. Basic `#[test]` is previewed before chapter 8.
No lifetime annotations, `?`, or custom error architecture is required.

## Gate, acceptance, and deliverables

Submit the crate and `docs/submissions/week-01.md`. Record versions, successful
test output, and one deliberate compiler error plus its correction. Do not
leave uncompilable code in the crate. The greeting must be produced by your
function, not separately hard-coded in main and tests.

**Gate:** before week 2, produce both greetings and explain why the final
expression in a block has no semicolon. Repair failed required checks first.

**Rubric:** Starter 15 + greeting behavior 30 + tests 25 +
tooling/reflection/formatting 30 = 100.

**C++ questions:** (1) How does a Rust expression return a value without
`return`? (2) What did the compiler reject before execution that your C++
experience helped explain?

**Optional Challenge (ungraded):** classify 0..=3 with a loop in a test,
without changing CLI output. Go installation can wait until week 7.
