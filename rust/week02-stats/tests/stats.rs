use week02_stats::{calculate, checked_add, classify, Dataset, Stats, StatsError};

#[test]
fn starter_classifies_slices() {
    let array = [4, 9];
    let mut vector = vec![4];
    vector.push(9);
    let slice = &vector[..];
    match classify(&[]) {
        Dataset::Empty => {}
        Dataset::Nonempty(count) => panic!("empty slice classified as {count} elements"),
    }
    for values in [&array[..], slice] {
        match classify(values) {
            Dataset::Empty => panic!("nonempty slice classified as empty"),
            Dataset::Nonempty(count) => assert_eq!(count, 2),
        }
    }
}

#[test]
fn checked_add_boundaries() {
    for (accumulator, value, expected) in [
        (0, 0, Ok(0)),
        (0, u64::MAX, Ok(u64::MAX)),
        (u64::MAX, 0, Ok(u64::MAX)),
        (u64::MAX - 1, 1, Ok(u64::MAX)),
        (u64::MAX - 2, 3, Err(StatsError::Overflow)),
        (u64::MAX, 1, Err(StatsError::Overflow)),
    ] {
        assert_eq!(checked_add(accumulator, value), expected);
    }
}

#[test]
fn statistics_cases() {
    type Case = (&'static str, &'static [u32], u64, Option<u32>, Option<u32>);
    let cases: &[Case] = &[
        ("empty", &[], 0, None, None),
        ("singleton", &[7], 7, Some(7), Some(7)),
        ("zero", &[0], 0, Some(0), Some(0)),
        ("zeros", &[0, 0, 0], 0, Some(0), Some(0)),
        ("duplicates", &[4, 4, 4], 12, Some(4), Some(4)),
        ("ascending", &[0, 1, 2, 3], 6, Some(0), Some(3)),
        ("descending", &[3, 2, 1, 0], 6, Some(0), Some(3)),
        ("mixed", &[4, 0, 9, 4], 17, Some(0), Some(9)),
        (
            "max",
            &[u32::MAX],
            4_294_967_295,
            Some(u32::MAX),
            Some(u32::MAX),
        ),
        (
            "wide sum",
            &[u32::MAX, u32::MAX, 0],
            8_589_934_590,
            Some(0),
            Some(u32::MAX),
        ),
    ];
    for &(name, values, sum, min, max) in cases {
        assert_eq!(
            calculate(values),
            Ok(Stats {
                count: values.len() as u64,
                sum,
                min,
                max,
            }),
            "{name}"
        );
    }
}

#[test]
fn caller_keeps_vector() {
    let values = vec![4, 0, 9, 4];
    assert_eq!(calculate(&values).unwrap().sum, 17);
    assert_eq!(values, [4, 0, 9, 4]);
    assert_eq!(values.iter().sum::<u32>(), 17);
}
