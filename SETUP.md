# Setup: Windows first

Install tools yourself in your chosen study environment. Commands below are
PowerShell commands, not actions already performed. Administrator permission
may be needed by installers. Reopen the terminal after installation.

## Rust

Our week 1 follows Microsoft's
[getting-started chapter](https://microsoft.github.io/RustTraining/c-cpp-book/ch02-getting-started.html).
Use its main/Cargo explanation, but follow the Windows prerequisites below.
Skip its optional REPL installation, external crates, and publishing commands.
Later workspace examples are reference only: retain our independent crates.

Use the [official Rust installer](https://www.rust-lang.org/tools/install)
and [Rust Book installation instructions](https://doc.rust-lang.org/book/ch01-01-installation.html).
Choose the stable MSVC toolchain on ordinary x64 Windows. Rust's MSVC linker
requires the Visual Studio C++ build tools and Windows SDK; the installer may
offer guidance. If missing, use
[Visual Studio Build Tools](https://visualstudio.microsoft.com/visual-cpp-build-tools/)
and select "Desktop development with C++", including the MSVC tools and a
Windows SDK. Existing C++ tools may already satisfy this requirement.
Do not switch toolchains merely to hide a linker error.

```powershell
rustup show
rustc --version
cargo --version
rustup component add rustfmt clippy
```

If `link.exe` cannot be found during a build, check the C++ workload first.
Use the appropriate Visual Studio Developer PowerShell if your installation
requires it. Record actual tool versions in your week 1 reflection; do not
pin to version numbers printed in an old tutorial.

## Go

Follow the [official Go download and installation guide](https://go.dev/doc/install).
Use the Windows MSI for a current stable release.

```powershell
go version
go env GOOS GOARCH
```

Use a normal user-writable directory. Do not put this course inside GOROOT,
and do not set GOPATH specially. Go modules handle these exercises.

The race detector needs a supported platform and a C compiler. Windows/amd64
has additional compiler requirements; consult
[race detector documentation](https://go.dev/doc/articles/race_detector.html)
and [Windows C compiler guidance](https://go.dev/wiki/Windows).
A Rust MSVC linker installation alone does not supply the GCC-compatible
compiler Go's race tooling may require. If unavailable, record the exact
failure and run race tests on a supported environment you control when
possible; passing ordinary tests is not evidence of race freedom.

## Directory plan

The documentation and fixtures exist; the following implementation directories
are created by you when assigned:

```text
rust-go-course/
  assignments/
  docs/submissions/
  fixtures/
  rust/
    week01-hello/       independent Cargo binary crate
    week02-stats/       independent Cargo binary crate
    logscope/           typed-record library week 3; CLI added week 4
  go/
    stats/             independent Go module, week 7
    logscope/          Go module, weeks 8-12
      cmd/logscope/    sequential CLI
      cmd/logserve/    HTTP command added in week 11
```

Each Rust directory owns its Cargo.toml and Cargo.lock. Each Go directory owns
its go.mod. Commit Cargo.lock for these application exercises. With no external
Go packages, go.sum may not exist; that is normal.
Do not rerun cargo new in logscope in week 4; extend week 3's library with
src/main.rs and declared modules. No source-guide rand/serde/thiserror or
property-testing dependency is required.

Create the submission directory yourself once, from your copied course root:

```powershell
New-Item -ItemType Directory -Force .\docs\submissions
```

Representative checks, after the corresponding module exists:

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
Pop-Location
```

`gofmt -l .` lists unformatted files but does not fail for formatting differences;
the acceptance criterion is an empty listing. Use `gofmt -w .` to format.
Use `cargo fmt` for Rust. These commands do not initialize modules for you.
Do not add dependencies, workspaces, CI, or a build system to complete a week.

Optional editor support: rust-analyzer and Go language tooling. A text editor
and terminal are sufficient. Download editor extensions only if you want them.

## Executables and fixtures

Prefer `cargo run` and `go run` while learning. Cargo puts artifacts under each
crate's target directory. For Go, explicitly build to `bin` when requested:
`go build -o .\bin\logscope.exe .\cmd\logscope`. The course ignores bin
directories and `.exe` files. Do not rely on bare `go build` output names.

The shared fixture paths from either analyzer module are
`..\..\fixtures\sample.log`, `empty.log`, and `malformed.log`.
Tests must locate fixtures relative to their source/module, not your personal
absolute path. Rust integration tests can use `CARGO_MANIFEST_DIR`; Go tests
in the module root can use `filepath.Join("..", "..", "fixtures", ...)`.
Fixtures use LF; empty.log is intentionally zero bytes, not a placeholder.
The included .gitattributes preserves LF on Git checkouts. Tests must also
construct CRLF byte inputs in memory.
