use crate::{ParseError, Record};

/// Parse one normalized nonempty ASCII line, borrowing its service bytes.
/// Do not split LF, strip CR, trim, allocate service storage, or count errors.
pub fn parse_line(line: &[u8]) -> Result<Record<'_>, ParseError> {
    // TODO W4-1: Validate the complete grammar in docs/logscope-contract.md.
    let _ = line;
    Err(ParseError::Unimplemented)
}
