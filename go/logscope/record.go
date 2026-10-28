// Learner scaffold, not a solution. See docs/exercises/go-logscope.md.
package logscope

import "errors"

var (
	ErrInvalidRecord = errors.New("invalid log record")
	ErrOverflow      = errors.New("analysis arithmetic overflow")
	ErrWorkers       = errors.New("workers must be between 1 and 4")
	ErrUnimplemented = errors.New("TODO: implement the requested exercise step")
)

type Level string

const (
	LevelInfo  Level = "INFO"
	LevelWarn  Level = "WARN"
	LevelError Level = "ERROR"
)

type Record struct {
	Level      Level
	Service    string
	DurationMS uint64
}

// ParseError describes a malformed record; errors.Is recognizes ErrInvalidRecord.
type ParseError struct {
	Reason string
}

func (err *ParseError) Error() string {
	return "invalid log record: " + err.Reason
}

func (err *ParseError) Unwrap() error {
	return ErrInvalidRecord
}

// ParseLine accepts one normalized nonempty line, with no trailing CR or LF.
// The returned service string owns its bytes independently of line.
func ParseLine(line []byte) (Record, error) {
	// TODO G8.1: Validate the exact ASCII grammar; return Record{} and *ParseError
	// for malformed input. Do not trim or split on Unicode whitespace.
	return Record{}, ErrUnimplemented
}
