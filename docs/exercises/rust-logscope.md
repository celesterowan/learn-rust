# Rust logscope: one scaffold, weeks 3-6

Use the existing `rust/logscope` directory; do not run `cargo new`, create a
workspace, or copy a project for each week. The library evolves in place.
The package, binary, and library import are all named `logscope`.
Only std is used. Reading budgets, rubrics, gates, submission paths, and
reflection prompts remain in [week 3](../../assignments/week-03.md),
[week 4](../../assignments/week-04.md), [week 5](../../assignments/week-05.md),
and [week 6](../../assignments/week-06.md).

The [frozen shared contract](../logscope-contract.md) is authoritative for
bytes, malformed records, arithmetic, deterministic output, and CLI behavior.
These types and tests make it executable; they do not supply implementations.
Build and test compilation must pass before starting. Active acceptance tests
intentionally fail until the corresponding TODO is complete. There are no
ignored/future-disabled tests. **Select the current week's targets**, not
`cargo test` without selectors in week 3 or 4.

## Layout and milestones

Paths in this table are relative to `rust/logscope`.

| File | Work due |
| --- | --- |
| `src/model.rs` | Week 3: typed borrowed/owned records, W3-1..3 |
| `tests/week03.rs` | Week 3: ownership tests; W3-S1 experiments and W3-S2 Drop |
| `src/parser.rs` | Week 4 Starter: strict normalized-line parser, W4-1 |
| `src/summary.rs` | Week 4: W4-2..6 checked updates, normalization, analysis, rendering |
| `src/error.rs` | Supplied typed errors/display/source; week 5 W5-2 conversion |
| `src/command.rs`, `src/main.rs` | Supplied sequential CLI; week 5 W5-5 reader wiring |
| `src/reader.rs` | Week 5: W5-3 generic reader |
| `src/summary.rs` | Week 5: W5-1 Default, W5-4 iterator refactor |
| `src/parallel.rs` | Week 6: W6-2 joins, W6-3 bounded analysis |
| `src/summary.rs` | Week 6: W6-1 checked summary merge |
| `tests/week06.rs` | Week 6: W6-S two-thread Starter and acceptance |
| `tests/week04_parser.rs`, `week04_summary.rs`, `week04_cli.rs` | Week 4 and later regression targets |
| `tests/week05.rs` | Week 5 and later reader/conversion regressions |
| `tests/common/mod.rs` | Supplied portable shared-fixture paths |
| `src/lib.rs`, `Cargo.toml` | Supplied module declarations, re-exports, manifest |

Future modules are already declared so the final project always compiles.
Do not implement them early. A `todo!` in a non-Result function signals unfinished
work; a Result function returns a typed `Unimplemented` error. These are never
successful defaults. In particular, do not count `ParseError::Unimplemented`
as malformed input: propagate it as `AnalyzeError::Unimplemented` so an
unfinished parser cannot produce a plausible success report.

## Week 3: ownership without parsing

Complete only W3 tasks; the future CLI is supplied but not due yet.

1. **W3-S1**, `tests/week03.rs`: replace the Starter TODO with passing
   assertions demonstrating String move, u32 copy, String clone, and a shared
   borrow followed by mutable borrowing after its last use. Predict and
   temporarily observe a use-after-move diagnostic, then remove broken code
   and record the diagnostic/correction in the reflection.
2. **W3-S2**, `DropMarker::drop` in that test: set the borrowed boolean.
   The provided test explicitly drops the marker and checks the flag.
   Keep this test-only marker separate from Copy types.
3. **W3-1**, `record(level, service, duration) -> Record<'_>`: construct
   the supplied record without allocation or validation. Inputs are already
   typed, so even an empty service or a large u32 duration is allowed here.
4. **W3-2**, `service_name<'a>(&Record<'a>) -> &'a str`: preserve the
   service-storage lifetime, not the wrapper's shorter lifetime.
5. **W3-3**, `Record::into_owned(self) -> OwnedRecord`: consume the record,
   keep level/duration, allocate a String for the service. Only this conversion
   owns new service storage.

Focus with `cargo test --test week03 starter_` or
`cargo test --test week03 owned_conversion_outlives_storage`.
Gate: `cargo test --test week03` and `cargo fmt --check`.
The tests check all levels, zero and large duration, pointer identity, owned
storage surviving the source, and borrowed service surviving the wrapper.
Explain why a borrow still cannot outlive its original String. No byte parsing,
Result error handling, malformed matrix, Default, readers, or threads is due.

