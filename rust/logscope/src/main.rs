use std::io;

fn main() {
    let arguments: Vec<_> = std::env::args_os().skip(1).collect();
    let status = logscope::command::execute(
        &arguments,
        &mut io::stdout().lock(),
        &mut io::stderr().lock(),
    );
    std::process::exit(status);
}
