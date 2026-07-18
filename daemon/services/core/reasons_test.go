package core

import (
	"testing"
)

func TestRsyncExitReasonKnownCode(t *testing.T) {
	reason := rsyncExitReason(23)

	if reason != "Partial transfer due to error (rsync exit code 23)" {
		t.Fatalf("unexpected reason: %s", reason)
	}
}

func TestRsyncExitReasonUnknownCode(t *testing.T) {
	reason := rsyncExitReason(42)

	if reason != "unknown error (rsync exit code 42)" {
		t.Fatalf("unexpected reason: %s", reason)
	}
}

func TestGetErrorParsesExitStatus(t *testing.T) {
	msg := getError("exit status 11", reRsync, rsyncErrors)

	if msg != "Error in file I/O" {
		t.Fatalf("unexpected message: %s", msg)
	}
}

func TestGetErrorUnknownLine(t *testing.T) {
	msg := getError("unsafe command path: boom", reRsync, rsyncErrors)

	if msg != "unknown error" {
		t.Fatalf("unexpected message: %s", msg)
	}
}
