// Week 7 learner scaffold. See docs/exercises/stats.md at the course root.
package main

import "errors"

var (
	ErrOverflow        = errors.New("arithmetic overflow")
	ErrInvalidArgument = errors.New("invalid numeric argument")
	ErrUnimplemented   = errors.New("TODO: implement the requested exercise step")
)

// Stats uses HasMin/HasMax to distinguish absent extrema from a present zero.
// When absent, the associated numeric field must be zero.
type Stats struct {
	Count  uint64
	Sum    uint64
	Min    uint32
	Max    uint32
	HasMin bool
	HasMax bool
}

// CheckedAdd returns ErrOverflow rather than wrapping a uint64 sum.
func CheckedAdd(accumulator, value uint64) (uint64, error) {
	// TODO G1: Check the remaining range before adding; error results return zero.
	return 0, ErrUnimplemented
}

// Calculate accepts the full uint32 range and leaves values unchanged.
// Empty input has zero count/sum and absent extrema. Check count and sum.
// On overflow return Stats{} and ErrOverflow, not a partial aggregate.
func Calculate(values []uint32) (Stats, error) {
	// TODO G2: Accumulate with CheckedAdd and set extrema presence explicitly.
	return Stats{}, ErrUnimplemented
}
