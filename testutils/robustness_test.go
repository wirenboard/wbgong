package testutils

import (
	"os"
	"testing"
	"time"
)

// A full record buffer must not block the producer (previously the driver
// goroutine deadlocked mid-transaction when a test generated more records
// than it Verified). Overflow drops with a log marker instead.
func TestRecorderOverflowDoesNotBlock(t *testing.T) {
	rec := NewRecorder(t)
	done := make(chan struct{})
	go func() {
		defer close(done)
		for i := 0; i < 1500; i++ {
			rec.Rec("item %d", i)
		}
	}()
	select {
	case <-done:
	case <-time.After(10 * time.Second):
		t.Fatal("Rec blocked on a full buffer")
	}
	// the buffered 1000 records are still there for Verify
	for i := 0; i < 1000; i++ {
		rec.Verify("item " + itoa(i))
	}
	rec.VerifyEmpty()
}

func itoa(n int) string {
	if n == 0 {
		return "0"
	}
	var b [8]byte
	i := len(b)
	for n > 0 {
		i--
		b[i] = byte('0' + n%10)
		n /= 10
	}
	return string(b[i:])
}

// SetupTempDir chdirs into the temp dir; if a test fails hard (FailNow or
// panic) its explicit cleanup never runs, and before this fix the whole
// process stayed chdir'd into a removed directory, silently breaking
// cwd-relative logic in every later test. t.Cleanup must restore it.
func TestSetupTempDirRestoresCwdWithoutExplicitCleanup(t *testing.T) {
	wd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	t.Run("leaky", func(t *testing.T) {
		dir, _ := SetupTempDir(t)
		if dir == "" {
			t.Fatal("no temp dir")
		}
		// deliberately no cleanup call - simulates FailNow/panic paths
	})
	after, err := os.Getwd()
	if err != nil {
		t.Fatalf("cwd stranded in a removed directory: %v", err)
	}
	if after != wd {
		t.Fatalf("cwd not restored: %q != %q", after, wd)
	}
}
