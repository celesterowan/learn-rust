# Setup: Windows and macOS

Follow only your platform's section: [Windows PowerShell](#windows-powershell)
or [macOS Terminal](#macos-terminal). These are instructions for you to run,
not tools already installed by the course. Verify existing installations before
downloading another copy. Installers may request administrator permission;
ordinary course commands should run as your normal user.

## What you need to download

| Item | Windows | macOS | When needed |
| --- | --- | --- | --- |
| Git | Git for Windows | Included with Apple's Command Line Tools or Xcode | To clone and later push |
| Native compiler/linker | MSVC C++ tools and Windows SDK | Apple's Command Line Tools or Xcode | Before compiling Rust |
| Rust | Stable rustup installer; includes Cargo | Stable rustup installer; includes Cargo | Week 1 |
| Go | Stable Windows MSI | Stable macOS PKG for your processor | By week 7 |
| Course materials | Clone this repository or extract its ZIP | Same | Before week 1 |
| Editor | Any plain-text/code editor | Any plain-text/code editor | For writing your work |

You need internet access for these downloads and online readings. Rust's local
reference documentation is available through `rustup doc` after installation.
The graded exercises use only the standard libraries: there is no package
bundle, completed solution, framework, or Microsoft RustTraining checkout to
download. All required projects are included as buildable scaffolds with learner
TODOs and active acceptance tests. See the
[project guide index](README.md#project-specs-and-implementation-guides).
A GitHub account is needed to own a fork and push, not to read or
download this public course. GitHub CLI, Homebrew, full IDEs, and editor
extensions are optional, not prerequisites.

## Windows PowerShell

Use PowerShell, not Command Prompt or Git Bash, for the Windows examples.
Check **Settings > System > About > System type** before selecting installers.
Use x64/amd64 downloads on x64 Windows and native ARM64 downloads on ARM64
Windows where supported. Do not mix architectures to work around errors.

1. Install [Git for Windows](https://git-scm.com/install/windows), if
   `git --version` is not already available. Keep Git accessible from
   PowerShell through the installer's PATH option.
2. Check your existing Visual Studio installation. Rust's MSVC toolchain needs
   the **Desktop development with C++** workload, including the MSVC tools for
   your architecture and a Windows SDK. Use the
   [official MSVC prerequisite guide](https://rust-lang.github.io/rustup/installation/windows-msvc.html)
   and an appropriately licensed Visual Studio or
   [Build Tools](https://visualstudio.microsoft.com/visual-cpp-build-tools/)
   installation. You do not need to reinstall tools your C++ environment
   already provides.
3. Download and run [rustup-init](https://www.rust-lang.org/tools/install).
   Choose the stable MSVC toolchain for your machine. It installs `rustup`,
   `rustc`, and `cargo`; Cargo is not a separate download.
4. By week 7, download the current stable **Windows MSI** from
   [Go downloads](https://go.dev/dl/) for your architecture. Follow the
   [Windows installation tab](https://go.dev/doc/install), then close and
   reopen PowerShell.

In the new PowerShell window, verify the tools. Add Rust's formatting and
linting components if not already installed:

```powershell
git --version
rustup show
rustc --version
cargo --version
rustup component add rustfmt clippy
go version
go env GOOS GOARCH
```

If you deferred Go, defer its two verification commands too. Record actual
versions in your reflections rather than matching old tutorial version numbers.
Proceed to [download the course](#download-the-course).

## macOS Terminal

Use Terminal's normal zsh shell. Open **Apple menu > About This Mac**:
an Apple chip uses arm64 downloads; an Intel processor uses amd64/x86_64.
Use native Terminal rather than an Intel/Rosetta session on Apple silicon.
`uname -m` reports the running environment (`arm64` or `x86_64`); compare it
with About This Mac if unsure.

1. If Git and the native toolchain are not installed, request Apple's
   [Command Line Tools](https://developer.apple.com/documentation/xcode/installing-the-command-line-tools):

   ```sh
   xcode-select --install
   ```

   Finish the installation dialog before continuing. The package supplies
   Git, Clang, a linker, and the macOS SDK. Full Xcode is not required; if
   Xcode already supplies working tools, no separate installation is needed.

2. Follow the [official Rust installation instructions](https://doc.rust-lang.org/book/ch01-01-installation.html).
   This official command downloads and executes the rustup installer:

   ```sh
   curl --proto '=https' --tlsv1.2 https://sh.rustup.rs -sSf | sh
   ```

   Review the installer prompts and select the native stable toolchain.
   It installs Cargo along with Rust. Open a new Terminal window afterward;
   alternatively, load its environment in the current shell:

   ```sh
   . "$HOME/.cargo/env"
   ```

3. By week 7, download the current stable **macOS PKG** from
   [Go downloads](https://go.dev/dl/): `darwin-arm64` for Apple silicon or
   `darwin-amd64` for Intel. Open the package and follow the
   [Mac installation tab](https://go.dev/doc/install). The installer adds
   `/usr/local/go/bin` to PATH through `/etc/paths.d/go`; open a new Terminal
   window before verifying it. Homebrew is not required.

Verify the tools and selected architecture:

```sh
uname -m
xcode-select -p
clang --version
git --version
rustup show
rustc --version
cargo --version
rustup component add rustfmt clippy
go version
go env GOOS GOARCH
```

If you deferred Go, defer its two verification commands too. Compiling the
week 1 project checks the linker as well as PATH; a version string alone does
not prove a working build. Record actual tool versions in your reflection.

## Download the course

If you already have a checkout containing this README, SETUP.md, assignments,
and fixtures, use it. Do not clone another repository inside it or inside
Chromium. The **course root** is the directory containing those files, not
the parent Projects directory.

For a new checkout, choose a user-writable directory. The examples use
`Projects` under your home directory and create `learn-rust` inside it.
The URL downloads the main course repository. If you intend to push your
assignments, first fork it on GitHub under an account you can use, then replace
the clone URL with **your fork's Code > HTTPS URL**. Cloning creates the local
checkout and an `origin` remote; it does not create a fork or grant write access.

**Windows PowerShell**

```powershell
New-Item -ItemType Directory -Force "$HOME\Projects" -ErrorAction Stop
Set-Location "$HOME\Projects" -ErrorAction Stop
git clone https://github.com/celesterowan/learn-rust.git learn-rust
if ($LASTEXITCODE -ne 0) { throw 'Clone failed; resolve the error before continuing.' }
Set-Location .\learn-rust -ErrorAction Stop
```

**macOS Terminal**

```sh
mkdir -p "$HOME/Projects" &&
cd "$HOME/Projects" &&
git clone https://github.com/celesterowan/learn-rust.git learn-rust &&
cd learn-rust
```

If `learn-rust` already exists, inspect/use that checkout or choose a different
destination; do not delete it to make cloning work. Confirm that the downloaded
branch contains the handouts and fixtures. If not, select the published course
branch/version on GitHub rather than creating missing handouts yourself.

**Without Git:** on the repository's GitHub page, select the course branch,
then **Code > Download ZIP**. On Windows, choose **Extract All** in File
Explorer. On macOS, double-click the ZIP in Finder. Move the fully extracted
folder to your chosen Projects directory and open the folder containing
README.md, not the ZIP or an extra wrapper folder. ZIP contents have no Git
history or remotes. Prefer a clone if you want to commit and push later.

Read [week 1](assignments/week-01.md) and its
[greeting spec](docs/exercises/rust-hello.md) next. The download includes actual
source scaffolds, tests, assignment specifications, and fixtures. Implement
the learner TODOs; do not recreate or overwrite these projects.

## Shell conventions and directory plan

Run only your platform's command block. All assignment blocks start at the
course root unless they explicitly say otherwise. Windows blocks use
Push-Location/Pop-Location; macOS blocks use a subshell `(cd ... && ...)` so
your parent shell stays at the course root. Do not paste PowerShell paths or
commands into zsh. Portable Cargo/Go commands are marked as shared.

Read each block before running it. Open the supplied project, implement the
assignment, then run its focused build/test commands. Stop at a failing command and
record the error; the last successful command is not proof that earlier
commands passed. Do not run course builds with sudo or as administrator.

```text
learn-rust/
  assignments/         provided handouts
  docs/submissions/    create your own reflections
  fixtures/            provided inputs and expected outputs
  rust/
    week01-hello/      independent Cargo binary crate
    week02-stats/      independent Cargo binary crate
    logscope/          staged library/CLI/worker TODOs, weeks 3-6
  go/
    stats/             independent Go module, week 7
    logscope/          Go module, weeks 8-12
      cmd/logscope/    sequential CLI
      cmd/logserve/    supplied HTTP command, TODOs due week 11
```

There is no root Cargo workspace or Go module. Each Rust project owns its
Cargo.toml/Cargo.lock; each Go project owns its go.mod. Keep Cargo.lock for
these application exercises. Without external Go packages, go.sum may not
exist; that is normal. Do not run cargo new or go mod init over any supplied
project. Do not set GOROOT/GOPATH or put the course inside a toolchain.

Use the [project specs](README.md#project-specs-and-implementation-guides) for
build-only, test-compilation, and per-week acceptance commands. Compilation
should pass before implementation; acceptance failures from explicit TODO
stubs are expected until you complete them. The logscope projects evolve
across several weeks: use focused selectors until the final stage.
The capstone harness has an explicit `parity` build tag and requires both
toolchains; ordinary earlier Go tests do not require Rust.

Microsoft's [getting-started chapter](https://microsoft.github.io/RustTraining/c-cpp-book/ch02-getting-started.html)
supplies the main/Cargo explanation. Its optional REPL, external crates,
workspaces, and publishing exercises are not prerequisites here.

Create your submission directory once, from the course root.

**Windows PowerShell**

```powershell
New-Item -ItemType Directory -Force .\docs\submissions
```

**macOS Terminal**

```sh
mkdir -p docs/submissions
```

During a staged assignment, use its spec's focused acceptance selectors.
After all stages in a crate/module are due, run the full checks below.
From inside a Rust crate, these commands work in both shells:

```text
cargo fmt --check
cargo test
cargo clippy --all-targets -- -D warnings
```

From inside a Go module, in either shell:

```text
gofmt -l .
go vet ./...
go test ./...
```

For the week 12 cross-language gate, also run
`go test -tags parity -count=1 -run '^TestParity$' .` inside `go/logscope`.
See the [capstone spec](docs/exercises/parity.md) for compile-only commands
and the additional toolchain requirement.

`gofmt -l .` lists unformatted files but does not fail for formatting
differences; the acceptance criterion is an empty listing. Use `gofmt -w .`
to format Go and `cargo fmt` for Rust. Do not add dependencies or a build
system. A code editor with rust-analyzer and Go support is optional.

## Race detector prerequisites

For concurrent Go exercises, follow the
[official race detector requirements](https://go.dev/doc/articles/race_detector.html#requirements).
Check configuration in either shell:

```text
go env GOOS GOARCH CGO_ENABLED CC
```

The race detector requires cgo enabled. It supports `darwin/amd64` and
`darwin/arm64` on macOS, and `windows/amd64` on Windows. Native Windows ARM64
race testing is not listed as supported; do not claim an unavailable run passed.
Apple's tools already satisfy this course's macOS native-toolchain setup;
the race guide does not require a separate non-Apple compiler on macOS.

On Windows/amd64, race builds require a compatible C compiler incorporating
mingw-w64 runtime libraries version 8 or later. MSVC alone is not that compiler.
Follow the official requirements for your chosen compiler rather than changing
Rust's toolchain to fix Go. If cgo was explicitly disabled, resolve that setting
and compiler availability before retrying; enabling it alone cannot install
missing tools.

After setup, run `go test -race ./...` inside the assigned Go module. Record
unavailable checks as incomplete and retry in a supported environment you
control. Ordinary tests do not establish race freedom; the race detector
only observes executed paths.

## Executables, exit status, and fixtures

Use `cargo run` and `go run` for normal practice. To check a program's own exit
code, build and invoke the executable directly. In particular, `go run` can
report a wrapper exit code instead of the program's code.

From the course root after implementing the Go analyzer:

**Windows PowerShell**

```powershell
Push-Location .\go\logscope
New-Item -ItemType Directory -Force .\bin
go build -o .\bin\logscope.exe .\cmd\logscope
if ($LASTEXITCODE -ne 0) { Pop-Location; throw 'Build failed.' }
.\bin\logscope.exe ..\..\fixtures\sample.log
$LASTEXITCODE
Pop-Location
```

**macOS Terminal**

```sh
(
  cd go/logscope &&
  mkdir -p bin &&
  go build -o bin/logscope ./cmd/logscope || exit
  ./bin/logscope ../../fixtures/sample.log
  printf 'exit=%s\n' "$?"
)
```

Inspect `$LASTEXITCODE` on Windows or `$?` on macOS immediately after the
executable, before another command changes it. For intentional failure tests,
do not use `&&` between the invocation and its status display. The example
above should report zero. A status display is evidence, not an automated test.
Both platforms' Go binaries stay under ignored bin directories; do not use
bare `go build` and accidentally commit a macOS executable at the module root.

Fixture paths from either analyzer module are `..\..\fixtures\sample.log` on
Windows and `../../fixtures/sample.log` on macOS. The other fixtures are
empty.log and malformed.log in the same directory. Tests must build paths
portably, not embed one machine's absolute paths or Windows separators:
Rust can use `Path::join` with `CARGO_MANIFEST_DIR`; Go can use
`filepath.Join("..", "..", "fixtures", "sample.log")`.

Rust CLI integration tests use Cargo's `CARGO_BIN_EXE_logscope` path, which
includes the platform's executable suffix. Go process tests select `.exe`
only on Windows (`runtime.GOOS == "windows"`); macOS executables have no suffix.
These are filename conventions, not a reason to change expected output.

Fixtures and expected reports use LF on both platforms; .gitattributes preserves
it. empty.log is deliberately zero bytes. Do not rewrite expected files through
shell redirection. Tests must also construct CRLF and invalid-byte input in
memory, and compare output bytes without normalizing newlines.

## Troubleshooting

- **Command not found:** reopen the terminal after installation. In PowerShell,
  `Get-Command git, rustup, cargo, go` shows resolved tools. On macOS,
  `command -v git rustup cargo go` does the same. Do not add another tool
  installation until you know which one PATH selects.
- **Rust PATH:** Windows rustup normally uses `%USERPROFILE%\.cargo\bin`.
  On macOS, reopen Terminal or source `$HOME/.cargo/env`. Follow rustup's
  installation output rather than replacing the entire PATH.
- **Go PATH:** the Windows MSI updates it; the macOS PKG adds
  `/usr/local/go/bin` through `/etc/paths.d/go`. Restart your terminal.
  GOPATH/bin is only needed for additional `go install` tools, not this course.
- **Rust linker errors:** on Windows, check the C++ workload and Windows SDK
  using the MSVC guide; an appropriate Visual Studio Developer PowerShell may
  be needed. On macOS, check `xcode-select -p` and `clang --version`, finish
  any pending tools installation, and check Software Update after OS upgrades.
- **Architecture mismatch:** compare OS hardware information with `rustup show`
  and `go env GOOS GOARCH`; use the native installers. Do not silently switch
  to Rosetta or a different target to hide a build error.
- **Cannot push:** downloading a public repository does not grant write access.
  Check `git remote -v` and the authenticated account. Work locally until you
  can authenticate to your own fork; do not force-push to solve permissions.
