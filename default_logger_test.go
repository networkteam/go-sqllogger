package sqllogger

import (
	"fmt"
	"testing"
)

type testLogger []string

func (tl *testLogger) Printf(format string, args ...interface{}) {
	*tl = append([]string(*tl), fmt.Sprintf(format, args...))
}

func TestNewDefaultSQLLogger(t *testing.T) {
	var l testLogger

	defaultSQLLogger := NewDefaultSQLLogger(&l)
	defaultSQLLogger.Connect(123)

	expectedEntries := []string{
		"Connect → CONN(123)",
	}

	if len(l) != len(expectedEntries) {
		t.Fatalf("expect %d log entries, but got %d", len(expectedEntries), len(l))
	}

	for i, entry := range l {
		if entry != expectedEntries[i] {
			t.Errorf("log entry at index %d expected to be %q, but got %q", i, expectedEntries[i], entry)
		}
	}
}
