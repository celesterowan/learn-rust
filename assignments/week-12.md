# Week 12: Comparison capstone and portfolio handoff

**Prerequisite:** weeks 1-11. **Objectives:** prove parity, repair small gaps,
and explain language trade-offs using your own evidence.

**Reading (30m):** reread the
[contract](../docs/logscope-contract.md) and
[C++ comparison guide](../CPP-TO-RUST-AND-GO.md).
Revisit selected Microsoft chapters 7-13 through your week 3-6 notes.
Consult [Rust testing](https://doc.rust-lang.org/book/ch11-03-test-organization.html)
and [Go testing](https://pkg.go.dev/testing) only for remaining test gaps.

**Six-hour budget:** review/reading 30m + parity harness 90m +
small fixes 90m + final checks 60m + comparison/submission 90m = 360m.

## Assignment

No new product. Preserve the two sequential CLIs, library parallel extensions,
and Go HTTP command. Use the supplied `go/logscope/parity_test.go` and
[capstone instructions](../docs/exercises/parity.md). Trace the harness,
run it, and repair implementation gaps rather than writing process-launch
infrastructure from scratch.

The harness builds both CLIs into temporary locations, with an explicit Cargo
target directory. It uses filepath.Join, native executable suffixes, argument
arrays, byte captures, and command timeouts. Both toolchains must be on PATH.
Build failures fail the test. The `parity` build tag keeps this cross-language
requirement out of earlier Go gates; final acceptance must include the tag.

The active cases compare each binary against all three expected fixture files,
not merely against each other. They also cover no/two arguments, nonexistent
files, filenames with spaces or a leading dash, temporary CRLF input, and
invalid bytes with manually calculated expectations. Fatal error wording may
differ after `error:`.

## Final checks

From the course root, run only your OS's block.

**Windows PowerShell**

```powershell
Push-Location .\rust\logscope
cargo fmt --check
cargo test
cargo clippy --all-targets -- -D warnings
Pop-Location
Push-Location .\go\logscope
gofmt -l .
go vet ./...
go test ./...
go test -tags parity -count=1 -run '^TestParity$' .
go test -race ./...
Pop-Location
```

**macOS Terminal**

```sh
(
  cd rust/logscope &&
  cargo fmt --check &&
  cargo test &&
  cargo clippy --all-targets -- -D warnings
) &&
(
  cd go/logscope &&
  gofmt -l . &&
  go vet ./... &&
  go test ./... &&
  go test -tags parity -count=1 -run '^TestParity$' . &&
  go test -race ./...
)
```

Also rerun checks in the independent week01-hello, week02-stats, and Go stats
modules; use their earlier assignment commands. The parity test now requires
both toolchains. Race tooling limitations must remain visible.

## Deliverables and grading

Submit source/tests, `docs/submissions/week-12.md`, and
`docs/comparison.md` (700-1000 words). The comparison must address ownership
and aliases, error propagation, deterministic ordering, resource cleanup,
testing seams, and worker cancellation/join behavior. Cite one example from
your work for each topic. Explain which language you would choose for a small
local analyzer and why; distinguish preference from measured evidence.
No benchmark is required, and elapsed test time is not an algorithm benchmark.

Within the existing comparison/reflection budget, briefly discuss one design
choice from Microsoft's
[ch16 translation cases](https://microsoft.github.io/RustTraining/c-cpp-book/ch16-case-studies.html),
such as a closed enum versus a class hierarchy. Connect it to your Level type
or an ownership decision; do not implement another case-study project.
Account for the sequence you followed: typed borrowed records in week 3,
Result-based parsing/modules in week 4, and trait/conversion refactoring in
week 5. The broader post-course roadmap remains optional and ungraded.

Review your copied repository for artifacts, secrets, and stale instructions.
Update the progress checklist and document remaining gaps honestly.
**Rubric:** byte/status parity harness 35 + full regression evidence 25 +
comparison with concrete examples 30 + clean navigable portfolio 10 = 100.

**C++ questions:** (1) Which C++ habits did you deliberately stop applying?
(2) What safety or lifecycle obligation still belongs to you in both languages?

**Optional stretch (ungraded):** demonstrate the portfolio to another learner
and record one documentation improvement, without publishing automatically.
