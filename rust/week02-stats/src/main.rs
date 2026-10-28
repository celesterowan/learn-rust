use std::io::{self, Write};

fn main() {
    // TODO R6: Trace this supplied wiring while implementing run in lib.rs.
    let arguments: Vec<String> = std::env::args().skip(1).collect();
    let output = week02_stats::command_output(week02_stats::run(&arguments));
    // Byte writes keep the report's LF unchanged on Windows.
    let _ = io::stdout().write_all(output.stdout.as_bytes());
    let _ = io::stderr().write_all(output.stderr.as_bytes());
    std::process::exit(output.status);
}
