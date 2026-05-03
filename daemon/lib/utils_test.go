package lib

import (
	"os"
	"path/filepath"
	"testing"

	"unbalance/daemon/domain"
)

func TestSaveEnvWritesRestrictivePermissions(t *testing.T) {
	location := filepath.Join(t.TempDir(), "unbalanced.env")
	if err := os.WriteFile(location, []byte("DRY_RUN=false\n"), 0o644); err != nil {
		t.Fatalf("seed env: %s", err)
	}

	err := SaveEnv(location, domain.Config{
		DryRun:         true,
		NotifyPlan:     1,
		NotifyTransfer: 2,
		ReservedAmount: 42,
		ReservedUnit:   "GB",
		RsyncArgs:      []string{"-X"},
		Verbosity:      1,
		RefreshRate:    1000,
		AuthPassword:   "$argon2id$v=19$m=1,t=1,p=1$c2FsdA$hash",
	})
	if err != nil {
		t.Fatalf("SaveEnv returned error: %s", err)
	}

	info, err := os.Stat(location)
	if err != nil {
		t.Fatalf("stat env: %s", err)
	}

	if got := info.Mode().Perm(); got != 0o600 {
		t.Fatalf("env permissions = %o, want 600", got)
	}
}
