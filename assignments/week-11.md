# Week 11: Serve one immutable report over local HTTP

**Prerequisite:** weeks 8-10. **Objectives:** use net/http, test handlers
without a real port, and shut down resources with context.

**Reading (45m):** [net/http](https://pkg.go.dev/net/http)
(Handler, Server, Server.Shutdown),
[httptest](https://pkg.go.dev/net/http/httptest) (NewRecorder/NewRequest),
and [signal.NotifyContext](https://pkg.go.dev/os/signal#NotifyContext).

**Six-hour budget:** reading 45m + endpoint design 30m +
implementation 135m + tests 105m + reflection/checks 45m = 360m.

## Assignment

Add `go/logscope/cmd/logserve/main.go` and a testable handler constructor in
the library. It accepts an immutable rendered report. Do not parse request
uploads or share mutable summary maps between requests.

`logserve` accepts exactly one file path. Wrong argument count exits 2;
file/analysis/listen failures exit 1, stderr starting `error:`. Analyze once
at startup using the sequential API. Do not listen on a bad file or refresh
it per request. Bind exactly `127.0.0.1:8080`, not all interfaces.

Endpoint behavior:

- GET `/summary`: status 200, Content-Type `text/plain; charset=utf-8`,
  body exactly the common CLI report bytes.
- Any other method at `/summary`, including HEAD: status 405,
  `Allow: GET`, empty body.
- Any other path, including `/summary/`: status 404, empty body.
  Unknown path takes precedence over method checking.
- Queries do not change the path: `/summary?x=1` behaves like `/summary`.

Set ReadHeaderTimeout to five seconds. Use a dedicated ServeMux or handler,
not global registrations. Handle Ctrl+C through signal.NotifyContext with
os.Interrupt; call Shutdown with a fresh two-second timeout context, then
Close if graceful shutdown fails. Wait for the serving goroutine to finish.
Treat http.ErrServerClosed as normal shutdown, not startup failure.
Successful shutdown exits 0. No TLS, persistence, or deployment work.

## Tests and commands

Test all endpoint cases with httptest, including exact header/body and empty
input report. Concurrent handler tests must not mutate captured data.

```powershell
Push-Location .\go\logscope
go test ./...
go test -race ./...
gofmt -l .
go vet ./...
go build -o .\bin\logserve.exe .\cmd\logserve
.\bin\logserve.exe ..\..\fixtures\sample.log
Pop-Location
```

The server blocks: in a second terminal run
`curl.exe -i http://127.0.0.1:8080/summary`, then Ctrl+C in the first.
If port 8080 is occupied, do not kill unrelated processes; record the conflict
and retry when free. Tests must not depend on that fixed port.

Submit module and `docs/submissions/week-11.md`, including shutdown evidence.
**Rubric:** endpoint contract 30 + startup/shutdown/resource handling 25 +
handler tests 30 + reflection/formatting 15 = 100.

**C++ questions:** (1) Why is immutable captured text safe to share here?
(2) Which resource lifetimes are explicit rather than garbage-collected?

**Optional stretch (ungraded):** a test server on an ephemeral loopback port.
