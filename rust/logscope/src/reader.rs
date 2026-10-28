use std::io::Read;

use crate::{AnalyzeError, Summary};

/// Read to EOF, then analyze all bytes. A read failure discards any prefix.
/// Own R, which can itself be File, Cursor, or a borrowed reader.
pub fn analyze_reader<R: Read>(reader: R) -> Result<Summary, AnalyzeError> {
    // TODO W5-3: Read bytes, use From with ?, and reuse analyze.
    let _ = reader;
    Err(AnalyzeError::Unimplemented("W5-3: analyze_reader"))
}
