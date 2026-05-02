package core

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/teris-io/shortid"

	"unbalance/daemon/common"
	"unbalance/daemon/domain"
)

const pendingPlanTTL = 2 * time.Hour

type planFlow string

const (
	planFlowScatter planFlow = "scatter"
	planFlowGather  planFlow = "gather"
)

type planTicket struct {
	ID        string
	Flow      planFlow
	Plan      domain.Plan
	CreatedAt time.Time
	ExpiresAt time.Time
}

type operationRef struct {
	OperationID string `json:"operationID"`
}

type commandRef struct {
	OperationID string `json:"operationID"`
	CommandID   string `json:"commandID"`
}

type planRef struct {
	PlanID string `json:"planID"`
	Target string `json:"target"`
}

func newPlanID() string {
	return shortid.MustGenerate()
}

func (c *Core) storePendingPlan(flow planFlow, plan *domain.Plan) (*domain.Plan, error) {
	if plan == nil {
		return nil, fmt.Errorf("missing plan")
	}

	stored, err := clonePlan(*plan)
	if err != nil {
		return nil, err
	}

	if stored.ID == "" {
		stored.ID = newPlanID()
	}

	now := time.Now()
	ticket := &planTicket{
		ID:        stored.ID,
		Flow:      flow,
		Plan:      stored,
		CreatedAt: now,
		ExpiresAt: now.Add(pendingPlanTTL),
	}

	c.pendingPlansMu.Lock()
	defer c.pendingPlansMu.Unlock()
	c.pruneExpiredPendingPlansLocked(now)
	c.pendingPlans[ticket.ID] = ticket

	view, err := clonePlan(stored)
	if err != nil {
		return nil, err
	}

	return &view, nil
}

func (c *Core) takePendingPlan(id string, flow planFlow) (domain.Plan, error) {
	if id == "" {
		return domain.Plan{}, fmt.Errorf("missing planID")
	}

	now := time.Now()
	c.pendingPlansMu.Lock()
	defer c.pendingPlansMu.Unlock()
	c.pruneExpiredPendingPlansLocked(now)

	ticket, ok := c.pendingPlans[id]
	if !ok {
		return domain.Plan{}, fmt.Errorf("pending plan not found or expired: %s", id)
	}
	if ticket.Flow != flow {
		return domain.Plan{}, fmt.Errorf("pending plan %s is for %s, not %s", id, ticket.Flow, flow)
	}

	delete(c.pendingPlans, id)
	return clonePlan(ticket.Plan)
}

func (c *Core) pruneExpiredPendingPlansLocked(now time.Time) {
	for id, ticket := range c.pendingPlans {
		if !ticket.ExpiresAt.After(now) {
			delete(c.pendingPlans, id)
		}
	}
}

func clonePlan(plan domain.Plan) (domain.Plan, error) {
	data, err := json.Marshal(plan)
	if err != nil {
		return domain.Plan{}, err
	}

	var cloned domain.Plan
	if err := json.Unmarshal(data, &cloned); err != nil {
		return domain.Plan{}, err
	}

	return cloned, nil
}

func cloneOperation(operation domain.Operation) (domain.Operation, error) {
	data, err := json.Marshal(operation)
	if err != nil {
		return domain.Operation{}, err
	}

	var cloned domain.Operation
	if err := json.Unmarshal(data, &cloned); err != nil {
		return domain.Operation{}, err
	}

	return cloned, nil
}

func (c *Core) historyOperation(id string) (domain.Operation, error) {
	if id == "" {
		return domain.Operation{}, fmt.Errorf("missing operationID")
	}
	if c.state.History == nil || c.state.History.Items == nil {
		return domain.Operation{}, fmt.Errorf("history is unavailable")
	}

	operation, ok := c.state.History.Items[id]
	if !ok || operation == nil {
		return domain.Operation{}, fmt.Errorf("operation not found: %s", id)
	}

	return cloneOperation(*operation)
}

func findCommand(operation *domain.Operation, id string) (*domain.Command, error) {
	if operation == nil {
		return nil, fmt.Errorf("missing operation")
	}
	if id == "" {
		return nil, fmt.Errorf("missing commandID")
	}

	for _, command := range operation.Commands {
		if command != nil && command.ID == id {
			return command, nil
		}
	}

	return nil, fmt.Errorf("command not found: %s", id)
}

func (c *Core) publishOperationError(format string, args ...any) {
	msg := fmt.Sprintf(format, args...)
	packet := &domain.Packet{Topic: common.EventOperationError, Payload: msg}
	c.ctx.Hub.Pub(packet, "socket:broadcast")
}
