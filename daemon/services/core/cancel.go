package core

import (
	"context"
)

// newPlanContext arms cancellation for a starting plan: the returned context
// is cancelled when the user stops the plan, killing any in-flight scan
// process immediately.
func (c *Core) newPlanContext() context.Context {
	ctx, cancel := context.WithCancel(context.Background())

	c.planMu.Lock()
	c.planCancel = cancel
	c.planMu.Unlock()

	return ctx
}

// cancelPlanContext kills the in-flight scan process, if any.
func (c *Core) cancelPlanContext() {
	c.planMu.Lock()
	if c.planCancel != nil {
		c.planCancel()
	}
	c.planMu.Unlock()
}

// clearPlanContext releases the plan context after a plan ends or is
// cancelled.
func (c *Core) clearPlanContext() {
	c.planMu.Lock()
	if c.planCancel != nil {
		c.planCancel()
		c.planCancel = nil
	}
	c.planMu.Unlock()
}
