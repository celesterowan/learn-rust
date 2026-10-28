# Week 12: cross-language acceptance and comparison

This is the final check of the existing Rust and Go analyzers, not a new
implementation project. Complete the week 1-11 gates first. The supplied
`go/logscope/parity_test.go` builds and executes both sequential CLIs and
compares each with independent expected bytes from the
[frozen contract](../logscope-contract.md). Two implementations agreeing on the
same mistake is not enough.

## What to do

1. Finish the earlier learner TODOs and rerun their focused tests. Do not copy
   an implementation between languages or change fixtures to match your output.
2. Read `TestParity` in `go/logscope/parity_test.go`. Trace where fixture paths,
   process argument arrays, stdout, stderr, status, and timeouts are handled.
   The harness itself is supplied testing infrastructure.
3. Run the tagged test below. Fix discrepancies in your implementations.
   Failures identify the case and language. Fatal diagnostic wording after
   `error:` may differ; success output, empty streams, and statuses may not.
4. Rerun every project's final acceptance checks, including Go race checks on
   a supported toolchain. A compilation-only run is not passing acceptance.
5. Write `docs/comparison.md` and the week 12 reflection as specified by the
   [handout](../../assignments/week-12.md). Cite concrete code/test evidence.
   The comparison and reflections are learner submissions, not supplied answers.

## Supplied cases

| Case | Expected behavior |
| --- | --- |
| Sample, empty, malformed fixtures | Exact corresponding `.expected.txt` bytes, status 0, empty stderr |
| Temporary filename containing spaces | Sample report, status 0, empty stderr |
| CRLF, blank line, final unterminated CR | Two api records, total 10, maximum 10 |
| Invalid UTF-8 byte and NUL line | One valid api record, two invalid lines |
| Filename beginning with a dash | A path, not a CLI option |
| No arguments or two arguments | Status 2, empty stdout, `error:` stderr |
| Missing file | Status 1, empty stdout, `error:` stderr |

The source contains the manually calculated expectations for temporary inputs.
No shell redirection or newline normalization is used. Each process has a
15-second limit; builds have a three-minute limit. Temporary builds and inputs
are removed by Go's test cleanup.

## Commands

Both `cargo` and `go` must be on PATH. The harness uses an explicit temporary
Cargo `--target-dir`, so it does not depend on a preexisting debug executable
or a globally configured target directory. Use the native host target, not a
cross-compilation target. Everything is standard-library-only and builds
offline. The executable suffix is selected for Windows or macOS automatically.

The `parity` build tag makes the cross-language dependency explicit. Earlier Go
weeks can run their tests without Rust; the capstone must run the tagged tests.
There are no skipped parity cases. Untagged `go test` does not run this harness.

### Windows PowerShell, from the course root

```powershell
Push-Location .\go\logscope
$env:GOTOOLCHAIN = 'local'
# Compile the complete tagged test suite, without running learner behavior.
New-Item -ItemType Directory -Force .\bin
go test -tags parity -c -o .\bin\logscope-parity.test.exe .
# Run all parity cases; TODO failures are not a completed capstone.
go test -tags parity -count=1 -run '^TestParity$' .
Pop-Location
```

### macOS Terminal, from the course root

```sh
(
  cd go/logscope || exit
  export GOTOOLCHAIN=local
  mkdir -p bin
  go test -tags parity -c -o bin/logscope-parity.test .
  go test -tags parity -count=1 -run '^TestParity$' .
)
```

If a build fails or a tool is missing, fix setup before interpreting behavioral
results. The unfinished starter CLIs intentionally fail acceptance. By week 12,
every case must pass without weakening tests, removing the tag from final
commands, suppressing errors, or returning a partial report.
