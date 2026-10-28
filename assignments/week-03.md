# Week 3: Ownership and a typed in-memory record

**Prerequisite:** week 2 gate. **Objectives:** distinguish move, Copy, Clone,
and Drop; tie a borrowed field to its owner's lifetime.

**Primary reading (60m):** Microsoft's
[ch7 ownership](https://microsoft.github.io/RustTraining/c-cpp-book/ch07-ownership-and-borrowing.html)
(moves, borrows, Copy/Clone/Drop), and
[lifetimes](https://microsoft.github.io/RustTraining/c-cpp-book/ch07-1-lifetimes-and-borrowing-deep-dive.html)
(one-input elision and references in structs). Skim the overview of
[smart pointers](https://microsoft.github.io/RustTraining/c-cpp-book/ch07-2-smart-pointers-and-interior-mutability.html)
only to distinguish Box from Rc; Cell/RefCell and graph exercises are deferred.
The budget does not include completing all source exercises.

**Six-hour budget:** reading 60m + Starter 30m + Intermediate 120m +
tests 90m + reflection/checks 60m = 360m.

## C++ mental model -> Rust concept

std::string_view and move constructors -> lifetime-checked borrows and moves
that make a non-Copy source unavailable. A lifetime annotation states a
relationship; it cannot keep a local allocation alive.

## Required Starter

Complete W3-S1 and W3-S2 in the supplied `tests/week03.rs`.
In test-sized experiments, move a String, copy a u32, and clone a String.
Record a use-after-move compiler diagnostic, then remove the broken code.
Show a shared borrow followed by a mutable borrow after its last use.
Use a separate test-only Drop marker to set a borrowed boolean on destruction;
assert the boolean after explicitly dropping the marker. Do not combine
Copy and Drop on that type.

## Required Intermediate project

Use the supplied `rust/logscope` crate; weeks 4-6 extend it in place.
Do not run `cargo new`. Follow the
[week 3 file/TODO map and focused commands](../docs/exercises/rust-logscope.md#week-3-ownership-without-parsing).
`cargo build` and `cargo test --no-run` must pass initially. Run only the
week03 acceptance target this week; future active tests intentionally fail.

In `src/model.rs`, the supplied types define Level with Info, Warn, Error
variants, a `Record<'a>` containing level, `service: &'a str`, and `duration: u32`,
and an OwnedRecord with an owned String service. Implement W3-1 through W3-3.

Required APIs: `record(level: Level, service: &str, duration: u32) -> Record<'_>`,
`service_name<'a>(record: &Record<'a>) -> &'a str`, and
`Record::into_owned(self) -> OwnedRecord`. Expose enough fields/accessors for
tests. These accept already-typed data: no validation or parsing is due.
Convert service storage exactly when producing OwnedRecord, not when borrowing.

## Gate, tests, and deliverables

After implementing, run from the course root.

**Windows PowerShell**

```powershell
Push-Location .\rust\logscope
cargo test --test week03
cargo fmt --check
Pop-Location
```

**macOS Terminal**

```sh
(
  cd rust/logscope &&
  cargo test --test week03 &&
  cargo fmt --check
)
```

Construct Info/api/10, Warn/api/20, and Error/worker/100 records in memory.
Verify fields, all levels, and zero duration. Compare the service pointer with
the original string slice's pointer to demonstrate borrowing. Verify an owned
conversion remains usable after the original String leaves scope. Verify a
returned service can outlive the temporary Record wrapper but not the String.

**Gate:** explain who owns service storage in both types, demonstrate the
borrow/owned tests, and explain why returning a local borrow is rejected.
Do not implement parse_line, Result-based validation, or the malformed-input
matrix early. Those are week 4's chapter 9 work.

Submit the crate and `docs/submissions/week-03.md`.
**Rubric:** Starter experiments 20 + record/lifetime design 35 +
ownership tests 25 + reflection/formatting 20 = 100.

**C++ questions:** (1) What obligation would std::string_view leave to the
caller? (2) How do Rust's unavailable moved-from binding and a valid-but-
unspecified moved-from C++ standard-library value differ?

**Optional Challenge (ungraded):** repeat an ownership experiment with Box;
do not add smart pointers to the record model merely to use them.
