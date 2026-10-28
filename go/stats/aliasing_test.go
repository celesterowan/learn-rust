package main

import (
	"reflect"
	"testing"
)

func TestAliasing(t *testing.T) {
	cases := []struct {
		name string
		lab  func(*testing.T)
	}{
		{"value_copy", func(t *testing.T) {
			original := Sample{Value: 4}
			copied := CopyScalar(original, 9)
			if original.Value != 4 || copied.Value != 9 {
				t.Fatalf("original=%+v copy=%+v", original, copied)
			}
		}},
		{"pointer_mutation", func(t *testing.T) {
			original := Sample{Value: 4}
			MutateScalar(&original, 9)
			if original.Value != 9 {
				t.Fatalf("original=%+v", original)
			}
		}},
		{"shared_slice", func(t *testing.T) {
			original := []uint32{4, 9}
			shared := ShareSlice(original)
			if !reflect.DeepEqual(shared, original) {
				t.Fatalf("view=%v original=%v", shared, original)
			}
			shared[0] = 7
			original[1] = 8
			if original[0] != 7 || shared[1] != 8 {
				t.Fatal("slice elements do not share storage")
			}
		}},
		{"copied_slice", func(t *testing.T) {
			original := []uint32{4, 9}
			copied := CopySlice(original)
			if !reflect.DeepEqual(copied, original) {
				t.Fatalf("copy=%v original=%v", copied, original)
			}
			copied[0] = 7
			original[1] = 8
			if original[0] != 4 || copied[1] != 9 {
				t.Fatal("slice copy still aliases")
			}
		}},
		{"shared_map", func(t *testing.T) {
			original := map[string]uint32{"a": 4}
			shared := ShareMap(original)
			if shared["a"] != 4 {
				t.Fatal("existing entry missing")
			}
			shared["a"] = 9
			original["b"] = 7
			if original["a"] != 9 || shared["b"] != 7 {
				t.Fatal("map entries do not share storage")
			}
		}},
		{"map_presence", func(t *testing.T) {
			values := map[string]uint32{"zero": 0, "positive": 9}
			for _, tc := range []struct {
				key     string
				value   uint32
				present bool
			}{
				{"missing", 0, false}, {"zero", 0, true}, {"positive", 9, true},
			} {
				value, present := Lookup(values, tc.key)
				if value != tc.value || present != tc.present {
					t.Fatalf("%q: got (%d, %t), want (%d, %t)",
						tc.key, value, present, tc.value, tc.present)
				}
			}
		}},
		{"append_capacity", func(t *testing.T) {
			original := make([]uint32, 1, 3)
			original[0] = 4
			window := original[:cap(original)]
			appended := AppendWithinCapacity(original, 9)
			if len(original) != 1 || len(appended) != 2 {
				t.Fatal("append must return a longer header without changing the original length")
			}
			if cap(appended) != cap(original) || appended[0] != 4 || appended[1] != 9 {
				t.Fatalf("unexpected append: %v (cap %d)", appended, cap(appended))
			}
			if window[1] != 9 {
				t.Fatal("append did not use existing spare capacity")
			}
			appended[0] = 7
			if original[0] != 7 {
				t.Fatal("append unexpectedly copied storage")
			}
		}},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			// Keep every lab visible even while another still has its starter panic.
			defer func() {
				if value := recover(); value != nil {
					t.Fatalf("lab panicked (starter TODOs must be implemented): %v", value)
				}
			}()
			tc.lab(t)
		})
	}
}
