use std::process::Command;
use week01_hello::{classify, greeting, HelloError};

#[test]
fn starter_classification() {
    assert_eq!(classify(0), "zero");
    for value in [1, 7, u32::MAX] {
        assert_eq!(classify(value), "nonzero");
    }
}

#[test]
fn greeting_contract() {
    for (arguments, expected) in [
        (vec![], "Hello, learner!\n"),
        (vec!["Ada"], "Hello, Ada!\n"),
        (vec!["Ada Lovelace"], "Hello, Ada Lovelace!\n"),
        (vec![" Ada "], "Hello,  Ada !\n"),
        (vec![" "], "Hello,  !\n"),
        (vec!["\u{4e16}\u{754c}"], "Hello, \u{4e16}\u{754c}!\n"),
    ] {
        let owned: Vec<String> = arguments.iter().map(|s| s.to_string()).collect();
        assert_eq!(greeting(&owned), Ok(expected.to_string()));
        let output = Command::new(env!("CARGO_BIN_EXE_week01-hello"))
            .args(&arguments)
            .output()
            .unwrap();
        assert_eq!(output.status.code(), Some(0));
        assert_eq!(output.stdout, expected.as_bytes());
        assert!(output.stderr.is_empty());
    }
}

#[test]
fn invalid_arguments_are_not_panics() {
    for arguments in [vec![""], vec!["Ada", "extra"], vec!["", "extra"]] {
        let owned: Vec<String> = arguments.iter().map(|s| s.to_string()).collect();
        assert_eq!(greeting(&owned), Err(HelloError::InvalidArguments));
        let output = Command::new(env!("CARGO_BIN_EXE_week01-hello"))
            .args(&arguments)
            .output()
            .unwrap();
        assert_eq!(output.status.code(), Some(2));
        assert!(output.stdout.is_empty());
        assert!(output.stderr.starts_with(b"error:"));
    }
}

#[test]
fn starter_mutability_and_shadowing() {
    // TODO H3: Write a passing mutation/shadowing experiment with assertions.
    // Record the deliberate compiler diagnostic in your reflection, not here.
    todo!("H3: mutability and shadowing experiment")
}
