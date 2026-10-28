use std::path::PathBuf;

pub fn fixture(name: &str) -> PathBuf {
    PathBuf::from(env!("CARGO_MANIFEST_DIR"))
        .join("..")
        .join("..")
        .join("fixtures")
        .join(name)
}

pub fn bytes(name: &str) -> Vec<u8> {
    std::fs::read(fixture(name)).unwrap()
}
