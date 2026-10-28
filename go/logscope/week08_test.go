package logscope

import (
	"bytes"
	"errors"
	"math"
	"reflect"
	"strings"
	"testing"
)

func TestWeek08ParseLine(t *testing.T) {
	valid := []struct {
		line string
		want Record
	}{
		{"INFO a 0", Record{LevelInfo, "a", 0}},
		{"WARN abcdefghijklmnop 60000", Record{LevelWarn, "abcdefghijklmnop", 60000}},
		{"ERROR a0-z 00001", Record{LevelError, "a0-z", 1}},
		{"INFO a 00000", Record{LevelInfo, "a", 0}},
		{"INFO a 9", Record{LevelInfo, "a", 9}},
	}
	for _, test := range valid {
		t.Run(test.line, func(t *testing.T) {
			defer catchTODO(t)
			input := []byte(test.line)
			before := append([]byte(nil), input...)
			got, err := ParseLine(input)
			if err != nil || got != test.want {
				t.Fatalf("ParseLine = %#v, %v; want %#v", got, err, test.want)
			}
			if !bytes.Equal(input, before) {
				t.Fatal("parser mutated caller bytes")
			}
			for index := range input {
				input[index] = 'x'
			}
			if got.Service != test.want.Service {
				t.Fatal("record service aliases mutable input")
			}
		})
	}
	invalid := []struct {
		name string
		line []byte
	}{
		{"empty", nil},
		{"17_byte_service", []byte("INFO abcdefghijklmnopq 1")},
		{"empty_service", []byte("INFO  1")},
		{"uppercase_service", []byte("INFO Api 1")},
		{"digit_first", []byte("INFO 1api 1")},
		{"hyphen_first", []byte("INFO -api 1")},
		{"underscore", []byte("INFO api_x 1")},
		{"mixed_case_level", []byte("Info api 1")},
		{"unknown_level", []byte("DEBUG api 1")},
		{"60001", []byte("INFO api 60001")},
		{"six_digits", []byte("INFO api 000001")},
		{"negative", []byte("INFO api -1")},
		{"positive_sign", []byte("INFO api +1")},
		{"decimal", []byte("INFO api 1.0")},
		{"unit", []byte("INFO api 1ms")},
		{"missing_field", []byte("INFO api")},
		{"extra_field", []byte("INFO api 1 extra")},
		{"empty_duration", []byte("INFO api ")},
		{"tabs", []byte("INFO\tapi\t1")},
		{"double_space", []byte("INFO  api 1")},
		{"leading_space", []byte(" INFO api 1")},
		{"trailing_space", []byte("INFO api 1 ")},
		{"whitespace_only", []byte(" \t ")},
		{"non_ascii", []byte("INFO caf\xc3\xa9 1")},
		{"unicode_digit", []byte("INFO api \xef\xbc\x91")},
		{"invalid_utf8", []byte{'I', 'N', 'F', 'O', ' ', 0xff, ' ', '1'}},
		{"nul", []byte("INFO api 1\x00")},
		{"cr", []byte("INFO api 1\r")},
		{"lf", []byte("INFO api 1\n")},
	}
	for _, test := range invalid {
		t.Run(test.name, func(t *testing.T) {
			defer catchTODO(t)
			got, err := ParseLine(test.line)
			var parseError *ParseError
			if got != (Record{}) || !errors.Is(err, ErrInvalidRecord) || !errors.As(err, &parseError) {
				t.Fatalf("malformed line = %#v, %v; want zero record and *ParseError wrapping ErrInvalidRecord", got, err)
			}
			if parseError.Reason == "" {
				t.Error("typed parse error needs an informative reason")
			}
		})
	}
}

func TestWeek08CheckedAdd(t *testing.T) {
	for _, test := range []struct {
		left, right, want uint64
		overflow          bool
	}{
		{0, 0, 0, false},
		{40, 2, 42, false},
		{math.MaxUint64, 0, math.MaxUint64, false},
		{math.MaxUint64 - 1, 1, math.MaxUint64, false},
		{math.MaxUint64, 1, 0, true},
		{1, math.MaxUint64, 0, true},
	} {
		got, err := checkedAdd(test.left, test.right)
		if got != test.want || (test.overflow && !errors.Is(err, ErrOverflow)) || (!test.overflow && err != nil) {
			t.Errorf("checkedAdd(%d,%d) = %d, %v", test.left, test.right, got, err)
		}
	}
}

func TestWeek08NormalizedLines(t *testing.T) {
	defer catchTODO(t)
	input := []byte("\n\r\nINFO a 1\r\n\nWARN b 2\r\r\n\xff\n\r")
	before := append([]byte(nil), input...)
	got := normalizedLines(input)
	want := [][]byte{[]byte("INFO a 1"), []byte("WARN b 2\r"), {0xff}}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("normalized lines = %q; want %q", got, want)
	}
	if !bytes.Equal(input, before) {
		t.Error("normalization mutated input")
	}
	if len(normalizedLines(nil)) != 0 || len(normalizedLines([]byte("\r\n\n\r"))) != 0 {
		t.Error("empty normalized records must be omitted")
	}
}

