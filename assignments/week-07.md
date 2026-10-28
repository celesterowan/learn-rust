# Week 7: Go statistics and the aliasing lab

**Prerequisite:** weeks 1-6; Go setup. **Objectives:** learn Go syntax,
value/pointer behavior, slice backing arrays, and map zero values.

**Reading (60m):** [A Tour of Go: basics](https://go.dev/tour/basics/1)
(first 8 pages), [pointers through slices](https://go.dev/tour/moretypes/1)
(pages 1-16), and [maps](https://go.dev/tour/moretypes/19)
(pages 19-22). Run the small examples, not every Tour exercise.

**Six-hour budget:** reading 60m + module/tool setup 30m +
implementation 105m + tests/aliasing lab 105m + reflection/checks 60m = 360m.

## Assignment

Use the supplied module, executable wiring, and tests in `go/stats`; do not
run `go mod init` over it. Read the [shared statistics spec](../docs/exercises/stats.md)
for ordered G1-G6 TODOs, lab contracts, and separate compilation/acceptance
commands. Active tests intentionally fail until implementation.

Reimplement week 2's CLI grammar, output, and error statuses in Go.
Use uint32 inputs and uint64 count/sum; the pure function takes a slice
and returns `(Stats, error)`.
Use the supplied `HasMin` and `HasMax` presence booleans for absent extrema,
not a sentinel numeric value. Keep the scaffold's Stats fields stable and
check addition.

Reuse the week 2 cases, including a pure-function value of uint32's maximum.
Keep validation separate from main. Do not use panic for bad input.
No argument means an empty dataset, not usage failure.

## Aliasing lab and acceptance

Use and extend the supplied tests showing: a struct copied by value preserves
independent scalar fields; a pointer mutation affects the original; two slices share an element;
an explicitly copied slice does not; two map variables share entries; and a
missing map key is distinguishable from a present zero. Include a test for
append with deliberately sufficient capacity, not an assumed growth policy.
Do not commit a test that intentionally panics on nil-map assignment.

The sample output must remain `count=4 sum=17 min=0 max=9` plus LF.
After implementing source and tests, run from the course root.

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

Formatting check lists no files; all tests and vet pass.
Submit module and `docs/submissions/week-07.md`.
**Rubric:** stats/CLI parity 35 + aliasing lab 25 + edge-case tests 25 +
reflection/formatting 15 = 100.

**C++ questions:** (1) Which Go assignments resemble copying a vector, and
which resemble copying a span? (2) Why does garbage collection not guarantee
that two callers cannot mutate the same slice?

**Optional stretch (ungraded):** repeat the slice test with a struct containing
a slice and explain why the struct copy is not a deep copy.
