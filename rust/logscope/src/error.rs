use std::{fmt, io};

#[derive(Debug, Clone, Copy, PartialEq, Eq)]
pub enum ParseError {
    Malformed,
    Unimplemented,
}

impl fmt::Display for ParseError {
    fn fmt(&self, f: &mut fmt::Formatter<'_>) -> fmt::Result {
        match self {
            Self::Malformed => write!(f, "malformed record"),
            Self::Unimplemented => write!(f, "TODO: implement W4-1: parse_line"),
        }
    }
}

impl std::error::Error for ParseError {}

#[derive(Debug)]
pub enum AnalyzeError {
    Overflow,
    Io(io::Error),
    InvalidWorkers,
    WorkerPanicked,
    Unimplemented(&'static str),
}

impl fmt::Display for AnalyzeError {
    fn fmt(&self, f: &mut fmt::Formatter<'_>) -> fmt::Result {
        match self {
            Self::Overflow => write!(f, "arithmetic overflow"),
            Self::Io(error) => write!(f, "I/O: {error}"),
            Self::InvalidWorkers => write!(f, "workers must be in 1..=4"),
            Self::WorkerPanicked => write!(f, "worker panicked"),
            Self::Unimplemented(step) => write!(f, "TODO: implement {step}"),
        }
    }
}

impl std::error::Error for AnalyzeError {
    fn source(&self) -> Option<&(dyn std::error::Error + 'static)> {
        match self {
            Self::Io(error) => Some(error),
            _ => None,
        }
    }
}

impl From<io::Error> for AnalyzeError {
    fn from(error: io::Error) -> Self {
        // TODO W5-2: Wrap, do not stringify or replace, the original error.
        let _ = error;
        Self::Unimplemented("W5-2: From<io::Error>")
    }
}