func TestWeek08AddRecord(t *testing.T) {
	defer catchTODO(t)
	var got Summary
	for _, record := range []Record{{LevelInfo, "api", 0}, {LevelWarn, "api", 7}, {LevelError, "worker", 3}} {
		if err := addRecord(&got, record); err != nil {
			t.Fatal(err)
		}
	}
	want := Summary{
		Valid:    3,
		Levels:   map[Level]uint64{LevelInfo: 1, LevelWarn: 1, LevelError: 1},
		Services: map[string]ServiceSummary{"api": {2, 7, 7}, "worker": {1, 3, 3}},
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("aggregate = %#v; want %#v", got, want)
	}
}

func TestWeek08Overflow(t *testing.T) {
	for _, field := range []string{"valid", "level", "service_count", "service_total", "invalid"} {
		t.Run(field, func(t *testing.T) {
			defer catchTODO(t)
			summary := Summary{Levels: map[Level]uint64{}, Services: map[string]ServiceSummary{}}
			switch field {
			case "valid":
				summary.Valid = math.MaxUint64
			case "level":
				summary.Levels[LevelInfo] = math.MaxUint64
			case "service_count":
				summary.Services["api"] = ServiceSummary{Count: math.MaxUint64}
			case "service_total":
				summary.Services["api"] = ServiceSummary{TotalMS: math.MaxUint64}
			case "invalid":
				summary.Invalid = math.MaxUint64
			}
			line := []byte("INFO api 1")
			if field == "invalid" {
				line = []byte("bad")
			}
			if err := accumulateLine(&summary, line); !errors.Is(err, ErrOverflow) {
				t.Errorf("overflow in %s = %v; want ErrOverflow", field, err)
			}
		})
	}
}

func TestWeek08AccumulateLine(t *testing.T) {
	defer catchTODO(t)
	var summary Summary
	if err := accumulateLine(&summary, []byte("bad")); err != nil {
		t.Fatal(err)
	}
	if summary.Invalid != 1 || summary.Valid != 0 || len(summary.Services) != 0 {
		t.Fatalf("malformed record contaminated aggregate: %#v", summary)
	}
	if err := accumulateLine(&summary, []byte("ERROR api 5")); err != nil {
		t.Fatal(err)
	}
	if summary.Valid != 1 || summary.Invalid != 1 || summary.Levels[LevelError] != 1 || summary.Services["api"] != (ServiceSummary{1, 5, 5}) {
		t.Errorf("accumulation = %#v", summary)
	}
}

func TestWeek08Analyze(t *testing.T) {
	for _, test := range []struct {
		name, input    string
		valid, invalid uint64
	}{
		{"empty", "", 0, 0},
		{"blank_lines", "\n\r\n\r", 0, 0},
		{"crlf", "INFO api 1\r\nWARN api 2\r\n", 2, 0},
		{"unterminated", "ERROR api 3", 1, 0},
		{"unterminated_cr", "INFO api 1\r", 1, 0},
		{"one_cr_only", "INFO api 1\r\r\n\r\r", 0, 2},
		{"invalid_bytes", "\xff\n\x00\nINFO api 1\n", 1, 2},
		{"long_malformed", strings.Repeat("x", 128*1024) + "\nINFO api 1\n", 1, 1},
	} {
		t.Run(test.name, func(t *testing.T) {
			defer catchTODO(t)
			input := []byte(test.input)
			before := append([]byte(nil), input...)
			got, err := Analyze(input)
			if err != nil {
				t.Fatal(err)
			}
			if got.Valid != test.valid || got.Invalid != test.invalid {
				t.Errorf("Analyze = %#v; want valid=%d invalid=%d", got, test.valid, test.invalid)
			}
			if !bytes.Equal(input, before) {
				t.Error("Analyze mutated caller bytes")
			}
		})
	}
}

func TestWeek08Render(t *testing.T) {
	for _, test := range []struct {
		name    string
		summary Summary
		want    string
	}{
		{"empty", Summary{}, "valid=0 invalid=0\nlevel INFO=0 WARN=0 ERROR=0\n"},
		{"all_levels_and_sorting", Summary{
			Valid: 3, Invalid: 2,
			Levels:   map[Level]uint64{LevelError: 1, LevelInfo: 2},
			Services: map[string]ServiceSummary{"z": {1, 0, 0}, "a-": {1, 8, 8}, "a": {1, 9, 9}},
		}, "valid=3 invalid=2\nlevel INFO=2 WARN=0 ERROR=1\nservice a count=1 total_ms=9 max_ms=9\nservice a- count=1 total_ms=8 max_ms=8\nservice z count=1 total_ms=0 max_ms=0\n"},
		{"full_range", Summary{Invalid: math.MaxUint64}, "valid=0 invalid=18446744073709551615\nlevel INFO=0 WARN=0 ERROR=0\n"},
	} {
		t.Run(test.name, func(t *testing.T) {
			defer catchTODO(t)
			for repeat := 0; repeat < 5; repeat++ {
				if got := Render(test.summary); got != test.want {
					t.Errorf("Render bytes = %q; want %q", got, test.want)
				}
			}
		})
	}
}

func TestWeek08Fixtures(t *testing.T) {
	for _, name := range []string{"sample", "empty", "malformed"} {
		t.Run(name, func(t *testing.T) {
			defer catchTODO(t)
			summary, err := Analyze(fixture(t, name+".log"))
			if err != nil {
				t.Fatal(err)
			}
			if got, want := Render(summary), string(fixture(t, name+".expected.txt")); got != want {
				t.Errorf("report bytes = %q; want %q", got, want)
			}
		})
	}
}
