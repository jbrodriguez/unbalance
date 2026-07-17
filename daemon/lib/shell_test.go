package lib

import (
	"context"
	"testing"
	"time"
)

func TestStreamContextKillsChild(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		time.Sleep(100 * time.Millisecond)
		cancel()
	}()

	start := time.Now()
	err := StreamContext(ctx, "sleep", []string{"5"}, func(string) {})
	elapsed := time.Since(start)

	if err == nil {
		t.Fatalf("expected an error from a killed child")
	}

	if elapsed > 3*time.Second {
		t.Fatalf("cancellation took too long: %s", elapsed)
	}
}

func TestStreamHandlesLongLines(t *testing.T) {
	// 200KB on one line: beyond bufio's 64KB default, within our limit
	err := Shell("head -c 200000 /dev/zero | tr '\\0' a", "", func(line string) {
		if len(line) != 200000 {
			t.Fatalf("expected a 200000-char line, got %d", len(line))
		}
	})
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}
}

func TestStreamDoesNotDeadlockOnOversizedLine(t *testing.T) {
	// a line beyond maxLineSize aborts the scan; the child keeps writing and
	// must still be reaped instead of deadlocking Wait
	done := make(chan error, 1)
	go func() {
		done <- Shell("head -c 4000000 /dev/zero | tr '\\0' a; echo tail", "", func(string) {})
	}()

	select {
	case err := <-done:
		if err == nil {
			t.Fatalf("expected an error for an oversized line")
		}
	case <-time.After(10 * time.Second):
		t.Fatalf("Shell deadlocked on an oversized line")
	}
}
