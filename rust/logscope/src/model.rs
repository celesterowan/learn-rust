#[derive(Debug, Clone, Copy, PartialEq, Eq)]
pub enum Level {
    Info,
    Warn,
    Error,
}

#[derive(Debug, PartialEq, Eq)]
pub struct Record<'a> {
    pub level: Level,
    pub service: &'a str,
    pub duration: u32,
}

#[derive(Debug, PartialEq, Eq)]
pub struct OwnedRecord {
    pub level: Level,
    pub service: String,
    pub duration: u32,
}

/// Borrow already-typed data without validation or allocation (week 3).
pub fn record(level: Level, service: &str, duration: u32) -> Record<'_> {
    // TODO W3-1: Construct a borrowed record; do not copy service storage.
    let _ = (level, service, duration);
    todo!("W3-1: record")
}

/// The returned borrow lasts as long as service storage, not the wrapper.
pub fn service_name<'a>(record: &Record<'a>) -> &'a str {
    // TODO W3-2: Return the service borrow with its original lifetime.
    let _ = record;
    todo!("W3-2: service_name")
}

impl Record<'_> {
    /// Consume the wrapper and allocate independently owned service storage.
    pub fn into_owned(self) -> OwnedRecord {
        // TODO W3-3: Preserve fields while making the service an owned String.
        todo!("W3-3: Record::into_owned")
    }
}
