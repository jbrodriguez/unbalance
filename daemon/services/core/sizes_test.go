package core

import (
	"os"
	"path/filepath"
	"testing"

	"unbalance/daemon/domain"
)

func writeSizedFile(t *testing.T, path string, size int) {
	t.Helper()

	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatalf("unable to create dir %s: %s", filepath.Dir(path), err)
	}

	if err := os.WriteFile(path, make([]byte, size), 0o644); err != nil {
		t.Fatalf("unable to write %s: %s", path, err)
	}
}

func TestEntrySizeFile(t *testing.T) {
	dir := t.TempDir()
	file := filepath.Join(dir, "movie.mkv")
	writeSizedFile(t, file, 1500)

	fi, err := os.Lstat(file)
	if err != nil {
		t.Fatal(err)
	}

	size, err := entrySize(file, fi)
	if err != nil {
		t.Fatal(err)
	}

	if size != 1500 {
		t.Fatalf("expected 1500, got %d", size)
	}
}

func TestEntrySizeNestedFolder(t *testing.T) {
	dir := t.TempDir()
	writeSizedFile(t, filepath.Join(dir, "films", "alien", "alien.mkv"), 1000)
	writeSizedFile(t, filepath.Join(dir, "films", "alien", "poster.jpg"), 200)
	writeSizedFile(t, filepath.Join(dir, "films", "aliens", "aliens.mkv"), 3000)

	folder := filepath.Join(dir, "films")
	fi, err := os.Lstat(folder)
	if err != nil {
		t.Fatal(err)
	}

	size, err := entrySize(folder, fi)
	if err != nil {
		t.Fatal(err)
	}

	if size != 4200 {
		t.Fatalf("expected 4200, got %d", size)
	}
}

func TestEntrySizeEmptyFolder(t *testing.T) {
	dir := t.TempDir()
	folder := filepath.Join(dir, "empty")
	if err := os.MkdirAll(folder, 0o755); err != nil {
		t.Fatal(err)
	}

	fi, err := os.Lstat(folder)
	if err != nil {
		t.Fatal(err)
	}

	size, err := entrySize(folder, fi)
	if err != nil {
		t.Fatal(err)
	}

	if size != 0 {
		t.Fatalf("expected 0, got %d", size)
	}
}

func TestSizeAcrossDisks(t *testing.T) {
	disk1 := t.TempDir()
	disk2 := t.TempDir()
	disk3 := t.TempDir()

	writeSizedFile(t, filepath.Join(disk1, "films", "alien", "alien.mkv"), 1000)
	writeSizedFile(t, filepath.Join(disk2, "films", "alien", "extras.mkv"), 250)

	c := &Core{state: &domain.State{Unraid: &domain.Unraid{Disks: []*domain.Disk{
		{Name: "disk1", Path: disk1},
		{Name: "disk2", Path: disk2},
		{Name: "disk3", Path: disk3},
	}}}}

	sizes := c.Size("/mnt/user/films/alien")

	if len(sizes.Disks) != 2 {
		t.Fatalf("expected 2 disks, got %d (%v)", len(sizes.Disks), sizes.Disks)
	}

	if sizes.Disks["disk1"] != 1000 {
		t.Fatalf("expected disk1 to hold 1000, got %d", sizes.Disks["disk1"])
	}

	if sizes.Disks["disk2"] != 250 {
		t.Fatalf("expected disk2 to hold 250, got %d", sizes.Disks["disk2"])
	}

	if sizes.Total != 1250 {
		t.Fatalf("expected total 1250, got %d", sizes.Total)
	}
}

func TestSizeMissingEntry(t *testing.T) {
	disk1 := t.TempDir()

	c := &Core{state: &domain.State{Unraid: &domain.Unraid{Disks: []*domain.Disk{
		{Name: "disk1", Path: disk1},
	}}}}

	sizes := c.Size("/mnt/user/does/not/exist")

	if len(sizes.Disks) != 0 || sizes.Total != 0 {
		t.Fatalf("expected empty result, got %v", sizes)
	}
}
