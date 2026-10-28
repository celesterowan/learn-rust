mod common;

use std::ffi::OsString;
use std::io::{self, Write};
use std::path::PathBuf;
use std::process::Command;
use std::sync::atomic::{AtomicU64, Ordering};

use logscope::command::execute;

struct TempDirectory(PathBuf);

impl TempDirectory {
    fn new() -> Self {
        static NEXT: AtomicU64 = AtomicU64::new(0);
        let path = std::env::temp_dir().join(format!(
            "rust-logscope-{}-{}-{}",
            std::process::id(),
            std::time::SystemTime::now()
                .duration_since(std::time::UNIX_EPOCH)
                .unwrap()
                .as_nanos(),
            NEXT.fetch_add(1, Ordering::Relaxed)
        ));
        std::fs::create_dir(&path).unwrap();
        Self(path)
    }
}

impl Drop for TempDirectory {
    fn drop(&mut self) {
        let _ = std::fs::remove_dir_all(&self.0);
    }
}

#[test]
fn compiled_binary_fixtures() {
    for name in ["sample", "empty", "malformed"] {
        let output = Command::new(env!("CARGO_BIN_EXE_logscope"))
            .arg(common::fixture(&format!("{name}.log")))
            .output()
            .unwrap();
        assert_eq!(output.status.code(), Some(0), "{output:?}");
        assert_eq!(
            output.stdout,
            common::bytes(&format!("{name}.expected.txt"))
        );
        assert!(output.stderr.is_empty());
    }
}

#[test]
fn usage_and_file_errors() {
    for arguments in [vec![], vec![OsString::from("one"), OsString::from("two")]] {
        let output = Command::new(env!("CARGO_BIN_EXE_logscope"))
            .args(arguments)
            .output()
            .unwrap();
        assert_eq!(output.status.code(), Some(2));
        assert!(output.stdout.is_empty());
        assert!(output.stderr.starts_with(b"error:"));
    }
    let directory = TempDirectory::new();
    for path in [directory.0.join("missing.log"), directory.0.clone()] {
        let output = Command::new(env!("CARGO_BIN_EXE_logscope"))
            .arg(path)
            .output()
            .unwrap();
        assert_eq!(output.status.code(), Some(1));
        assert!(output.stdout.is_empty());
        assert!(output.stderr.starts_with(b"error:"));
    }
}

#[test]
fn spaces_dash_and_invalid_utf8_are_file_input() {
    let directory = TempDirectory::new();
    for filename in ["input with spaces.log", "-input.log"] {
        std::fs::write(directory.0.join(filename), b"INFO api 1\n\xff\n").unwrap();
        let output = Command::new(env!("CARGO_BIN_EXE_logscope"))
            .current_dir(&directory.0)
            .arg(filename)
            .output()
            .unwrap();
        assert_eq!(output.status.code(), Some(0), "{output:?}");
        assert_eq!(output.stdout, b"valid=1 invalid=1\nlevel INFO=1 WARN=0 ERROR=0\nservice api count=1 total_ms=1 max_ms=1\n");
        assert!(output.stderr.is_empty());
    }
}

struct FailingWriter {
    fail_flush: bool,
}

impl Write for FailingWriter {
    fn write(&mut self, bytes: &[u8]) -> io::Result<usize> {
        if self.fail_flush {
            Ok(bytes.len())
        } else {
            Err(io::Error::other("write injected"))
        }
    }
    fn flush(&mut self) -> io::Result<()> {
        Err(io::Error::other("flush injected"))
    }
}

#[test]
fn output_failures_exit_one_and_stderr_is_best_effort() {
    let args = [common::fixture("sample.log").into_os_string()];
    for fail_flush in [false, true] {
        let mut diagnostic = Vec::new();
        assert_eq!(
            execute(&args, &mut FailingWriter { fail_flush }, &mut diagnostic),
            1
        );
        assert!(diagnostic.starts_with(b"error:"));
        // Distinguish a real output failure from an unfinished analyzer error.
        assert!(String::from_utf8(diagnostic).unwrap().contains("injected"));
    }
    let mut output = Vec::new();
    assert_eq!(
        execute(&[], &mut output, &mut FailingWriter { fail_flush: false }),
        2
    );
    assert!(output.is_empty());
}
