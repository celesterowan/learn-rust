package logscope

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

// Keep a TODO panic visible as a failed test rather than aborting later gates.
func catchTODO(t *testing.T) {
	t.Helper()
	if value := recover(); value != nil {
		t.Errorf("unexpected panic (learner TODOs must be implemented): %v", value)
	}
}

func fixture(t *testing.T, name string) []byte {
	t.Helper()
	data, err := os.ReadFile(filepath.Join("..", "..", "fixtures", name))
	if err != nil {
		t.Fatal(err)
	}
	return data
}

func requireZeroSummary(t *testing.T, summary Summary) {
	t.Helper()
	if !reflect.DeepEqual(summary, Summary{}) {
		t.Errorf("error returned a partial summary: %#v", summary)
	}
}