## Week 4: parser, checked summary, sequential command

Complete W4 in order; retain the week 3 tests.

1. **W4-1, Starter**, `parse_line(&[u8]) -> Result<Record<'_>, ParseError>`:
   accept exactly one normalized nonempty line from the shared grammar.
   Empty input is `Malformed`. Do not split LF, strip CR, or trim here.
   Validate ASCII and all fields before converting the service field to &str;
   borrow those original bytes. Every bad line returns `Malformed`, not a
   panic or a duration of zero. First run
   `cargo test --test week04_parser starter_parser`, then the whole target.
2. **W4-2**, `summary::checked_add(u64, u64) -> Result<u64, AnalyzeError>`:
   mathematical addition or `Overflow`. Test near-limit values directly with
   `cargo test --test week04_summary checked_add_boundaries`.
3. **W4-3**, `Summary::add_record(&Record)`, `Summary::add_invalid()`:
   count valid/invalid lines and all per-level counts using u64. Own service
   keys in the supplied BTreeMap. Check service count and total_ms additions;
   max_ms is the greatest duration, not a sum. The helper accepts typed u32
   durations, including u32::MAX, so test wide sums independently of parsing.
   On error, callers discard the summary; transactional rollback is not required.
   Use `Summary::empty()` for now; Default is a week 5 task.
4. **W4-4**, `summary::process_line(&mut Summary, &[u8])`: receive one raw
   LF-split segment, remove at most one final CR, ignore it if now empty,
   otherwise parse and update the summary. Only `Malformed` increments invalid
   once; arithmetic and unfinished-function errors abort.
5. **W4-5**, `analyze(&[u8]) -> Result<Summary, AnalyzeError>`: split bytes
   on LF and reuse process_line, including an unterminated final segment.
   Do not decode a whole file, ignore invalid UTF-8, or lose malformed counts.
   The returned Summary owns its keys and can outlive the input.
6. **W4-6**, `render(&Summary) -> String`: build the complete deterministic
   report before writing anything. Include every level and sorted services
   with exactly one final LF. BTreeMap supplies key order; explain the
   ordering rather than relying on an unordered map's iteration.
   Gather service lines with a simple collection loop now, leaving a focused
   map/collect refactor for week 5. Use at least one explicit match and one
   appropriate `?` in your learner implementation.
7. Trace the supplied `command::report_path`, `command::execute`, and main.
   They accept exactly one OS path (even a dash-leading name), read bytes,
   call your library, and write/flush the completed report. Wrong arity exits
   2; read/analysis/write/flush failures exit 1; success exits 0. Error paths
   have `error:` stderr and no success report; write failure may leave partial
   stdout, and stderr failure is best-effort. No second parser belongs here.

The three `week04_*` targets cover fixture bytes, zero-byte input, normalization,
CRLF/no final LF/single-CR stripping, field/byte grammar, sorted owned keys,
all update overflow paths, wide sums, compiled-binary streams/status,
missing files, directories, spaces/dash-leading paths, and injected write/
flush failures. Matrix failures must be errors or invalid counts, never panics.
Gate: week03 plus all three week04 targets, formatting, and clippy. Future
week05/week06 tests are not required yet.

## Week 5: traits and a generic reader

1. **W5-1, Starter**, `Default for Summary`: produce the same empty state as
   `Summary::empty()`. `cargo test --test week05 default_is_empty_report`.
2. **W5-2, Starter**, `From<io::Error> for AnalyzeError`: wrap the original
   io::Error in `Io`, retaining its kind, message, and inspectable source.
   Display/source implementations are supplied; do not stringify the cause.
   `cargo test --test week05 from_and_into_preserve_io_cause`.
3. **W5-3**, `analyze_reader<R: Read>(reader: R) -> Result<Summary, AnalyzeError>`:
   read to EOF into bytes, then reuse analyze. Use `?` with your From
   conversion. A failure at any point returns an error, never a summary of a
   valid prefix. The argument can itself be File, Cursor, or a borrowed reader.
   This is a whole-file testability boundary, not a streaming-memory promise.
4. **W5-4**, `render`: refactor one collection loop to a closure with
   map/collect without changing bytes or ordering. Do not discard malformed
   input with filter_map. Explain iter/into_iter and closure captures.
5. **W5-5**, `command::report_path`: replace the supplied byte read with
   File::open followed by analyze_reader. Preserve render and the command
   error/exit behavior. Update imports to avoid unused imports. No new CLI
   flags, Cargo features, or duplicate aggregation are needed.

