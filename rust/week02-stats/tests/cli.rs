use std::process::Command;
use week02_stats::{command_output, format_report, parse_arg, run, Stats, StatsError};

const INVALID: &[&str] = &[
    "-1", "1001", "1.0", "", " 1", "1 ", "+1", "\t1", "1\t", "1\n", "10000", "00000", "\u{0661}",
    "\u{ff11}", "1a",
];

fn valid_cases() -> Vec<(Vec<&'static str>, &'static str)> {
    vec![
        (vec![], "count=0 sum=0 min=none max=none\n"),
        (vec!["4", "0", "9", "4"], "count=4 sum=17 min=0 max=9\n"),
        (vec!["0"], "count=1 sum=0 min=0 max=0\n"),
        (vec!["1000"], "count=1 sum=1000 min=1000 max=1000\n"),
        (vec!["0001"], "count=1 sum=1 min=1 max=1\n"),
        (
            vec!["0000", "0010", "0999"],
            "count=3 sum=1009 min=0 max=999\n",
        ),
    ]
}

#[test]
fn parse_grammar() {
    for (argument, expected) in [
        ("0", 0),
        ("1000", 1000),
        ("0001", 1),
        ("0000", 0),
        ("0010", 10),
        ("999", 999),
        ("0999", 999),
    ] {
        assert_eq!(parse_arg(argument), Ok(expected), "{argument:?}");
    }
    for argument in INVALID {
        assert_eq!(
            parse_arg(argument),
            Err(StatsError::InvalidArgument),
            "{argument:?}"
        );
    }
}

#[test]
fn report_bytes() {
    for (stats, expected) in [
        (
            Stats {
                count: 0,
                sum: 0,
                min: None,
                max: None,
            },
            "count=0 sum=0 min=none max=none\n",
        ),
        (
            Stats {
                count: 4,
                sum: 17,
                min: Some(0),
                max: Some(9),
            },
            "count=4 sum=17 min=0 max=9\n",
        ),
        (
            Stats {
                count: 2,
                sum: 8_589_934_590,
                min: Some(u32::MAX),
                max: Some(u32::MAX),
            },
            "count=2 sum=8589934590 min=4294967295 max=4294967295\n",
        ),
    ] {
        assert_eq!(
            format_report(&stats).unwrap().as_bytes(),
            expected.as_bytes()
        );
    }
}

#[test]
fn runner_streams() {
    for (arguments, expected) in valid_cases() {
        let arguments: Vec<String> = arguments.iter().map(|value| value.to_string()).collect();
        let output = command_output(run(&arguments));
        assert_eq!(output.status, 0, "{arguments:?}: {}", output.stderr);
        assert_eq!(output.stdout.as_bytes(), expected.as_bytes());
        assert!(output.stderr.is_empty());
    }
    for argument in INVALID {
        let arguments = vec!["4".to_string(), argument.to_string(), "9".to_string()];
        let output = command_output(run(&arguments));
        assert_eq!(output.status, 2, "{argument:?}");
        assert!(output.stdout.is_empty());
        assert!(output.stderr.starts_with("error:"));
    }
}

#[test]
fn overflow_adapter() {
    // This tests the real adapter, not an unreachable oversized CLI invocation.
    let output = command_output(Err(StatsError::Overflow));
    assert_eq!(output.status, 1);
    assert!(output.stdout.is_empty());
    assert!(output.stderr.starts_with("error:"));
}

#[test]
fn executable_streams_and_status() {
    let binary = env!("CARGO_BIN_EXE_week02-stats");
    for (arguments, expected) in valid_cases() {
        let output = Command::new(binary).args(&arguments).output().unwrap();
        assert_eq!(output.status.code(), Some(0), "{arguments:?}: {output:?}");
        assert_eq!(output.stdout, expected.as_bytes());
        assert!(output.stderr.is_empty());
    }
    for argument in INVALID {
        let output = Command::new(binary)
            .args(["4", argument, "9"])
            .output()
            .unwrap();
        assert_eq!(output.status.code(), Some(2), "{argument:?}: {output:?}");
        assert!(output.stdout.is_empty(), "{argument:?}");
        assert!(output.stderr.starts_with(b"error:"), "{argument:?}");
    }
}
