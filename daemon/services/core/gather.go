package core

import (
	"fmt"
	"path/filepath"
	"strings"
	"time"

	"github.com/teris-io/shortid"

	"unbalance/daemon/algorithm"
	"unbalance/daemon/common"
	"unbalance/daemon/domain"
	"unbalance/daemon/lib"
	"unbalance/daemon/logger"
)

func (c *Core) gatherPlanPrepare(setup domain.GatherSetup) {
	now := time.Now()

	if c.state.Status != common.OpNeutral {
		logger.Yellow("unbalance is busy: %d", c.state.Status)
		return
	}

	c.state.Status = common.OpGatherPlan
	c.state.Unraid = c.refreshUnraid()

	plan := &domain.Plan{
		Started:       now,
		ChosenFolders: setup.Selected,
		VDisks:        make(map[string]*domain.VDisk),
	}

	for _, disk := range c.state.Unraid.Disks {
		plan.VDisks[disk.Path] = &domain.VDisk{
			Path:        disk.Path,
			CurrentFree: disk.Free,
			PlannedFree: disk.Free,
			Bin:         nil,
		}
	}

	go c.gatherPlan(plan)
}

func (c *Core) gatherPlan(plan *domain.Plan) {
	c.gatherPlanStart(plan)
	c.gatherPlanEnd(plan)
}

func (c *Core) gatherPlanStart(plan *domain.Plan) {
	logger.Blue("Running gather planner ...")

	packet := &domain.Packet{Topic: common.EventGatherPlanStarted, Payload: "Planning started"}
	c.ctx.Hub.Pub(packet, "socket:broadcast")

	c.printDisks(c.state.Unraid.Disks, c.state.Unraid.BlockSize)

	items, ownerIssue, groupIssue, folderIssue, fileIssue := c.getItemsAndIssues(c.state.Status, c.state.Unraid.BlockSize, reItems, reStat, c.state.Unraid.Disks, plan.ChosenFolders)

	// // no items found, no sense going on, just end this planning
	// if len(items) == 0 {
	// 	p.endPlan(state.Status, plan, state.Unraid.Disks, items, toBeTransferred)
	// 	p.bus.Pub(&pubsub.Message{Payload: plan}, common.IntScatterPlanFinished)
	// 	return
	// }

	plan.OwnerIssue = ownerIssue
	plan.GroupIssue = groupIssue
	plan.FolderIssue = folderIssue
	plan.FileIssue = fileIssue

	logger.Blue("gatherPlan:items(%d)", len(items))

	for _, item := range items {
		logger.Blue("gatherPlan:found(%s):size(%d)", filepath.Join(item.Location, item.Path), item.Size)

		msg := fmt.Sprintf("Found %s (%s)", filepath.Join(item.Location, item.Path), lib.ByteSize(item.Size))
		packet = &domain.Packet{Topic: common.EventGatherPlanProgress, Payload: msg}
		c.ctx.Hub.Pub(packet, "socket:broadcast")
	}

	logger.Blue("gatherPlan:issues:owner(%d),group(%d),folder(%d),file(%d)", plan.OwnerIssue, plan.GroupIssue, plan.FolderIssue, plan.FileIssue)

	// Initialize fields
	plan.BytesToTransfer = 0

	for _, disk := range c.state.Unraid.Disks {
		msg := fmt.Sprintf("Trying to allocate items to %s ...", disk.Name)
		packet = &domain.Packet{Topic: common.EventGatherPlanProgress, Payload: msg}
		c.ctx.Hub.Pub(packet, "socket:broadcast")
		logger.Blue("gatherPlan:%s", msg)

		reserved := c.getReservedAmount(disk.Size)

		ceil := lib.Max(common.ReservedSpace, reserved)
		logger.Blue("scatterPlan:ItemsLeft(%d):ReservedSpace(%d)", len(items), ceil)

		packer := algorithm.NewGreedy(disk, items, ceil, c.state.Unraid.BlockSize)
		bin := packer.FitAll()
		if bin != nil {
			plan.VDisks[disk.Path].Bin = bin
			plan.VDisks[disk.Path].PlannedFree -= bin.Size

			plan.BytesToTransfer += bin.Size
		}
	}

	c.endPlan(c.state.Status, plan, c.state.Unraid.Disks, items, make([]*domain.Item, 0))
	// p.bus.Pub(&pubsub.Message{Payload: plan}, common.IntScatterPlanFinished)
	c.gatherPlanEnd(plan)
}

func (c *Core) gatherPlanEnd(plan *domain.Plan) {
	c.state.Status = common.OpNeutral

	packet := &domain.Packet{Topic: common.EventGatherPlanEnded, Payload: plan}
	c.ctx.Hub.Pub(packet, "socket:broadcast")

	// TODO: finish this implementation
	// // only send the perm issue msg if there's actually some work to do (BytesToTransfer > 0)
	// // and there actually perm issues
	// if plan.BytesToTransfer > 0 && (plan.OwnerIssue+plan.GroupIssue+plan.FolderIssue+plan.FileIssue > 0) {
	// 	outbound = &dto.Packet{Topic: common.WsGatherPlanIssues, Payload: fmt.Sprintf("%d|%d|%d|%d", plan.OwnerIssue, plan.GroupIssue, plan.FolderIssue, plan.FileIssue)}
	// 	c.bus.Pub(&pubsub.Message{Payload: outbound}, "socket:broadcast")
	// }
}

// // GATHER TRANSFER
func (c *Core) gatherMove(plan domain.Plan) {
	c.state.Status = common.OpGatherMove
	c.state.Unraid = c.refreshUnraid()
	c.state.Operation = c.createGatherOperation(plan)

	go c.runOperation("Move")
}

func (c *Core) createGatherOperation(plan domain.Plan) *domain.Operation {
	operation := &domain.Operation{
		ID:     shortid.MustGenerate(),
		OpKind: c.state.Status,
		DryRun: c.ctx.DryRun,
	}

	operation.RsyncArgs = append([]string{common.RsyncArgs}, c.ctx.RsyncArgs...)

	// user may have changed dry-run setting, adjust for it
	if operation.DryRun {
		operation.RsyncArgs = append(operation.RsyncArgs, "--dry-run")
	}
	operation.RsyncStrArgs = strings.Join(operation.RsyncArgs, " ")

	operation.Commands = make([]*domain.Command, 0)

	for _, disk := range c.state.Unraid.Disks {
		vdisk := plan.VDisks[disk.Path]

		// only one disk will be destination (target)
		if vdisk.Path != plan.Target {
			continue
		}

		// user chose a target disk, adjust bytestotransfer to the size of its bin, since
		// that's the amount of data we need to transfer. Also remove bin from all other disks,
		operation.BytesToTransfer = vdisk.Bin.Size

		for _, item := range vdisk.Bin.Items {
			var entry string

			// this a double check. item.Path should be in the form of films/bluray or tvshows/Billions
			// if for some reason it starts with a "/", we strip it
			if item.Path[0] == filepath.Separator {
				entry = item.Path[1:]
			} else {
				entry = item.Path
			}

			operation.Commands = append(operation.Commands, &domain.Command{
				ID:     shortid.MustGenerate(),
				Src:    item.Location,
				Dst:    vdisk.Path + string(filepath.Separator),
				Entry:  entry,
				Size:   item.Size,
				Status: common.CmdPending,
			})
		}
	}

	return operation
}
