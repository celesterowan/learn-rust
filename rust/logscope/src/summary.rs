use std::collections::BTreeMap;

use crate::{AnalyzeError, Record};

#[derive(Debug, Clone, PartialEq, Eq)]
pub struct ServiceSummary {
    pub count: u64,
    pub total_ms: u64,
    pub max_ms: u32,
}

#[derive(Debug, Clone, PartialEq, Eq)]
pub struct Summary {
    pub valid: u64,
    pub invalid: u64,
    pub info: u64,
    pub warn: u64,
    pub error: u64,
    pub services: BTreeMap<String, ServiceSummary>,
}

impl Summary {
    /// Week 4 empty state; do not depend on the week 5 Default TODO yet.
    pub fn empty() -> Self {
        Self {
            valid: 0,
            invalid: 0,
            info: 0,
            warn: 0,
            error: 0,
            services: BTreeMap::new(),
        }
    }

    /// Count one valid record, checking every affected count and sum.
    /// On error the caller discards the summary; rollback is not required.
    pub fn add_record(&mut self, record: &Record<'_>) -> Result<(), AnalyzeError> {
        // TODO W4-3: Own service keys; preserve the maximum duration.
        let _ = record;
        Err(AnalyzeError::Unimplemented("W4-3: Summary::add_record"))
    }

    /// Count one malformed nonempty line, with checked arithmetic.
    pub fn add_invalid(&mut self) -> Result<(), AnalyzeError> {
        // TODO W4-3: Check the invalid counter before incrementing.
        Err(AnalyzeError::Unimplemented("W4-3: Summary::add_invalid"))
    }

    /// Merge all counters and services; check additions, take maximum max_ms.
    /// On failure both operands are discarded; no partial-result API is promised.
    pub fn merge(&mut self, other: Summary) -> Result<(), AnalyzeError> {
        // TODO W6-1: Handle both new and overlapping owned service keys.
        let _ = other;
        Err(AnalyzeError::Unimplemented("W6-1: Summary::merge"))
    }
}

impl Default for Summary {
    fn default() -> Self {
        // TODO W5-1: Delegate to the existing empty state.
        todo!("W5-1: Summary::default")
    }
}

/// Return the mathematical sum or Overflow, never wrap.
pub fn checked_add(left: u64, right: u64) -> Result<u64, AnalyzeError> {
    // TODO W4-2: Use checked arithmetic; tests inject near-limit values.
    let _ = (left, right);
    Err(AnalyzeError::Unimplemented("W4-2: checked_add"))
}

/// Process one raw LF-split segment: remove at most one final CR, ignore empty,
/// otherwise parse and update valid or invalid. Never swallow a TODO as input.
pub fn process_line(summary: &mut Summary, line: &[u8]) -> Result<(), AnalyzeError> {
    // TODO W4-4: Count only ParseError::Malformed; propagate Unimplemented.
    let _ = (summary, line);
    Err(AnalyzeError::Unimplemented("W4-4: process_line"))
}

/// Analyze bytes, not a decoded whole-file string. Return an owned summary.
pub fn analyze(input: &[u8]) -> Result<Summary, AnalyzeError> {
    // TODO W4-5: Split on LF and reuse process_line, including the final segment.
    let _ = input;
    Err(AnalyzeError::Unimplemented("W4-5: analyze"))
}

/// Render the complete frozen report, sorted by ASCII service name, with LF.
pub fn render(summary: &Summary) -> String {
    // TODO W4-6: Include all three levels and exactly one final LF.
    // TODO W5-4: Later refactor one collection loop to map/collect.
    let _ = summary;
    todo!("W4-6: render")
}
