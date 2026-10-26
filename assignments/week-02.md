# Week 2: Collections, structs, enums, and statistics

**Prerequisite:** week 1 gate. **Objectives:** use arrays, Vec, slices,
structs, enums, and exhaustive matching before the ownership deep dive.

**Primary reading (60m):** Microsoft's
[ch5 data structures](https://microsoft.github.io/RustTraining/c-cpp-book/ch05-data-structures.html)
(arrays, slices, String vs &str, structs, Vec; skim HashMap) and
[ch6 enums/match](https://microsoft.github.io/RustTraining/c-cpp-book/ch06-enums-and-pattern-matching.html)
(enum variants, exhaustive matching, associated methods).
Skip the long reference deep dives and source exercises this week.

**Six-hour budget:** reading 60m + Starter 30m + Intermediate 120m +
tests 90m + reflection/checks 60m = 360m.

## C++ mental model -> Rust concept

std::vector/span and variant -> Vec/borrowed slices and data-carrying enums.
Use a shared slice here as chapter 5 introduces it; moves and exclusive
borrowing become the subject of week 3.

## Required Starter

Define a small enum representing an empty dataset or a nonempty dataset's
count. Write a function classifying a slice and match both alternatives in
tests for empty and two-element inputs. Inspect one array, Vec, and slice:
record which owns storage and which is a view. This is a test-sized exercise,
not a second CLI.

## Required Intermediate project

```powershell
cargo new --vcs none .\rust\week02-stats
Push-Location .\rust\week02-stats
cargo run -- 4 0 9 4
cargo test
cargo fmt --check
Pop-Location
```

Implement a pure function taking `&[u32]` and returning a Stats struct with
count, sum, optional minimum, and optional maximum. Use u64 for count/sum
and Option<u32> for extrema. This function accepts any u32 values, not just
the CLI's smaller range. Use checked addition and return an error on overflow.
An empty slice is valid.

**Small previews of chapter 9:** Some/None model absent extrema;
checked_add returns an Option; numeric parse returns a Result. Use explicit
matches and a tiny error value to reject bad CLI input or overflow. Limit this
preview to these APIs; `?`, error composition, and the log parser wait for
week 4. Do not read all of chapter 9 now or build another validation framework.

The CLI takes zero or more numbers, each exactly 1-4 ASCII digits with numeric
value 0..1000. Leading zeros are allowed; signs and whitespace inside an
argument are not. Validate every argument before printing. Invalid input:
exit 2, empty stdout, diagnostic starting `error:`. Success: exit 0, empty
stderr, one LF-terminated line. The example prints:

```text
count=4 sum=17 min=0 max=9
```

No arguments prints `count=0 sum=0 min=none max=none`. Report internal
arithmetic failure with exit 1 and `error:` rather than panic.

## Gate, tests, and deliverables

Test empty, singleton, duplicates, zero, ascending/descending inputs, and
`u32::MAX` in the pure function. Test rejected `-1`, `1001`, `1.0`, `""`,
`" 1"`, and five digits at the CLI parsing boundary; accept `0001`.
Show that the caller can still read its vector after calculating statistics.
Test a checked-add helper with a near-limit u64 accumulator. No enormous
allocation or ownership-error experiment is required this week.

Submit crate plus `docs/submissions/week-02.md`.
**Gate:** explain empty versus zero-valued input and match every enum variant;
all stated stats/validation tests pass before week 3.
**Rubric:** Starter 15 + statistics/CLI 40 + tests 25 +
reflection/formatting 20 = 100.

**C++ questions:** (1) Why choose a slice over an owned Vec here?
(2) Compare the dataset enum and Option with std::variant and std::optional.

**Optional Challenge (ungraded):** compute a mean as an Option<f64>; do not
change the graded output.
