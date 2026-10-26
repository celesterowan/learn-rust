# Week 9: Interfaces and resource lifetimes you can test

**Prerequisite:** week 8. **Objectives:** use implicit interfaces, defer
correctly, and automate command-level failure cases.

**Reading (45m):** [interfaces](https://go.dev/tour/methods/9)
(Tour methods pages 9-12), [defer](https://go.dev/tour/flowcontrol/12)
(pages 12-13), and [io.ReadAll](https://pkg.go.dev/io#ReadAll).

**Six-hour budget:** reading 45m + design 30m + implementation 105m +
tests 135m + reflection/checks 45m = 360m.

## Assignment

Add `AnalyzeReader(r io.Reader) (Summary, error)`. Read all bytes and reuse
Analyze. If reading returns bytes and an error, return an error, not a partial
summary. Preserve causes using wrapping. Add
`AnalyzeFile(path string) (Summary, error)` which opens the file and defers
Close after a successful Open. The CLI uses AnalyzeFile.

Explicit policy for this read-only helper: report a Close error if analysis
succeeded; preserve the earlier read/analysis error if both fail. A named
return or small helper can implement this; explain your choice. Add an
unexported helper taking io.ReadCloser so tests can inject Close outcomes.
Do not invent a large filesystem interface.

## Required tests and checks

Use table-driven tests with named cases for the shared matrix.
Use strings.NewReader or bytes.NewReader for success; test-only types implement
io.Reader/io.ReadCloser without an inheritance declaration.
Cover immediate read failure, partial bytes plus error, one-byte reads,
successful close, close failure, and simultaneous read/close failure.
Assert Close is called exactly once on every successfully opened test handle.

Automate CLI tests with os/exec: build once into a t.TempDir path, then run
that executable for success, zero/multiple arguments, a nonexistent file,
and a temporary path with spaces. Compare captured bytes, stderr prefix, and
exit status. Fail the test if the binary build fails. Do not assert exact
platform-specific OS messages or run through a shell.

```powershell
Push-Location .\go\logscope
go test ./...
gofmt -l .
go vet ./...
Pop-Location
```

Submit module and `docs/submissions/week-09.md`.
**Rubric:** reader/interface boundary 25 + resource/error policy 25 +
table/process tests 35 + reflection/formatting 15 = 100.

**C++ questions:** (1) How is defer's function scope different from a local
RAII object's block scope? (2) Why does a tiny io.Reader consumer need no
explicit "implements" declaration?

**Optional stretch (ungraded):** add a typed-nil interface demonstration as a
test with a documented expected result, without introducing it into production
paths.
