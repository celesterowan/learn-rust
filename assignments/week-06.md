# Week 6: Bounded parallel analysis with scoped threads

**Prerequisite:** week 5 gate. **Objectives:** borrow across safe scoped threads,
merge independent results, and make worker failures visible.

**Primary reading (45m):** Microsoft's
[ch13 concurrency](https://microsoft.github.io/RustTraining/c-cpp-book/ch13-concurrency.html)
(spawn, scope, move captures, joins, Send/Sync; skim Arc versus Mutex).
Then consult [thread::scope](https://doc.rust-lang.org/std/thread/fn.scope.html)
and [scoped join handles](https://doc.rust-lang.org/std/thread/struct.ScopedJoinHandle.html).
Skip the source's one-thread-per-line challenge and lock/atomic/channel
exercises. Our bounded project replaces those exercises, rather than adding
to them. No unsafe or async.

**Six-hour budget:** reading 45m + Starter 30m + partition/merge design 30m +
Intermediate 105m + tests 105m + reflection/checks 45m = 360m.

## C++ mental model -> Rust concept

std::thread lifetime obligations -> scoped borrowing, move closures, and
Send/Sync constraints. Shared ownership alone does not synchronize mutation.

## Required Starter

Complete W6-S in the supplied `tests/week06.rs`: use two scoped threads on
fixed, complete in-memory records and join them.
Give each a private count, then assert the sum after both joins. Explain why
the scope can borrow a local buffer. This is a small test feeding the project,
not a second command or shared-map implementation.

## Required Intermediate project

Implement the supplied library-only extension:
`analyze_parallel(input: &[u8], workers: usize) -> Result<Summary, AnalyzeError>`.
Accept workers 1..4; otherwise return an error. Keep the sequential CLI
unchanged. Normalize and partition complete lines, never arbitrary byte chunks
that can cut a record. Empty lines still follow the shared contract.

Follow the [week 6 TODO/helper map and focused selectors](../docs/exercises/rust-logscope.md#week-6-bounded-scoped-concurrency-unchanged-cli).
Implement checked merge, join reduction, then partition/analysis in that order.
The supplied acceptance tests inject worker failures and merge overflow; do
not replace these with skipped tests or silently successful empty results.

Use at most workers scoped threads, each with its own Summary and borrowed
line slices. A simple fixed partition of a vector of line slices is sufficient;
you do not need a task queue. No thread per line and no shared mutable map.
Reuse line processing and checked merge helpers. Empty partitions are valid.

Explicitly join every handle, inspect every result, and return an error if
any worker returns an error or panics. Retain all handles and join them even
after finding an error. Never silently drop a worker's data. Convert a joined
panic to an analysis error; do not intentionally panic in ordinary processing.
Do not promise recovery when the process is built with panic=abort.

## Gate and required tests

Compare sequential and parallel rendered bytes for all fixtures, worker
counts 1..4, zero-byte input, fewer lines than workers, CRLF, an unterminated
last line, and mixed valid/invalid lines. Test rejected worker counts.
Unit-test checked summary merging near overflow, including overlapping
service keys. Test your join/reduction helper using test-only workers that
return Err or panic; verify no success summary escapes.

From the course root, run only your OS's block.

**Windows PowerShell**

```powershell
Push-Location .\rust\logscope
cargo test
cargo fmt --check
cargo clippy --all-targets -- -D warnings
Pop-Location
```

**macOS Terminal**

```sh
(
  cd rust/logscope &&
  cargo test &&
  cargo fmt --check &&
  cargo clippy --all-targets -- -D warnings
)
```

Submit crate and `docs/submissions/week-06.md` with a diagram of ownership
and join lifetime.
**Gate:** all equivalence/error tests pass, no thread-per-line creation,
every handle joined even after failure. Explain Send versus Sync before Go.
**Rubric:** Starter 15 + bounded ownership/partitioning 25 +
checked merge/failure handling 25 + tests 20 + reflection/formatting 15 = 100.

**C++ questions:** (1) What prevents borrowed input from dying before a worker?
(2) Compare the failure behavior with joining std::thread workers.

**Optional Challenge (ungraded):** time repeated input sequentially and in
parallel; explain overhead even when parallel is slower.
