package core

import (
	"fmt"
	"path/filepath"
	"slices"
	"sort"
	"strings"
	"time"

	"github.com/teris-io/shortid"

	"unbalance/daemon/algorithm"
	"unbalance/daemon/common"
	"unbalance/daemon/domain"
	"unbalance/daemon/lib"
	"unbalance/daemon/logger"
)

// SCATTER PLANNER
func (c *Core) scatterPlanPrepare(setup domain.ScatterSetup) {
	now := time.Now()

	if c.state.Status != common.OpNeutral {
		logger.Yellow("unbalance is busy: %d", c.state.Status)
		return
	}

	c.state.Status = common.OpScatterPlan
	c.state.Unraid = c.refreshUnraid()

	plan := &domain.Plan{
		Started:       now,
		ChosenFolders: setup.Selected,
		VDisks:        make(map[string]*domain.VDisk),
	}

	targets := make([]string, 0)
	for _, target := range setup.Targets {
		targets = append(targets, filepath.Join("/", "mnt", target))
	}

	for _, disk := range c.state.Unraid.Disks {
		plan.VDisks[disk.Path] = &domain.VDisk{
			Path:        disk.Path,
			CurrentFree: disk.Free,
			PlannedFree: disk.Free,
			Bin:         nil,
			Src:         filepath.Join("/", "mnt", setup.Source) == disk.Path,
			Dst:         slices.Contains(targets, disk.Path),
		}
	}

	// logger.Green("%+v", c.state.Plan)
	// for _, disk := range c.state.Unraid.Disks {
	// 	logger.Green("%+v", c.state.Plan.VDisks[disk.Path])
	// }

	// c.bus.Pub(&pubsub.Message{Payload: c.state.Plan}, common.IntScatterPlanStarted)
	// c.bus.Pub(&pubsub.Message{Payload: c.state}, common.IntScatterPlanStarted)

	// c.actor.Tell(common.IntScatterPlan, c.state)
	go c.scatterPlan(plan)
}

func (c *Core) scatterPlan(plan *domain.Plan) {
	c.scatterPlanStart(plan)
	c.scatterPlanEnd(plan)
}

