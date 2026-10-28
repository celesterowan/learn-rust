use logscope::{parse_line, Level, ParseError};

#[test]
fn starter_parser() {
    let record = parse_line(b"INFO api 10").unwrap();
    assert_eq!(record.level, Level::Info);
    assert_eq!(record.service, "api");
    assert_eq!(record.duration, 10);
    for line in [b"INFO api".as_slice(), b"INFO api -1"] {
        assert_eq!(parse_line(line), Err(ParseError::Malformed));
    }
}

#[test]
fn valid_grammar_and_borrowed_service() {
    for (line, level, service, duration) in [
        ("INFO a 0", Level::Info, "a", 0),
        (
            "WARN abcdefghijklmnop 60000",
            Level::Warn,
            "abcdefghijklmnop",
            60000,
        ),
        ("ERROR a0-9 00001", Level::Error, "a0-9", 1),
        ("INFO api 00000", Level::Info, "api", 0),
        ("INFO api 59999", Level::Info, "api", 59999),
    ] {
        let value = parse_line(line.as_bytes()).unwrap();
        assert_eq!(value.level, level);
        assert_eq!(value.service, service);
        assert_eq!(value.duration, duration);
        let offset = line.find(' ').unwrap() + 1;
        assert_eq!(value.service.as_ptr(), line.as_bytes()[offset..].as_ptr());
    }
}

#[test]
fn malformed_byte_matrix_never_panics() {
    let cases: &[&[u8]] = &[
        b"",
        b" ",
        b"\t",
        b"INFO",
        b"INFO api",
        b"INFO api 1 extra",
        b"INFO  api 1",
        b"INFO api  1",
        b" INFO api 1",
        b"INFO api 1 ",
        b"INFO\tapi 1",
        b"INFO api\t1",
        b"INFO api 1\n",
        b"INFO api 1\r",
        b"INFO api 1\r\r",
        b"info api 1",
        b"Info api 1",
        b"DEBUG api 1",
        b"WARN Api 1",
        b"ERROR 1api 1",
        b"INFO -api 1",
        b"INFO a_b 1",
        b"INFO abcdefghijklmnopq 1",
        b"INFO  1",
        b"INFO api ",
        b"INFO api 60001",
        b"INFO api 100000",
        b"INFO api 000001",
        b"INFO api -1",
        b"INFO api +1",
        b"INFO api 1.0",
        b"INFO api 1ms",
        b"INFO api 999999999999999999999999999999999999",
        b"INFO ap\0i 1",
        b"INFO api 1\0",
        b"\0",
        b"INFO ap\xffi 1",
        b"\xffINFO api 1",
        b"INFO api \xff",
        b"INFO api \xc3\xa9",
        b"INFO \xc3\xa9 1",
        b"INFO api \xd9\xa1",
        b"INFO api \xef\xbc\x91",
        b"INFO\xc2\xa0api 1",
        b"INFO api 1\x0b",
        b"INFO api 1\x0c",
    ];
    for line in cases {
        assert_eq!(parse_line(line), Err(ParseError::Malformed), "{line:?}");
    }
}
