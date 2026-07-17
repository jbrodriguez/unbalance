package algorithm

import (
	"testing"

	"unbalance/daemon/domain"
)

func TestBestFitBytesSkipsDiskBelowReserve(t *testing.T) {
	disk := &domain.Disk{Name: "disk1", Path: "/mnt/disk1", Free: 500 * 1024 * 1024}
	items := []*domain.Item{{Name: "films", Size: 100 * 1024 * 1024}}

	// reserve larger than the disk's free space: without the guard the
	// unsigned subtraction wraps and the item is allocated anyway
	k := NewKnapsack(disk, items, 1024*1024*1024, 0)

	if bin := k.BestFit(); bin != nil {
		t.Fatalf("expected no bin on a disk below the reserve, got %+v", bin)
	}
}

func TestBestFitBlocksSkipsDiskBelowReserve(t *testing.T) {
	disk := &domain.Disk{Name: "disk1", Path: "/mnt/disk1", Free: 500 * 1024 * 1024, BlocksFree: 100}
	items := []*domain.Item{{Name: "films", Size: 100 * 4096, BlocksUsed: 100}}

	// reserve of 1GiB at 4KiB blocks is far more than 100 free blocks
	k := NewKnapsack(disk, items, 1024*1024*1024, 4096)

	if bin := k.BestFit(); bin != nil {
		t.Fatalf("expected no bin on a disk below the block reserve, got %+v", bin)
	}
}

func TestBestFitAllocatesWhenSpaceAvailable(t *testing.T) {
	disk := &domain.Disk{Name: "disk1", Path: "/mnt/disk1", Free: 10 * 1024 * 1024 * 1024, BlocksFree: 2621440}
	items := []*domain.Item{{Name: "films", Size: 1024 * 1024 * 1024, BlocksUsed: 262144}}

	k := NewKnapsack(disk, items, 1024*1024*1024, 4096)

	bin := k.BestFit()
	if bin == nil || len(bin.Items) != 1 {
		t.Fatalf("expected the item to be allocated, got %+v", bin)
	}
}
