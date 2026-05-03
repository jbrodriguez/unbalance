package core

import (
	"os"
	"path/filepath"
	"testing"

	"unbalance/daemon/domain"
)

// TestHelperSecurityProbe is a manual probe harness for the no-shell migration
// of getIssues/getItems. It is skipped unless UNBALANCE_PROBE_DISK and
// UNBALANCE_PROBE_PATH are set, so it never runs in normal test runs.
//
//	UNBALANCE_PROBE_DISK=/mnt/disk1/testing/security-helper-fixtures \
//	UNBALANCE_PROBE_PATH=root \
//	./unbalance-core.test -test.v -test.run TestHelperSecurityProbe
//
// The probe calls the same unexported helpers the planner uses, so we exercise
// the real find/stat and find/du paths against the fixture without standing up
// the full daemon.
func TestHelperSecurityProbe(t *testing.T) {
	diskPath := os.Getenv("UNBALANCE_PROBE_DISK")
	subPath := os.Getenv("UNBALANCE_PROBE_PATH")
	if diskPath == "" || subPath == "" {
		t.Skip("set UNBALANCE_PROBE_DISK and UNBALANCE_PROBE_PATH to run the security probe")
	}

	full := filepath.Join(diskPath, subPath)
	if _, err := os.Stat(full); err != nil {
		t.Fatalf("fixture %s not accessible: %v", full, err)
	}

	disk := &domain.Disk{Path: diskPath}

	t.Logf("probe disk=%s sub=%s full=%s", diskPath, subPath, full)

	ownerIssue, groupIssue, folderIssue, fileIssue, err := getIssues(reStat, disk, subPath)
	if err != nil {
		t.Fatalf("getIssues error: %v", err)
	}
	t.Logf("issues: owner=%d group=%d folder=%d file=%d", ownerIssue, groupIssue, folderIssue, fileIssue)

	items, total, err := getItems(0, reItems, diskPath, subPath)
	if err != nil {
		t.Fatalf("getItems error: %v", err)
	}
	t.Logf("items: count=%d total_bytes=%d", len(items), total)
	for _, it := range items {
		t.Logf("  item size=%d name=%s path=%s", it.Size, it.Name, it.Path)
	}
}
