package main

type Sample struct {
	Value uint32
}

// CopyScalar returns a changed value copy without changing the caller's scalar.
func CopyScalar(original Sample, replacement uint32) Sample {
	// TODO G6a: Copy by value, then change only the returned copy.
	panic("TODO G6a: CopyScalar")
}

// MutateScalar changes the original through a non-nil pointer.
func MutateScalar(original *Sample, replacement uint32) {
	// TODO G6b: Mutate through the pointer rather than a value copy.
	panic("TODO G6b: MutateScalar")
}

// ShareSlice returns a view with the same length and shared backing array.
func ShareSlice(values []uint32) []uint32 {
	// TODO G6c: Return a slice view without copying its elements.
	panic("TODO G6c: ShareSlice")
}

// CopySlice returns equal elements in independent storage.
func CopySlice(values []uint32) []uint32 {
	// TODO G6d: Allocate a separate backing array and copy the elements.
	panic("TODO G6d: CopySlice")
}

// ShareMap returns an alias of the same initialized map.
func ShareMap(values map[string]uint32) map[string]uint32 {
	// TODO G6e: Share entries rather than cloning the map.
	panic("TODO G6e: ShareMap")
}

// Lookup distinguishes a missing key from a present key whose value is zero.
func Lookup(values map[string]uint32, key string) (uint32, bool) {
	// TODO G6f: Use the two-result form of map lookup.
	panic("TODO G6f: Lookup")
}

// AppendWithinCapacity returns a slice one element longer, sharing storage.
// Precondition: len(values) < cap(values). Do not allocate a replacement array.
func AppendWithinCapacity(values []uint32, value uint32) []uint32 {
	// TODO G6g: Append and return the new header; the caller's length is unchanged.
	panic("TODO G6g: AppendWithinCapacity")
}
