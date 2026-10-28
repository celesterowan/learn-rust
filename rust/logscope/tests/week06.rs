mod common;

use std::sync::atomic::{AtomicUsize, Ordering};
use std::thread;

use logscope::parallel::join_all;
use logscope::{analyze, analyze_parallel, render, AnalyzeError, ServiceSummary, Summary};

#[test]
fn starter_two_borrowing_threads() {
    // TODO W6-S: Use two scoped threads borrowing fixed complete records.
    // Give each a private count, explicitly join both, and assert their sum.
    todo!("W6-S: two scoped thread experiment")
}

#[test]
fn parallel_equivalence() {
    let mut cases = vec![
        Vec::new(),
        b"\n\r\n\r".to_vec(),
        b"INFO api 1".to_vec(),
        b"INFO api 1\r\nWARN api 2\r\nERROR worker 0\r".to_vec(),
        b"\nINFO z 60000\nbad\n\xff\n\0\nINFO a 00001\nINFO z 1\r\r\n".to_vec(),
    ];
    for name in ["sample", "empty", "malformed"] {
        cases.push(common::bytes(&format!("{name}.log")));
    }
    cases.push(b"INFO api 1\nbad\nERROR worker 2\n".repeat(1000));
    for bytes in cases {
        let sequential = analyze(&bytes).unwrap();
        for workers in 1..=4 {
            let parallel = analyze_parallel(&bytes, workers).unwrap();
            assert_eq!(parallel, sequential, "workers={workers}");
            assert_eq!(render(&parallel).as_bytes(), render(&sequential).as_bytes());
        }
    }
}

#[test]
fn invalid_worker_counts_even_with_no_records() {
    for workers in [0, 5, usize::MAX] {
        for input in [b"".as_slice(), b"INFO api 1"] {
            assert!(matches!(
                analyze_parallel(input, workers),
                Err(AnalyzeError::InvalidWorkers)
            ));
        }
    }
}

fn service(count: u64, total_ms: u64, max_ms: u32) -> ServiceSummary {
    ServiceSummary {
        count,
        total_ms,
        max_ms,
    }
}

#[test]
fn merge_empty_new_and_overlapping_keys() {
    let mut left = analyze(b"INFO z 10\nWARN a 0\nbad").unwrap();
    let right = analyze(b"ERROR z 20\nINFO b 7\nbad").unwrap();
    left.merge(right).unwrap();
    assert_eq!(
        (left.valid, left.invalid, left.info, left.warn, left.error),
        (4, 2, 2, 1, 1)
    );
    assert_eq!(left.services["z"], service(2, 30, 20));
    assert_eq!(left.services["a"], service(1, 0, 0));
    assert_eq!(left.services["b"], service(1, 7, 7));
    let before = left.clone();
    left.merge(Summary::empty()).unwrap();
    assert_eq!(left, before);
    let mut empty = Summary::empty();
    empty.merge(before.clone()).unwrap();
    assert_eq!(empty, before);
}

#[test]
fn merge_checks_every_addition() {
    for field in 0..7 {
        let mut left = Summary::empty();
        let mut right = Summary::empty();
        match field {
            0 => {
                left.valid = u64::MAX;
                right.valid = 1;
            }
            1 => {
                left.invalid = u64::MAX;
                right.invalid = 1;
            }
            2 => {
                left.info = u64::MAX;
                right.info = 1;
            }
            3 => {
                left.warn = u64::MAX;
                right.warn = 1;
            }
            4 => {
                left.error = u64::MAX;
                right.error = 1;
            }
            5 => {
                left.services.insert("api".into(), service(u64::MAX, 0, 0));
                right.services.insert("api".into(), service(1, 0, 0));
            }
            _ => {
                left.services.insert("api".into(), service(1, u64::MAX, 1));
                right.services.insert("api".into(), service(1, 1, 1));
            }
        }
        assert!(
            matches!(left.merge(right), Err(AnalyzeError::Overflow)),
            "field={field}"
        );
    }
}

#[test]
fn merge_accepts_exact_u64_boundary_and_keeps_larger_maximum() {
    let mut left = Summary::empty();
    left.valid = u64::MAX - 1;
    left.invalid = u64::MAX - 1;
    left.info = u64::MAX - 1;
    left.warn = u64::MAX - 1;
    left.error = u64::MAX - 1;
    left.services
        .insert("api".into(), service(u64::MAX - 1, u64::MAX - 1, 60000));
    let mut right = Summary::empty();
    right.valid = 1;
    right.invalid = 1;
    right.info = 1;
    right.warn = 1;
    right.error = 1;
    right.services.insert("api".into(), service(1, 1, 1));
    left.merge(right).unwrap();
    assert_eq!(
        (left.valid, left.invalid, left.info, left.warn, left.error),
        (u64::MAX, u64::MAX, u64::MAX, u64::MAX, u64::MAX)
    );
    assert_eq!(left.services["api"], service(u64::MAX, u64::MAX, 60000));
}

#[test]
fn join_all_success_and_empty() {
    thread::scope(|scope| {
        let handles = vec![
            scope.spawn(|| analyze(b"INFO api 1")),
            scope.spawn(|| analyze(b"ERROR api 2")),
        ];
        let summary = join_all(handles).unwrap();
        assert_eq!(summary, analyze(b"INFO api 1\nERROR api 2").unwrap());
        assert_eq!(join_all(Vec::new()).unwrap(), Summary::empty());
    });
}

#[test]
fn join_all_drains_after_error_and_panic() {
    for panic_first in [false, true] {
        let finished = AtomicUsize::new(0);
        // An unjoined panicking handle also makes scope panic, so this catches
        // accidental early returns even though scope auto-joins dropped handles.
        thread::scope(|scope| {
            let first = scope.spawn(|| {
                if panic_first {
                    panic!("injected first worker");
                }
                Err(AnalyzeError::Overflow)
            });
            let second = scope.spawn(|| -> Result<Summary, AnalyzeError> {
                panic!("injected later worker");
            });
            let last = scope.spawn(|| {
                thread::sleep(std::time::Duration::from_millis(10));
                finished.fetch_add(1, Ordering::SeqCst);
                Ok(Summary::empty())
            });
            let result = join_all(vec![first, second, last]);
            if panic_first {
                assert!(matches!(result, Err(AnalyzeError::WorkerPanicked)));
            } else {
                assert!(matches!(result, Err(AnalyzeError::Overflow)));
            }
            assert_eq!(
                finished.load(Ordering::SeqCst),
                1,
                "must join before returning"
            );
        });
    }
}

#[test]
fn join_all_merge_overflow_still_joins_later_panic() {
    thread::scope(|scope| {
        let first = scope.spawn(|| {
            let mut summary = Summary::empty();
            summary.valid = u64::MAX;
            Ok(summary)
        });
        let second = scope.spawn(|| {
            let mut summary = Summary::empty();
            summary.valid = 1;
            Ok(summary)
        });
        let last = scope.spawn(|| -> Result<Summary, AnalyzeError> {
            panic!("injected after merge overflow");
        });
        assert!(matches!(
            join_all(vec![first, second, last]),
            Err(AnalyzeError::Overflow)
        ));
    });
}