func (c *Core) scatterPlanStart(plan *domain.Plan) {
	logger.Blue("Running scatter planner ...")

	// plan := c.state.Plan
	// c.state.Plan.Started = time.Now()

	// create two slices
	// one of source disks, the other of destinations disks
	// in scatter srcDisks will contain only one element
	srcDisk, dstDisks := getSourceAndDestinationDisks(c.state.Unraid.Disks, plan)

	// get dest disks with more free space to the top
	sort.Slice(dstDisks, func(i, j int) bool { return dstDisks[i].Free < dstDisks[j].Free })

	// some logging
	logger.Blue("scatterPlan:source:(%+v)", srcDisk)
	for _, disk := range dstDisks {
		logger.Blue("scatterPlan:dest:(%+v)", disk)
	}

	// outbound := &dto.Packet{Topic: common.WsScatterPlanStarted, Payload: "Planning started"}
	// p.bus.Pub(&pubsub.Message{Payload: outbound}, "socket:broadcast")
	packet := &domain.Packet{Topic: common.EventScatterPlanStarted, Payload: "Planning started"}
	c.ctx.Hub.Pub(packet, "socket:broadcast")

	c.printDisks(c.state.Unraid.Disks, c.state.Unraid.BlockSize)

	items, ownerIssue, groupIssue, folderIssue, fileIssue := c.getItemsAndIssues(c.state.Status, c.state.Unraid.BlockSize, reItems, reStat, []*domain.Disk{srcDisk}, plan.ChosenFolders)

	toBeTransferred := make([]*domain.Item, 0)

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

	logger.Blue("scatterPlan:items(%d)", len(items))

	for _, item := range items {
		logger.Blue("scatterPlan:found(%s):size(%d)", filepath.Join(item.Location, item.Path), item.Size)

		msg := fmt.Sprintf("Found %s (%s)", filepath.Join(item.Location, item.Path), lib.ByteSize(item.Size))
		packet = &domain.Packet{Topic: common.EventScatterPlanProgress, Payload: msg}
		c.ctx.Hub.Pub(packet, "socket:broadcast")
	}

	logger.Blue("scatterPlan:issues:owner(%d),group(%d),folder(%d),file(%d)", plan.OwnerIssue, plan.GroupIssue, plan.FolderIssue, plan.FileIssue)

	// Initialize fields
	plan.BytesToTransfer = 0

	for _, disk := range dstDisks {
		msg := fmt.Sprintf("Trying to allocate items to %s ...", disk.Name)
		packet = &domain.Packet{Topic: common.EventScatterPlanProgress, Payload: msg}
		c.ctx.Hub.Pub(packet, "socket:broadcast")
		logger.Blue("scatterPlan:%s", msg)

		reserved := c.getReservedAmount(disk.Size)

		ceil := lib.Max(common.ReservedSpace, reserved)
		logger.Blue("scatterPlan:ItemsLeft(%d):ReservedSpace(%d)", len(items), ceil)

		packer := algorithm.NewKnapsack(disk, items, ceil, c.state.Unraid.BlockSize)
		bin := packer.BestFit()
		if bin != nil {
			plan.VDisks[disk.Path].Bin = bin
			plan.VDisks[disk.Path].PlannedFree -= bin.Size
			plan.VDisks[srcDisk.Path].PlannedFree += bin.Size

			plan.BytesToTransfer += bin.Size

			toBeTransferred = append(toBeTransferred, bin.Items...)
			items = removeItems(items, bin.Items)
		}
	}

	c.endPlan(c.state.Status, plan, c.state.Unraid.Disks, items, toBeTransferred)
	// p.bus.Pub(&pubsub.Message{Payload: plan}, common.IntScatterPlanFinished)
	c.scatterPlanEnd(plan)
}

func (c *Core) scatterPlanEnd(plan *domain.Plan) {
	c.state.Status = common.OpNeutral

	packet := &domain.Packet{Topic: common.EventScatterPlanEnded, Payload: plan}
	c.ctx.Hub.Pub(packet, "socket:broadcast")

	// TODO: finish implementation
	// // only send the perm issue msg if there's actually some work to do (BytesToTransfer > 0)
	// // and there actually perm issues
	// if plan.BytesToTransfer > 0 && (plan.OwnerIssue+plan.GroupIssue+plan.FolderIssue+plan.FileIssue > 0) {
	// 	outbound = &dto.Packet{Topic: common.WsScatterPlanIssues, Payload: fmt.Sprintf("%d|%d|%d|%d", plan.OwnerIssue, plan.GroupIssue, plan.FolderIssue, plan.FileIssue)}
	// 	c.bus.Pub(&pubsub.Message{Payload: outbound}, "socket:broadcast")
	// }
}

// // SCATTER TRANSFER
func (c *Core) scatterMove(plan domain.Plan) {
	c.state.Status = common.OpScatterMove
	c.state.Unraid = c.refreshUnraid()
	c.state.Operation = c.createScatterOperation(plan)

	go c.runOperation("Move")
}

func (c *Core) scatterCopy(plan domain.Plan) {
	c.state.Status = common.OpScatterCopy
	c.state.Unraid = c.refreshUnraid()
	c.state.Operation = c.createScatterOperation(plan)

	go c.runOperation("Copy")
}

func (c *Core) createScatterOperation(plan domain.Plan) *domain.Operation {
	operation := &domain.Operation{
		ID:              shortid.MustGenerate(),
		OpKind:          c.state.Status,
		BytesToTransfer: plan.BytesToTransfer,
		DryRun:          c.ctx.DryRun,
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

		if vdisk.Bin == nil || vdisk.Src {
			continue
		}

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