`tests/week05.rs` supplies real Read implementations for one-byte reads and
failure after a prefix (including immediate failure), plus Cursor and File.
No custom reader framework is required. Run `cargo test --test week05`;
retain all earlier targets. The iterator refactor and CLI reader call are
source-review/reflection gates in addition to the automated tests.
Future week06 tests remain intentionally unfinished.

## Week 6: bounded scoped concurrency, unchanged CLI

1. **W6-S, Starter**, `tests/week06.rs`: replace its TODO with two scoped
   threads borrowing fixed complete records. Give each a private count,
   explicitly join both, and assert the sum. No shared mutable map.
2. **W6-1**, `Summary::merge(other: Summary) -> Result<(), AnalyzeError>`:
   combine every count and duration sum with checked addition; merge both new
   and overlapping owned keys, taking the greater max_ms. Failure discards
   the aggregate, not a silently partial success. Test with
   `cargo test --test week06 merge_`.
3. **W6-2**, `parallel::join_all(Vec<ScopedJoinHandle<...>>)`: explicitly join
   every handle and inspect both join and worker results. Convert a panic to
   `WorkerPanicked`. Merge successes with your checked helper. Remember the
   first failure in handle order (including merge overflow) but still join
   all later handles. Only return an aggregate when every join/worker/merge
   succeeded; an empty handle list returns empty state.
   `cargo test --test week06 join_all_` injects worker errors, panics, delayed
   completion, and merge overflow. Dropping handles and relying on scope's
   implicit joins does not satisfy this contract.
4. **W6-3**, `analyze_parallel(&[u8], workers: usize)`: reject workers outside
   1..=4 even for empty input. Split on LF into complete borrowed line slices,
   then partition these slices across at most workers scoped threads.
   Each worker owns a private Summary and reuses process_line exactly once
   per raw segment; do not strip CR twice. Reuse join_all and merge.
   Empty partitions are valid; no arbitrary byte chunks, thread-per-line,
   unsafe, async, or shared mutable aggregate.

Run the entire week06 target, then `cargo test` for final acceptance.
The equivalence matrix covers all fixtures and worker counts, zero/fewer lines
than workers, CRLF, unterminated input, mixed valid/invalid bytes, and repeated
input. Source review confirms the thread bound and use of scoped borrowing.
Tests deliberately panic only inside injected test workers. Panic-to-error
recovery assumes normal unwinding, not a panic=abort build. The sequential
CLI still calls analyze_reader; do not add a parallel flag.

## Build, focused tests, and final acceptance commands

All commands start at the course root; run only your OS block.
Initial `cargo build`, `cargo test --no-run`, and formatting must pass.
Acceptance commands fail at explicit TODOs until that week is implemented.
Do not run past a failure and claim that the final command proves all passed.

**Windows PowerShell**

```powershell
Push-Location .\rust\logscope
cargo build
cargo test --no-run
cargo fmt --check
# Week 3 gate:
cargo test --test week03
# Week 4 gate: retain week03 and add these:
cargo test --test week04_parser --test week04_summary --test week04_cli
# Week 5 gate: retain earlier targets and add:
cargo test --test week05
# Week 6 gate and final acceptance:
cargo test
cargo clippy --all-targets -- -D warnings
cargo run -- ..\..\fixtures\sample.log
Pop-Location
```

**macOS Terminal** (choose only the current week's acceptance line)

```sh
(
  cd rust/logscope &&
  cargo build &&
  cargo test --no-run &&
  cargo fmt --check
)
# Run after completing the corresponding week:
(cd rust/logscope && cargo test --test week03)
(cd rust/logscope && cargo test --test week03 --test week04_parser --test week04_summary --test week04_cli)
(cd rust/logscope && cargo test --test week03 --test week04_parser --test week04_summary --test week04_cli --test week05)
(cd rust/logscope && cargo test && cargo clippy --all-targets -- -D warnings && cargo run -- ../../fixtures/sample.log)
```

Run `cargo fmt --check` at every gate; clippy is due from week 4 onward.
From week 4, run the sample CLI and later malformed.log as the handouts direct.
Tests derive shared fixtures by joining `CARGO_MANIFEST_DIR/../../fixtures`;
they never embed this checkout's absolute location. They compare bytes and
invoke `env!("CARGO_BIN_EXE_logscope")`, not a cargo-run wrapper.
Keep learner work and the supplied acceptance tests; extend tests if useful.
Do not disable assertions, replace TODO errors with empty successes, change
fixtures, add dependencies, or modify the frozen shared contract.
