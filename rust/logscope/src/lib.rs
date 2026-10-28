//! Evolving weeks 3-6 learner scaffold. See docs/exercises/rust-logscope.md.
//! Run only the current week's test targets until its gate is complete.

pub mod command;
pub mod error;
pub mod model;
pub mod parallel;
pub mod parser;
pub mod reader;
pub mod summary;

pub use error::{AnalyzeError, ParseError};
pub use model::{record, service_name, Level, OwnedRecord, Record};
pub use parallel::analyze_parallel;
pub use parser::parse_line;
pub use reader::analyze_reader;
pub use summary::{analyze, render, ServiceSummary, Summary};
