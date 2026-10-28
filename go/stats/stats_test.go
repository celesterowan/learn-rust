package main

import (
	"errors"
	"reflect"
	"testing"
)

func TestCheckedAdd(t *testing.T) {
	max := ^uint64(0)
	cases := []struct {
		name        string
		accumulator uint64
		value       uint64
		want        uint64
		wantErr     error
	}{
		{"zero", 0, 0, 0, nil},
		{"max value", 0, max, max, nil},
		{"max accumulator", max, 0, max, nil},
		{"at limit", max - 1, 1, max, nil},
		{"over limit", max - 2, 3, 0, ErrOverflow},
		{"overflow", max, 1, 0, ErrOverflow},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := CheckedAdd(tc.accumulator, tc.value)
			if got != tc.want || !errors.Is(err, tc.wantErr) {
				t.Fatalf("got (%d, %v), want (%d, %v)", got, err, tc.want, tc.wantErr)
			}
		})
	}
}

func TestCalculate(t *testing.T) {
	max := ^uint32(0)
	cases := []struct {
		name   string
		values []uint32
		want   Stats
	}{
		{"nil", nil, Stats{}},
		{"empty", []uint32{}, Stats{}},
		{"singleton", []uint32{7}, Stats{1, 7, 7, 7, true, true}},
		{"zero", []uint32{0}, Stats{1, 0, 0, 0, true, true}},
		{"zeros", []uint32{0, 0, 0}, Stats{3, 0, 0, 0, true, true}},
		{"duplicates", []uint32{4, 4, 4}, Stats{3, 12, 4, 4, true, true}},
		{"ascending", []uint32{0, 1, 2, 3}, Stats{4, 6, 0, 3, true, true}},
		{"descending", []uint32{3, 2, 1, 0}, Stats{4, 6, 0, 3, true, true}},
		{"mixed", []uint32{4, 0, 9, 4}, Stats{4, 17, 0, 9, true, true}},
		{"max", []uint32{max}, Stats{1, 4294967295, max, max, true, true}},
		{"wide sum", []uint32{max, max, 0}, Stats{3, 8589934590, 0, max, true, true}},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			before := append([]uint32(nil), tc.values...)
			got, err := Calculate(tc.values)
			if err != nil || got != tc.want {
				t.Fatalf("got (%+v, %v), want (%+v, nil)", got, err, tc.want)
			}
			// Compare elements, without confusing a nil and an empty slice.
			if len(tc.values) != len(before) {
				t.Fatal("input length changed")
			}
			for i := range before {
				if tc.values[i] != before[i] {
					t.Fatal("input elements changed")
				}
			}
		})
	}
}

func TestCallerKeepsSlice(t *testing.T) {
	values := []uint32{4, 0, 9, 4}
	_, err := Calculate(values)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(values, []uint32{4, 0, 9, 4}) {
		t.Fatalf("caller input changed: %v", values)
	}
}
