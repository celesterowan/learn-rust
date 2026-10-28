use std::thread::ScopedJoinHandle;

use crate::{AnalyzeError, Summary};

/// Library-only extension: workers must be 1..=4, even for empty input.
/// Borrow complete lines in at most workers scoped threads, with private state.
pub fn analyze_parallel(input: &[u8], workers: usize) -> Result<Summary, AnalyzeError> {
    // TODO W6-3: Split on LF without normalizing twice; reuse process_line.
    // Retain every handle and reduce through join_all.
    let _ = (input, workers);
    Err(AnalyzeError::Unimplemented("W6-3: analyze_parallel"))
}

/// Join EVERY handle, even after an error; return no partial summary.
/// First observed failure in handle order wins. A panic becomes WorkerPanicked.
/// An overflow from merge is a failure too, but must not stop later joins.
pub fn join_all<'scope>(
    handles: Vec<ScopedJoinHandle<'scope, Result<Summary, AnalyzeError>>>,
) -> Result<Summary, AnalyzeError> {
    // TODO W6-2: Accumulate privately, remember the first error, drain all joins.
    let _ = handles;
    Err(AnalyzeError::Unimplemented("W6-2: join_all"))
}
