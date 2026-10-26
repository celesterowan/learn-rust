# Week 5: Traits, conversions, closures, and a reader boundary

**Prerequisite:** week 4 gate. **Objectives:** use a useful generic bound,
explicit conversions and a small iterator refactor without changing behavior.

**Primary reading (60m):** Microsoft's
[ch10 traits](https://microsoft.github.io/RustTraining/c-cpp-book/ch10-traits.html)
(impl and static versus dynamic dispatch),
[generics](https://microsoft.github.io/RustTraining/c-cpp-book/ch10-1-generics.html)
(trait bounds only),
[ch11 conversions](https://microsoft.github.io/RustTraining/c-cpp-book/ch11-from-and-into-traits.html)
(From/Into, Default, fallible versus infallible), and
[ch12 closures/iterators](https://microsoft.github.io/RustTraining/c-cpp-book/ch12-closures.html)
(capture basics and map/collect). Use
[iterator tools](https://microsoft.github.io/RustTraining/c-cpp-book/ch12-1-iterator-power-tools.html)
as a lookup, not an exercise set. Skip operator overloads, type-state/builders,
custom iterators, and boxed closure collections. Look up
[Read](https://doc.rust-lang.org/std/io/trait.Read.html) and
[Cursor](https://doc.rust-lang.org/std/io/struct.Cursor.html) while coding.

**Six-hour budget:** reading 60m + Starter 30m + Intermediate 105m +
tests 105m + reflection/checks 60m = 360m.

## C++ mental model -> Rust concept

Constrained templates, conversion constructors, and lambdas -> trait bounds,
explicit From/Into conversions, and borrow-checked closure captures.

## Required Starter

Implement Default for Summary as the empty report state. Implement
`From<std::io::Error> for AnalyzeError`, preserving the original error.
Test Default rendering against empty.expected.txt and the conversion through
both From and Into. These conversions are infallible wrappers, not permission
to make invalid log input valid.

## Required Intermediate project

In the existing crate, add
`analyze_reader<R: std::io::Read>(reader: R) -> Result<Summary, AnalyzeError>`.
Extend your error type to preserve an I/O error. Read all bytes, then reuse
analyze. Use your From implementation to propagate I/O errors with `?`.
This is a testability boundary, not a streaming-memory promise.
The CLI now opens a File and calls this function. Preserve every shared
contract behavior and the earlier byte-slice API.

Write a tiny test-only reader implementing Read. It returns some bytes, then
an I/O error. Assert that analyze_reader returns an error rather than a
summary of the prefix. A second reader can return one byte per read; the
complete result must equal Cursor over the same input. Do not assume a single
read fills a requested buffer.

Refactor one renderer collection loop to a closure with map/collect (sorting
remains explicit). Keep checked aggregation as a loop if clearer. In the
reflection, explain iter versus into_iter and what the closure captures.
Do not use filter_map to discard malformed lines without incrementing invalid.

## Gate, checks, and acceptance

```powershell
Push-Location .\rust\logscope
cargo test
cargo fmt --check
cargo clippy --all-targets -- -D warnings
cargo run -- ..\..\fixtures\malformed.log
Pop-Location
```

Keep the existing parser, fixture, overflow, and CLI tests. Add tests for
Cursor containing empty input, sample bytes, an immediate read failure,
and a failure after a valid record. Test your error display without comparing
platform-specific OS error wording. Error variants or source chaining must
make the I/O cause inspectable. Do not create a custom reader framework or
duplicate aggregation for each reader type.

Submit updated crate and `docs/submissions/week-05.md`, including a brief
before/after call graph and evidence that output did not change.
**Gate:** show the generic reader working with both File and Cursor, injected
failure propagation, and unchanged fixture bytes after the iterator refactor.
**Rubric:** Starter/conversions 20 + generic reader/iterator refactor 30 +
regression/failure tests 30 + explanation/formatting 20 = 100.

**C++ questions:** (1) Compare `R: Read` with a template constraint and with
a virtual base class. (2) Who closes File, and what does this generic
function own when passed a File versus a mutable reference?

**Optional Challenge (ungraded):** call the API through `&mut dyn Read` and
explain why that call can work without redesigning the API.
