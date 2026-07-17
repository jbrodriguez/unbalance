package core

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/cskr/pubsub"

	"unbalance/daemon/common"
	"unbalance/daemon/domain"
)

func planTestCore(t *testing.T) (*Core, *domain.Disk) {
	t.Helper()

	dir := t.TempDir()

	folder := filepath.Join(dir, "films")
	if err := os.MkdirAll(folder, 0o755); err != nil {
		t.Fatal(err)
	}

	if err := os.WriteFile(filepath.Join(folder, "alien.mkv"), make([]byte, 100), 0o644); err != nil {
		t.Fatal(err)
	}

	c := &Core{
		ctx:   &domain.Context{Hub: pubsub.New(common.ChanCapacity)},
		state: &domain.State{Status: common.OpScatterPlan},
	}

	return c, &domain.Disk{Name: "disk1", Path: dir}
}

func TestGetItemsAndIssuesStopsWhenCancelled(t *testing.T) {
	c, disk := planTestCore(t)

	c.stopped.Store(true)

	items, _, _, _, _ := c.getItemsAndIssues(context.Background(), common.OpScatterPlan, 4096, reItems, reStat, []*domain.Disk{disk}, []string{"films"})

	if len(items) != 0 {
		t.Fatalf("expected no items after cancellation, got %d", len(items))
	}
}

func TestGetItemsAndIssuesAbortsOnContextCancel(t *testing.T) {
	c, disk := planTestCore(t)

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	items, _, _, _, _ := c.getItemsAndIssues(ctx, common.OpScatterPlan, 4096, reItems, reStat, []*domain.Disk{disk}, []string{"films"})

	if len(items) != 0 {
		t.Fatalf("expected no items after context cancellation, got %d", len(items))
	}
}

func TestGetItemsAndIssuesScansWhenNotCancelled(t *testing.T) {
	c, disk := planTestCore(t)

	items, _, _, _, _ := c.getItemsAndIssues(context.Background(), common.OpScatterPlan, 4096, reItems, reStat, []*domain.Disk{disk}, []string{"films"})

	if len(items) == 0 {
		t.Fatalf("expected items from an uncancelled scan")
	}
}

func TestPlanCancelledResetsStatus(t *testing.T) {
	c, _ := planTestCore(t)

	c.planCancelled(common.EventScatterPlanCancelled)

	if c.state.Status != common.OpNeutral {
		t.Fatalf("expected neutral status, got %d", c.state.Status)
	}
}

func TestCancelPlanContextKillsScan(t *testing.T) {
	c, disk := planTestCore(t)

	ctx := c.newPlanContext()

	// simulate the user pressing stop while the scan is in flight
	c.stopped.Store(true)
	c.cancelPlanContext()

	items, _, _, _, _ := c.getItemsAndIssues(ctx, common.OpScatterPlan, 4096, reItems, reStat, []*domain.Disk{disk}, []string{"films"})

	if len(items) != 0 {
		t.Fatalf("expected no items after stop, got %d", len(items))
	}
}
