mod common;

use std::error::Error;
use std::io::{self, Cursor, Read};

use logscope::{analyze, analyze_reader, render, AnalyzeError, Summary};

#[test]
fn default_is_empty_report() {
    let summary = Summary::default();
    assert_eq!(summary, Summary::empty());
    assert_eq!(
        render(&summary).as_bytes(),
        common::bytes("empty.expected.txt")
    );
}

fn assert_io(error: AnalyzeError) {
    assert!(!error.to_string().is_empty());
    assert!(
        error.source().is_some(),
        "original cause must remain inspectable"
    );
    match error {
        AnalyzeError::Io(cause) => {
            assert_eq!(cause.kind(), io::ErrorKind::PermissionDenied);
            assert_eq!(cause.to_string(), "injected read");
        }
        other => panic!("expected I/O cause, got {other:?}"),
    }
}

#[test]
fn from_and_into_preserve_io_cause() {
    assert_io(AnalyzeError::from(io::Error::new(
        io::ErrorKind::PermissionDenied,
        "injected read",
    )));
    let converted: AnalyzeError =
        io::Error::new(io::ErrorKind::PermissionDenied, "injected read").into();
    assert_io(converted);
}

struct OneByte<R>(R);

impl<R: Read> Read for OneByte<R> {
    fn read(&mut self, buffer: &mut [u8]) -> io::Result<usize> {
        if buffer.is_empty() {
            return Ok(0);
        }
        self.0.read(&mut buffer[..1])
    }
}

struct FailAfterPrefix(Cursor<Vec<u8>>);

impl Read for FailAfterPrefix {
    fn read(&mut self, buffer: &mut [u8]) -> io::Result<usize> {
        if buffer.is_empty() {
            return Ok(0);
        }
        if self.0.position() < self.0.get_ref().len() as u64 {
            self.0.read(buffer)
        } else {
            Err(io::Error::new(
                io::ErrorKind::PermissionDenied,
                "injected read",
            ))
        }
    }
}

#[test]
fn cursor_file_short_reads_and_bytes_agree() {
    for name in ["sample", "empty", "malformed"] {
        let bytes = common::bytes(&format!("{name}.log"));
        let expected = analyze(&bytes).unwrap();
        assert_eq!(analyze_reader(Cursor::new(&bytes)).unwrap(), expected);
        assert_eq!(
            analyze_reader(OneByte(Cursor::new(&bytes))).unwrap(),
            expected
        );
        let file = std::fs::File::open(common::fixture(&format!("{name}.log"))).unwrap();
        assert_eq!(analyze_reader(file).unwrap(), expected);
    }
    for bytes in [b"".as_slice(), b"INFO api 1\r\n\xff\nERROR a 60000"] {
        assert_eq!(
            analyze_reader(OneByte(Cursor::new(bytes))).unwrap(),
            analyze(bytes).unwrap()
        );
    }
}

#[test]
fn read_failure_discards_even_a_valid_prefix() {
    for prefix in [b"".as_slice(), b"INFO api 1\n"] {
        let reader = FailAfterPrefix(Cursor::new(prefix.to_vec()));
        assert_io(analyze_reader(reader).unwrap_err());
    }
}
