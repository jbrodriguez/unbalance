package core

import "testing"

func TestValidateRsyncArgsAllowsCommonLocalOptions(t *testing.T) {
	args := []string{
		"-avPR",
		"-X",
		"-A",
		"-H",
		"--numeric-ids",
		"--inplace",
		"--info=progress2",
		"--partial",
		"--bwlimit=50000",
		"--dry-run",
		"",
		"  ",
	}

	if err := validateRsyncArgs(args); err != nil {
		t.Fatalf("expected args to be allowed: %s", err)
	}
}

func TestValidateRsyncArgsBlocksDangerousOptions(t *testing.T) {
	tests := [][]string{
		{"--delete"},
		{"--delete-before"},
		{"--delete-after"},
		{"--delete-excluded"},
		{"--remove-source-files"},
		{"--rsync-path=/tmp/rsync"},
		{"--rsync-path", "/tmp/rsync"},
		{"-e", "ssh"},
		{"-essh"},
		{"--rsh=ssh"},
		{"--rsh", "ssh"},
	}

	for _, args := range tests {
		if err := validateRsyncArgs(args); err == nil {
			t.Fatalf("expected %v to be blocked", args)
		}
	}
}

func TestCleanRsyncArgsTrimsAndDropsEmptyArgs(t *testing.T) {
	got := cleanRsyncArgs([]string{" -X ", "", "\t", "--partial"})
	want := []string{"-X", "--partial"}

	if len(got) != len(want) {
		t.Fatalf("len = %d, want %d (%v)", len(got), len(want), got)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("got[%d] = %q, want %q", i, got[i], want[i])
		}
	}
}
