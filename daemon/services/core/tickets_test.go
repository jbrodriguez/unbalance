package core

import (
	"testing"
	"time"

	"unbalance/daemon/domain"
)

func testCoreWithTickets() *Core {
	return &Core{pendingPlans: make(map[string]*planTicket)}
}

func TestPendingPlanTicketStoresCloneAndConsumesByID(t *testing.T) {
	c := testCoreWithTickets()
	plan := &domain.Plan{
		BytesToTransfer: 10,
		VDisks: map[string]*domain.VDisk{
			"/mnt/disk1": {Path: "/mnt/disk1"},
		},
	}

	view, err := c.storePendingPlan(planFlowScatter, plan)
	if err != nil {
		t.Fatalf("storePendingPlan: %s", err)
	}
	if view.ID == "" {
		t.Fatalf("expected generated plan id")
	}

	// Mutating the original or returned view must not affect server-owned plan data.
	plan.BytesToTransfer = 99
	view.BytesToTransfer = 77

	stored, err := c.takePendingPlan(view.ID, planFlowScatter)
	if err != nil {
		t.Fatalf("takePendingPlan: %s", err)
	}
	if stored.BytesToTransfer != 10 {
		t.Fatalf("stored bytes = %d, want 10", stored.BytesToTransfer)
	}

	if _, err := c.takePendingPlan(view.ID, planFlowScatter); err == nil {
		t.Fatalf("expected consumed plan ticket to be unavailable")
	}
}

func TestPendingPlanTicketRejectsWrongFlowWithoutConsuming(t *testing.T) {
	c := testCoreWithTickets()
	view, err := c.storePendingPlan(planFlowScatter, &domain.Plan{})
	if err != nil {
		t.Fatalf("storePendingPlan: %s", err)
	}

	if _, err := c.takePendingPlan(view.ID, planFlowGather); err == nil {
		t.Fatalf("expected wrong flow to be rejected")
	}

	if _, err := c.takePendingPlan(view.ID, planFlowScatter); err != nil {
		t.Fatalf("expected ticket to remain after wrong-flow rejection: %s", err)
	}
}

func TestPendingPlanTicketPrunesExpiredPlans(t *testing.T) {
	c := testCoreWithTickets()
	c.pendingPlans["expired"] = &planTicket{ID: "expired", Flow: planFlowScatter, ExpiresAt: time.Now().Add(-time.Minute)}

	if _, err := c.takePendingPlan("expired", planFlowScatter); err == nil {
		t.Fatalf("expected expired ticket to be rejected")
	}
	if len(c.pendingPlans) != 0 {
		t.Fatalf("expected expired ticket to be pruned")
	}
}

func TestHistoryOperationReturnsClone(t *testing.T) {
	c := &Core{state: &domain.State{History: &domain.History{Items: map[string]*domain.Operation{
		"op-1": {ID: "op-1", Commands: []*domain.Command{{ID: "cmd-1", Entry: "original"}}},
	}}}}

	operation, err := c.historyOperation("op-1")
	if err != nil {
		t.Fatalf("historyOperation: %s", err)
	}
	operation.Commands[0].Entry = "mutated"

	if got := c.state.History.Items["op-1"].Commands[0].Entry; got != "original" {
		t.Fatalf("history operation was mutated through clone, got %q", got)
	}
}
