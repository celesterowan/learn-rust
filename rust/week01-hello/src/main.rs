use std::io::{self, Write};
use week01_hello::{greeting, HelloError};

fn main() {
    // Supplied adapter: names must be Unicode, unlike an OS file path.
    let arguments: Result<Vec<String>, _> = std::env::args_os()
        .skip(1)
        .map(|s| s.into_string())
        .collect();
    let result = match arguments {
        Ok(arguments) => greeting(&arguments),
        Err(_) => Err(HelloError::InvalidArguments),
    };
    let status = match result {
        Ok(report) => match io::stdout().write_all(report.as_bytes()) {
            Ok(()) => 0,
            Err(error) => {
                let _ = writeln!(io::stderr(), "error: {error}");
                1
            }
        },
        Err(error) => {
            let _ = writeln!(io::stderr(), "error: {error}");
            match error {
                HelloError::InvalidArguments => 2,
                HelloError::Unimplemented => 1,
            }
        }
    };
    std::process::exit(status);
}
