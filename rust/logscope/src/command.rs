//! Supplied resource/exit adapter; learner analysis stays in the library.

use std::ffi::OsString;
use std::io::Write;
use std::path::Path;

use crate::{analyze, render, AnalyzeError};

/// Week 4 byte-oriented file boundary. Week 5 replaces this body with
/// File::open + analyze_reader, preserving the public command contract.
pub fn report_path(path: &Path) -> Result<String, AnalyzeError> {
    // TODO W5-5: Open File and pass it to analyze_reader; use ? for I/O.
    let bytes = std::fs::read(path).map_err(AnalyzeError::Io)?;
    let summary = analyze(&bytes)?;
    Ok(render(&summary))
}

/// Arguments exclude argv[0]; OS strings preserve non-Unicode file paths.
/// stdout failures can leave partial output. Stderr is best-effort.
pub fn execute(arguments: &[OsString], stdout: &mut impl Write, stderr: &mut impl Write) -> i32 {
    if arguments.len() != 1 {
        let _ = writeln!(stderr, "error: expected exactly one file path");
        return 2;
    }
    match report_path(Path::new(&arguments[0])) {
        Ok(report) => match stdout
            .write_all(report.as_bytes())
            .and_then(|()| stdout.flush())
        {
            Ok(()) => 0,
            Err(error) => {
                let _ = writeln!(stderr, "error: {error}");
                1
            }
        },
        Err(error) => {
            let _ = writeln!(stderr, "error: {error}");
            1
        }
    }
}
