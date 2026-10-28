mod common;

use logscope::summary::{checked_add, process_line};
use logscope::{analyze, render, AnalyzeError, Level, Record, ServiceSummary, Summary};

#[test]
fn fixture_bytes() {
    for name in ["sample", "empty", "malformed"] {
        let summary = analyze(&common::bytes(&format!("{name}.log"))).unwrap();
        assert_eq!(
            render(&summary).as_bytes(),
            common::bytes(&format!("{name}.expected.txt")),
            "{name}"
        );
    }
}

#[test]
fn normalization_and_invalid_lines() {
    for (input, valid, invalid) in [
        (b"".as_slice(), 0, 0),
        (b"\n\r\n\r", 0, 0),
        (b"\r\r\n", 0, 1),
        (b"INFO api 1", 1, 0),
        (b"INFO api 1\n", 1, 0),
        (b"INFO api 1\r", 1, 0),
        (b"INFO api 1\r\nWARN api 2\r\n", 2, 0),
        (b"\nINFO api 1\n\nbad\n \n\xff\n\0", 1, 4),
        (b"INFO api 1\r\r\nINFO api 2\r\r", 0, 2),
    ] {
        let summary = analyze(input).unwrap();
        assert_eq!(
            (summary.valid, summary.invalid),
            (valid, invalid),
            "{input:?}"
        );
    }
}

#[test]
fn process_line_is_reusable() {
    let mut summary = Summary::empty();
    for line in [b"\r".as_slice(), b"INFO api 0\r", b"bad", b"ERROR worker 7"] {
        process_line(&mut summary, line).unwrap();
    }
    assert_eq!((summary.valid, summary.invalid), (2, 1));
    assert_eq!((summary.info, summary.warn, summary.error), (1, 0, 1));
    assert_eq!(summary.services["api"].max_ms, 0);
}

#[test]
fn owned_keys_and_sorted_complete_report() {
    let summary = {
        let input = b"ERROR z 2\nWARN a 0\nINFO z 5\nINFO z 1\n".to_vec();
        analyze(&input).unwrap()
    };
    assert_eq!(
        render(&summary),
        "valid=4 invalid=0\nlevel INFO=2 WARN=1 ERROR=1\n\
service a count=1 total_ms=0 max_ms=0\n\
service z count=3 total_ms=8 max_ms=5\n"
    );
    assert_eq!(
        render(&Summary::empty()),
        "valid=0 invalid=0\nlevel INFO=0 WARN=0 ERROR=0\n"
    );
}

fn overflow<T: std::fmt::Debug>(result: Result<T, AnalyzeError>) {
    assert!(matches!(result, Err(AnalyzeError::Overflow)), "{result:?}");
}

#[test]
fn checked_add_boundaries() {
    for (a, b, sum) in [
        (0, 0, 0),
        (u64::MAX, 0, u64::MAX),
        (u64::MAX - 1, 1, u64::MAX),
    ] {
        assert_eq!(checked_add(a, b).unwrap(), sum);
    }
    overflow(checked_add(u64::MAX, 1));
    overflow(checked_add(u64::MAX - 1, 2));
}

#[test]
fn every_counter_and_total_checks_overflow() {
    let record = Record {
        level: Level::Info,
        service: "api",
        duration: 1,
    };
    let mut summary = Summary::empty();
    summary.invalid = u64::MAX;
    overflow(summary.add_invalid());
    for field in 0..5 {
        let mut summary = Summary::empty();
        let level = match field {
            2 => Level::Warn,
            3 => Level::Error,
            _ => Level::Info,
        };
        match field {
            0 => summary.valid = u64::MAX,
            1 => summary.info = u64::MAX,
            2 => summary.warn = u64::MAX,
            3 => summary.error = u64::MAX,
            _ => {
                summary.services.insert(
                    "api".into(),
                    ServiceSummary {
                        count: u64::MAX,
                        total_ms: 0,
                        max_ms: 0,
                    },
                );
            }
        }
        overflow(summary.add_record(&Record { level, ..record }));
    }
    let mut summary = Summary::empty();
    summary.services.insert(
        "api".into(),
        ServiceSummary {
            count: 1,
            total_ms: u64::MAX,
            max_ms: 1,
        },
    );
    overflow(summary.add_record(&record));
}

#[test]
fn wide_total_and_invalid_updates() {
    let mut summary = Summary::empty();
    let record = Record {
        level: Level::Info,
        service: "api",
        duration: u32::MAX,
    };
    summary.add_record(&record).unwrap();
    summary.add_record(&record).unwrap();
    summary.add_invalid().unwrap();
    assert_eq!(summary.services["api"].total_ms, 2 * u64::from(u32::MAX));
    assert_eq!(summary.services["api"].max_ms, u32::MAX);
    assert_eq!((summary.valid, summary.invalid, summary.info), (2, 1, 2));
}
