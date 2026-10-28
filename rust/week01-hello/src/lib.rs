//! Week 1 learner scaffold. See docs/exercises/rust-hello.md.

use std::fmt;

#[derive(Debug, Clone, Copy, PartialEq, Eq)]
pub enum HelloError {
    InvalidArguments,
    Unimplemented,
}

impl fmt::Display for HelloError {
    fn fmt(&self, f: &mut fmt::Formatter<'_>) -> fmt::Result {
        match self {
            Self::InvalidArguments => write!(f, "expected zero or one nonempty name"),
            Self::Unimplemented => write!(f, "TODO: implement H2: greeting"),
        }
    }
}

impl std::error::Error for HelloError {}

/// Return "zero" for 0 and "nonzero" for every other u32.
pub fn classify(value: u32) -> &'static str {
    // TODO H1: Use an if expression as this function's final expression.
    let _ = value;
    todo!("H1: classify")
}

/// Arguments exclude argv[0]. Return one LF-terminated greeting, or an error.
/// No arguments uses learner; one nonempty name is preserved without trimming.
pub fn greeting(arguments: &[String]) -> Result<String, HelloError> {
    // TODO H2: Reject empty names and extra arguments; never print or panic here.
    let _ = arguments;
    Err(HelloError::Unimplemented)
}
