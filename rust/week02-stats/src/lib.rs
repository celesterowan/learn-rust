//! Week 2 learner scaffold. See docs/exercises/stats.md at the course root.

use std::fmt;

#[derive(Debug, Clone, Copy, PartialEq, Eq)]
pub struct Stats {
    pub count: u64,
    pub sum: u64,
    pub min: Option<u32>,
    pub max: Option<u32>,
}

#[derive(Debug, Clone, Copy, PartialEq, Eq)]
pub enum Dataset {
    Empty,
    Nonempty(usize),
}

#[derive(Debug, Clone, Copy, PartialEq, Eq)]
pub enum StatsError {
    InvalidArgument,
    Overflow,
    Unimplemented(&'static str),
}

impl fmt::Display for StatsError {
    fn fmt(&self, f: &mut fmt::Formatter<'_>) -> fmt::Result {
        match self {
            Self::InvalidArgument => write!(f, "invalid numeric argument"),
            Self::Overflow => write!(f, "arithmetic overflow"),
            Self::Unimplemented(step) => write!(f, "TODO: implement {step}"),
        }
    }
}

impl std::error::Error for StatsError {}

/// Classify a borrowed slice; the nonempty variant carries its length.
pub fn classify(values: &[u32]) -> Dataset {
    // TODO R1: Distinguish an empty slice from one with a nonzero length.
    let _ = values;
    todo!("R1: classify")
}

/// Add without wrapping; return Overflow if the mathematical sum exceeds u64.
pub fn checked_add(accumulator: u64, value: u64) -> Result<u64, StatsError> {
    // TODO R2: Match u64::checked_add explicitly, returning Ok or Overflow.
    let _ = (accumulator, value);
    Err(StatsError::Unimplemented("R2: checked_add"))
}

/// Summarize any u32 values without mutating or taking ownership of the input.
/// Empty input has count/sum zero and absent extrema. Check count and sum.
pub fn calculate(values: &[u32]) -> Result<Stats, StatsError> {
    // TODO R3: Accumulate into u64 with checked_add; handle both Option states.
    let _ = values;
    Err(StatsError::Unimplemented("R3: calculate"))
}

/// Parse one CLI token: 1-4 ASCII digits, numeric value 0..=1000.
/// Leading zeros are allowed. Every other token returns InvalidArgument.
pub fn parse_arg(argument: &str) -> Result<u32, StatsError> {
    // TODO R4: Validate bytes/length before parsing; reject out-of-range values.
    let _ = argument;
    Err(StatsError::Unimplemented("R4: parse_arg"))
}

/// Render a valid Stats as exactly one LF-terminated report, including empty.
pub fn format_report(stats: &Stats) -> Result<String, StatsError> {
    // TODO R5: Render absent extrema as "none", preserving field order and LF.
    let _ = stats;
    Err(StatsError::Unimplemented("R5: format_report"))
}

/// Validate ALL arguments, calculate, then render. Arguments exclude argv[0].
/// Return an error rather than printing, panicking, or returning a partial report.
pub fn run(arguments: &[String]) -> Result<String, StatsError> {
    // TODO R6: Call parse_arg, calculate, and format_report using explicit matches.
    let _ = arguments;
    Err(StatsError::Unimplemented("R6: run"))
}

/// Supplied process adapter; stream contents are testable without a subprocess.
#[derive(Debug, PartialEq, Eq)]
pub struct CommandOutput {
    pub status: i32,
    pub stdout: String,
    pub stderr: String,
}

pub fn command_output(result: Result<String, StatsError>) -> CommandOutput {
    match result {
        Ok(report) => CommandOutput {
            status: 0,
            stdout: report,
            stderr: String::new(),
        },
        Err(error) => CommandOutput {
            status: match error {
                StatsError::InvalidArgument => 2,
                StatsError::Overflow | StatsError::Unimplemented(_) => 1,
            },
            stdout: String::new(),
            stderr: format!("error: {error}\n"),
        },
    }
}
