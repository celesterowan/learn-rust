use logscope::{record, service_name, Level, Record};

#[test]
fn typed_records_borrow_without_validation() {
    let service = String::from("api");
    for (level, duration) in [
        (Level::Info, 10),
        (Level::Warn, 20),
        (Level::Error, 100),
        (Level::Info, 0),
        (Level::Error, u32::MAX),
    ] {
        let value = record(level, &service, duration);
        assert_eq!(value.level, level);
        assert_eq!(value.duration, duration);
        assert_eq!(value.service, "api");
        assert_eq!(value.service.as_ptr(), service.as_ptr());
    }
    // Week 3 accepts typed data, including names outside week 4's grammar.
    assert_eq!(record(Level::Info, "", 60001).service, "");
}

#[test]
fn service_borrow_outlives_wrapper() {
    let storage = String::from("worker");
    let name = {
        let wrapper = Record {
            level: Level::Error,
            service: &storage,
            duration: 100,
        };
        service_name(&wrapper)
    };
    assert_eq!(name, "worker");
    assert_eq!(name.as_ptr(), storage.as_ptr());
}

#[test]
fn owned_conversion_outlives_storage() {
    for level in [Level::Info, Level::Warn, Level::Error] {
        let owned = {
            let storage = String::from("worker");
            record(level, &storage, 100).into_owned()
        };
        assert_eq!(owned.level, level);
        assert_eq!(owned.service, "worker");
        assert_eq!(owned.duration, 100);
    }
}

struct DropMarker<'a>(&'a mut bool);

impl Drop for DropMarker<'_> {
    fn drop(&mut self) {
        // TODO W3-S2: Set the borrowed flag on destruction.
        let _ = &self.0;
        todo!("W3-S2: DropMarker::drop")
    }
}

#[test]
fn starter_drop_marker() {
    let mut dropped = false;
    let marker = DropMarker(&mut dropped);
    drop(marker);
    assert!(dropped);
}

#[test]
fn starter_moves_copies_clones_and_borrows() {
    // TODO W3-S1: Add assertions demonstrating a String move, u32 copy, String
    // clone, and shared then mutable borrow after the shared borrow's last use.
    // Record a deliberate use-after-move diagnostic outside this passing test.
    todo!("W3-S1: ownership experiments")
}
