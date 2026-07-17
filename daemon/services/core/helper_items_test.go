package core

import (
	"context"
	"os"
	"path/filepath"
	"testing"
	"time"

	"unbalance/daemon/domain"
)

func noTick(int) {}

func TestGetItemsSingleFileBlocksUsed(t *testing.T) {
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "movie.mkv"), make([]byte, 5000), 0o644); err != nil {
		t.Fatal(err)
	}

	items, total, err := getItems(context.Background(), 4096, reItems, dir, "movie.mkv", noTick)
	if err != nil {
		t.Fatal(err)
	}

	if len(items) != 1 || total != 5000 {
		t.Fatalf("expected one 5000-byte item, got %d items, total %d", len(items), total)
	}

	// 5000 bytes at 4096-byte blocks round up to 2 blocks; a zero value
	// makes the block-based allocator misplace single files (#128)
	if items[0].BlocksUsed != 2 {
		t.Fatalf("expected BlocksUsed 2, got %d", items[0].BlocksUsed)
	}
}

func TestGetItemsEmptyFolderBlocksUsed(t *testing.T) {
	dir := t.TempDir()
	if err := os.MkdirAll(filepath.Join(dir, "empty"), 0o755); err != nil {
		t.Fatal(err)
	}

	items, _, err := getItems(context.Background(), 4096, reItems, dir, "empty", noTick)
	if err != nil {
		t.Fatal(err)
	}

	if len(items) != 1 || items[0].Size != 1 || items[0].BlocksUsed != 1 {
		t.Fatalf("expected the empty-folder sentinel with Size 1 / BlocksUsed 1, got %+v", items[0])
	}
}

func TestGetItemsMissingFolderErrors(t *testing.T) {
	dir := t.TempDir()

	_, _, err := getItems(context.Background(), 4096, reItems, dir, "does-not-exist", noTick)
	if err == nil {
		t.Fatalf("expected an error for a missing folder")
	}
}

func TestProgressGuardsDegenerateInput(t *testing.T) {
	percent, left, speed := progress(0, 0, time.Second)
	if percent != 0 || left != 0 || speed != 0 {
		t.Fatalf("expected zeros for zero bytesToTransfer, got %f %s %f", percent, left, speed)
	}

	percent, left, speed = progress(100, 50, 0)
	if percent != 0 || left != 0 || speed != 0 {
		t.Fatalf("expected zeros for zero elapsed, got %f %s %f", percent, left, speed)
	}

	percent, _, _ = progress(100, 200, time.Second)
	if percent != 100 {
		t.Fatalf("expected clamped percent 100, got %f", percent)
	}
}

func TestRemoveItemsPreservesOrder(t *testing.T) {
	a := &domain.Item{Name: "a"}
	b := &domain.Item{Name: "b"}
	c := &domain.Item{Name: "c"}

	remaining := removeItems([]*domain.Item{a, b, c}, []*domain.Item{b})

	if len(remaining) != 2 || remaining[0] != a || remaining[1] != c {
		t.Fatalf("expected [a c], got %+v", remaining)
	}
}
