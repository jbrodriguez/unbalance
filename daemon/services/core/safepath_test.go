package core

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"unbalance/daemon/domain"
)

func TestSafeJoinRejectsEscapes(t *testing.T) {
	root := filepath.Join(string(filepath.Separator), "mnt", "disk1")

	tests := []string{
		"",
		".",
		"..",
		"../disk2/share",
		filepath.Join("..", "disk2", "share"),
		filepath.Join(string(filepath.Separator), "mnt", "disk2", "share"),
	}

	for _, entry := range tests {
		if _, _, err := safeJoin(root, entry); err == nil {
			t.Fatalf("expected %q to be rejected", entry)
		}
	}
}

func TestSafeJoinAllowsRelativeEntries(t *testing.T) {
	root := filepath.Join(string(filepath.Separator), "mnt", "disk1")
	path, entry, err := safeJoin(root, "share/movie/file.mkv")
	if err != nil {
		t.Fatalf("safeJoin returned error: %s", err)
	}

	want := filepath.Join(root, "share", "movie", "file.mkv")
	if path != want {
		t.Fatalf("path = %q, want %q", path, want)
	}
	if entry != filepath.Join("share", "movie", "file.mkv") {
		t.Fatalf("entry = %q", entry)
	}
}

func TestRemoveTransferredSourceRemovesOnlyValidatedSource(t *testing.T) {
	tmp := t.TempDir()
	srcRoot := filepath.Join(tmp, "disk1")
	dstRoot := filepath.Join(tmp, "disk2")
	entry := filepath.Join("share", "movie")

	mustMkdirAll(t, filepath.Join(srcRoot, entry))
	mustWriteFile(t, filepath.Join(srcRoot, entry, "file.mkv"), "source")
	mustMkdirAll(t, filepath.Join(dstRoot, entry))
	mustWriteFile(t, filepath.Join(dstRoot, entry, "file.mkv"), "dest")

	removed, pruned, err := removeTransferredSource(&domain.Command{
		Src:   srcRoot,
		Dst:   dstRoot,
		Entry: entry,
	}, true)
	if err != nil {
		t.Fatalf("removeTransferredSource returned error: %s", err)
	}

	if removed != filepath.Join(srcRoot, entry) {
		t.Fatalf("removed = %q", removed)
	}

	if _, err := os.Stat(filepath.Join(srcRoot, entry)); !os.IsNotExist(err) {
		t.Fatalf("source still exists or unexpected stat error: %v", err)
	}

	if _, err := os.Stat(filepath.Join(dstRoot, entry)); err != nil {
		t.Fatalf("destination should remain: %s", err)
	}

	if len(pruned) != 0 {
		t.Fatalf("top-level share parent should not be pruned, got %v", pruned)
	}
}

func TestRemoveTransferredSourcePrunesNestedEmptyParents(t *testing.T) {
	tmp := t.TempDir()
	srcRoot := filepath.Join(tmp, "disk1")
	dstRoot := filepath.Join(tmp, "disk2")
	entry := filepath.Join("share", "show", "season", "episode.mkv")

	mustMkdirAll(t, filepath.Dir(filepath.Join(srcRoot, entry)))
	mustWriteFile(t, filepath.Join(srcRoot, entry), "source")
	mustMkdirAll(t, filepath.Dir(filepath.Join(dstRoot, entry)))
	mustWriteFile(t, filepath.Join(dstRoot, entry), "dest")

	_, pruned, err := removeTransferredSource(&domain.Command{
		Src:   srcRoot,
		Dst:   dstRoot,
		Entry: entry,
	}, true)
	if err != nil {
		t.Fatalf("removeTransferredSource returned error: %s", err)
	}

	if len(pruned) == 0 {
		t.Fatalf("expected nested empty parents to be pruned")
	}

	if _, err := os.Stat(filepath.Join(srcRoot, "share")); err != nil {
		t.Fatalf("share boundary should remain: %s", err)
	}
}

func TestRemoveTransferredSourceRefusesMissingDestination(t *testing.T) {
	tmp := t.TempDir()
	srcRoot := filepath.Join(tmp, "disk1")
	dstRoot := filepath.Join(tmp, "disk2")
	entry := filepath.Join("share", "movie")

	mustMkdirAll(t, filepath.Join(srcRoot, entry))
	mustMkdirAll(t, dstRoot)

	_, _, err := removeTransferredSource(&domain.Command{
		Src:   srcRoot,
		Dst:   dstRoot,
		Entry: entry,
	}, false)
	if err == nil {
		t.Fatalf("expected missing destination to be rejected")
	}

	if _, statErr := os.Stat(filepath.Join(srcRoot, entry)); statErr != nil {
		t.Fatalf("source should remain after refused delete: %s", statErr)
	}
}

func TestRemoveTransferredSourceRefusesSymlinkEscape(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("symlink behavior differs on windows")
	}

	tmp := t.TempDir()
	srcRoot := filepath.Join(tmp, "disk1")
	dstRoot := filepath.Join(tmp, "disk2")
	outside := filepath.Join(tmp, "outside")
	entry := filepath.Join("share", "escape")

	mustMkdirAll(t, filepath.Join(srcRoot, "share"))
	mustMkdirAll(t, filepath.Join(dstRoot, entry))
	mustMkdirAll(t, outside)
	mustWriteFile(t, filepath.Join(outside, "file.txt"), "outside")

	if err := os.Symlink(outside, filepath.Join(srcRoot, entry)); err != nil {
		t.Fatalf("unable to create symlink: %s", err)
	}

	_, _, err := removeTransferredSource(&domain.Command{
		Src:   srcRoot,
		Dst:   dstRoot,
		Entry: entry,
	}, false)
	if err == nil {
		t.Fatalf("expected symlink escape to be rejected")
	}

	if _, statErr := os.Stat(filepath.Join(outside, "file.txt")); statErr != nil {
		t.Fatalf("outside target should remain: %s", statErr)
	}
}

func mustMkdirAll(t *testing.T, path string) {
	t.Helper()
	if err := os.MkdirAll(path, 0o755); err != nil {
		t.Fatalf("mkdir %s: %s", path, err)
	}
}

func mustWriteFile(t *testing.T, path string, content string) {
	t.Helper()
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatalf("write %s: %s", path, err)
	}
}
