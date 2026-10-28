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

Complete G11.1-G11.4 in the supplied
[Go logscope exercise spec](../docs/exercises/go-logscope.md). The separate
`cmd/logserve/main.go` signal/process adapter and active tests are included;
the handler, server configuration, startup, and shutdown remain learner TODOs.
No new module is needed. Follow the spec's dependency seams and ordered tests.

Implement the testable handler constructor in the library.
It accepts an immutable rendered report. Do not parse request
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
The supplied tests also exercise startup failures and real ephemeral-loopback
server cancellation, graceful completion, and forced-close fallback. Run and
inspect them, then extend the cases; tests never occupy port 8080 or leave a
server running. Shutdown failure remains an error even after fallback Close.

From the course root, run only your OS's block. Resolve any unavailable
[race-tool prerequisites](../SETUP.md#race-detector-prerequisites) or record
that check as incomplete.

**Windows PowerShell**

```powershell
Push-Location .\go\logscope
go test ./...
go test -race ./...
gofmt -l .
go vet ./...
Pop-Location
```

**macOS Terminal**

```sh
(
  cd go/logscope &&
  go test ./... &&
  go test -race ./... &&
  gofmt -l . &&
  go vet ./...
)
```

Then build and run the server from the course root.

**Windows PowerShell**

```powershell
Push-Location .\go\logscope
New-Item -ItemType Directory -Force .\bin
go build -o .\bin\logserve.exe .\cmd\logserve
if ($LASTEXITCODE -ne 0) { Pop-Location; throw 'Build failed.' }
.\bin\logserve.exe ..\..\fixtures\sample.log
$LASTEXITCODE
Pop-Location
```

**macOS Terminal**

```sh
(
  cd go/logscope &&
  mkdir -p bin &&
  go build -o bin/logserve ./cmd/logserve || exit
  ./bin/logserve ../../fixtures/sample.log
  printf 'exit=%s\n' "$?"
)
```

The server blocks. In a second terminal, from any directory, use your OS's
request command:

**Windows PowerShell**

```powershell
curl.exe -i http://127.0.0.1:8080/summary
```

**macOS Terminal**

```sh
curl -i http://127.0.0.1:8080/summary
```

Press Ctrl+C in the first terminal. Its next command displays the server's
exit status, which must be 0 on successful shutdown.
If port 8080 is occupied, do not kill unrelated processes; record the conflict
and retry when free. Tests must not depend on that fixed port.

Submit module and `docs/submissions/week-11.md`, including shutdown evidence.
**Rubric:** endpoint contract 30 + startup/shutdown/resource handling 25 +
handler tests 30 + reflection/formatting 15 = 100.

**C++ questions:** (1) Why is immutable captured text safe to share here?
(2) Which resource lifetimes are explicit rather than garbage-collected?

**Optional stretch (ungraded):** extend the supplied ephemeral-loopback tests
with additional concurrent-client cases.
