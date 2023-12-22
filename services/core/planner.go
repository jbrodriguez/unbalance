package core

// "unbalance/1old/server/algorithm"
// "unbalance/1old/server/dto"
// "unbalance/1old/server/lib"
// "unbalance/1old/server/algorithm"

// "github.com/jbrodriguez/mlog"
// "github.com/phonkee/go-pubsub"
// "fmt"
// "io/ioutil"
// "math"
// "os"
// "path/filepath"
// "regexp"
// "sort"
// "strconv"
// "strings"
// "time"
// "unbalance/algorithm"
// "unbalance/common"
// "unbalance/domain"
// "unbalance/dto"
// "unbalance/lib"
// "github.com/jbrodriguez/actor"
// "github.com/jbrodriguez/mlog"
// "github.com/jbrodriguez/pubsub"

// Planner -
// type Planner struct {
// 	bus      *pubsub.PubSub
// 	settings *lib.Settings
// 	actor    *actor.Actor

// 	reItems *regexp.Regexp
// 	reStat  *regexp.Regexp
// }

// // NewPlanner -
// func NewPlanner(bus *pubsub.PubSub, settings *lib.Settings) *Planner {
// 	plan := &Planner{
// 		bus:      bus,
// 		settings: settings,
// 		actor:    actor.NewActor(bus),
// 	}

// 	plan.reItems = regexp.MustCompile(`(\d+)\s+(.*?)$`)
// 	plan.reStat = regexp.MustCompile(`[-dclpsbD]([-rwxsS]{3})([-rwxsS]{3})([-rwxtT]{3})\|(.*?)\:(.*?)\|(.*?)\|(.*)`)

// 	return plan
// }

// Start -
// func (c *Core) Start() (err error) {
// 	mlog.Info("starting service Planner ...")

// 	p.actor.Register(common.IntScatterPlan, p.scatter)
// 	p.actor.Register(common.IntGatherPlan, p.gather)

// 	go p.actor.React()

// 	return nil
// }

// Stop -
// func (c *Core) Stop() {
// 	mlog.Info("stopped service Planner ...")
// }

// func (c *Core) gather(msg *pubsub.Message) {
// 	state := msg.Payload.(*domain.State)

// 	mlog.Info("Running gather planner ...")

// 	plan := state.Plan
// 	plan.Started = time.Now()

// 	outbound := &dto.Packet{Topic: common.WsGatherPlanStarted, Payload: "Planning Started"}
// 	p.bus.Pub(&pubsub.Message{Payload: outbound}, "socket:broadcast")

// 	p.printDisks(state.Unraid.Disks, state.Unraid.BlockSize)

// 	items, ownerIssue, groupIssue, folderIssue, fileIssue := p.getItemsAndIssues(state.Status, state.Unraid.BlockSize, p.reItems, p.reStat, state.Unraid.Disks, plan.ChosenFolders)

// 	// no items found, no sense going on, just end this planning
// 	if len(items) == 0 {
// 		p.endPlan(state.Status, plan, state.Unraid.Disks, items, make([]*domain.Item, 0))
// 		p.bus.Pub(&pubsub.Message{Payload: plan}, common.IntScatterPlanFinished)
// 		return
// 	}

// 	plan.OwnerIssue = ownerIssue
// 	plan.GroupIssue = groupIssue
// 	plan.FolderIssue = folderIssue
// 	plan.FileIssue = fileIssue

// 	mlog.Info("gatherPlan:items(%d)", len(items))

// 	for _, item := range items {
// 		mlog.Info("gatherPlan:found(%s):size(%d)", filepath.Join(item.Location, item.Path), item.Size)

// 		msg := fmt.Sprintf("Found %s (%s)", filepath.Join(item.Location, item.Path), lib.ByteSize(item.Size))
// 		outbound = &dto.Packet{Topic: common.WsGatherPlanProgress, Payload: msg}
// 		p.bus.Pub(&pubsub.Message{Payload: outbound}, "socket:broadcast")
// 	}

// 	mlog.Info("gatherPlan:issues:owner(%d),group(%d),folder(%d),file(%d)", plan.OwnerIssue, plan.GroupIssue, plan.FolderIssue, plan.FileIssue)

// 	// Initialize fields
// 	plan.BytesToTransfer = 0

// 	for _, disk := range state.Unraid.Disks {
// 		msg := fmt.Sprintf("Trying to allocate items to %s ...", disk.Name)
// 		outbound = &dto.Packet{Topic: common.WsGatherPlanProgress, Payload: msg}
// 		p.bus.Pub(&pubsub.Message{Payload: outbound}, "socket:broadcast")
// 		mlog.Info("gatherPlan:%s", msg)

// 		reserved := p.getReservedAmount(disk.Size)

// 		ceil := lib.Max(lib.ReservedSpace, reserved)
// 		mlog.Info("gatherPlan:ItemsLeft(%d):ReservedSpace(%d)", len(items), ceil)

// 		packer := algorithm.NewGreedy(disk, items, ceil, state.Unraid.BlockSize)
// 		bin := packer.FitAll()
// 		if bin != nil {
// 			plan.VDisks[disk.Path].Bin = bin
// 			plan.VDisks[disk.Path].PlannedFree -= bin.Size

// 			plan.BytesToTransfer += bin.Size
// 		}
// 	}

// 	p.endPlan(state.Status, plan, state.Unraid.Disks, items, make([]*domain.Item, 0))
// 	p.bus.Pub(&pubsub.Message{Payload: plan}, common.IntGatherPlanFinished)
// }
