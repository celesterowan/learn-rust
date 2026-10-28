# Week 10: A bounded worker pool that actually stops

**Prerequisite:** week 9. **Objectives:** coordinate channels, cancellation,
and joining without shared-map races or silently lost errors.

**Reading (60m):** [goroutines and channels](https://go.dev/tour/concurrency/1)
(pages 1-6), [Context](https://pkg.go.dev/context#Context),
and [race detector usage](https://go.dev/doc/articles/race_detector.html).

**Six-hour budget:** reading 60m + lifecycle design 45m +
implementation 120m + tests 90m + reflection/checks 45m = 360m.

## Assignment

Complete G10.1-G10.2 in the supplied
[Go logscope exercise spec](../docs/exercises/go-logscope.md). It fixes the
checked-merge API, processing-callback seam, ordered tests, and exact lifecycle
contract. Do not initialize a new module or alter the sequential command.
The focused cumulative gate excludes week-11 TODOs.

Add a library-only extension:
`AnalyzeParallel(ctx context.Context, input []byte, workers int) (Summary, error)`.
Workers must be 1..4. Validate this first; then check ctx.Err before starting.
The caller passes a non-nil context. The sequential CLI stays unchanged.

Split normalized complete lines using shared logic. One caller-side producer
sends line jobs through a channel buffered to at most workers. Start exactly
workers worker goroutines, each owning a local Summary. Workers send exactly
one final result, including errors, into a result channel buffered to workers.
The producer alone closes jobs. No extra producer goroutine is necessary.

Check cancellation when sending/receiving jobs and between records. On worker
error, cancel a derived context to unblock the producer and peers. Collect
every worker result before returning. Merge with checked arithmetic only after
joining. Never return a partial summary as success.

Return a worker/merge error if one occurred; otherwise return the caller's
ctx.Err if canceled, checked again just before success. Internal cancellation
must not replace the original worker failure with context.Canceled.
No panic recovery is required for programmer bugs.

## Required tests

Run and inspect the supplied active tests, then extend them with your own
cases. Their channel handshakes are synchronization, not timing assumptions.

Compare sequential bytes for fixtures, 1..4 workers, empty input, CRLF,
and fewer jobs than workers. Reject 0, negative, and 5 workers.
Test a pre-canceled context and checked-merge overflow. To test cancellation
in flight deterministically, use an unexported worker helper with a test-only
processing callback: signal entry, block on ctx.Done, cancel from the test,
then verify every worker signals exit before the function returns.
Use the same seam to inject a worker error. A generous test timeout is a
deadlock guard, not a sleep-based synchronization mechanism.

Check the [platform-specific race prerequisites](../SETUP.md#race-detector-prerequisites).
From the course root, run only your OS's block.

**Windows PowerShell**

```powershell
Push-Location .\go\logscope
go test -run '^TestWeek(08|09|10)' ./...
go test -race -run '^TestWeek(08|09|10)' ./...
gofmt -l .
go vet ./...
Pop-Location
```

**macOS Terminal**

```sh
(
  cd go/logscope &&
  go test -run '^TestWeek(08|09|10)' ./... &&
  go test -race -run '^TestWeek(08|09|10)' ./... &&
  gofmt -l . &&
  go vet ./...
)
```

Submit module and `docs/submissions/week-10.md`.
**Rubric:** bounded ownership 30 + cancellation/error/join correctness 30 +
equivalence/lifecycle tests 25 + reflection/race evidence 15 = 100.
If race tooling is unavailable, record the failure and mark race validation
incomplete; do not claim a race-clean run.

**C++ questions:** (1) Who closes each channel, and why? (2) Why is context
cancellation not equivalent to forcibly terminating a thread?

**Optional stretch (ungraded):** measure pool overhead without a speed target.
