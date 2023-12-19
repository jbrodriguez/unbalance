package core

import (
	"fmt"
	"io/ioutil"
	"math"
	"os"
	"path/filepath"
	"regexp"
	"slices"
	"sort"
	"strconv"
	"strings"
	"time"

	// "unbalance/1old/server/algorithm"
	// "unbalance/1old/server/dto"
	// "unbalance/1old/server/lib"
	// "unbalance/1old/server/algorithm"

	"unbalance/algorithm"
	"unbalance/common"
	"unbalance/domain"
	"unbalance/lib"
	"unbalance/logger"
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
)

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

// SCATTER PLANNER
func (c *Core) scatterPlanHandler() {
	for payload := range c.scatterPlanChan {
		now := time.Now()

		if c.state.Status != common.OpNeutral {
			logger.Yellow("unbalance is busy: %d", c.state.Status)
			continue
		}

		var setup domain.ScatterSetup
		err := lib.Bind(payload, &setup)
		if err != nil {
			logger.Red("unable to unmarshal packet: %+v (%s)", payload, err)
			continue
		}

		c.state.Status = common.OpScatterPlanning
		c.state.Unraid = c.refreshUnraid()

		c.state.Plan = &domain.Plan{
			Started:       now,
			ChosenFolders: setup.Selected,
			VDisks:        make(map[string]*domain.VDisk),
		}

		targets := make([]string, 0)
		for _, target := range setup.Targets {
			targets = append(targets, filepath.Join("/", "mnt", target))
		}

		for _, disk := range c.state.Unraid.Disks {
			c.state.Plan.VDisks[disk.Path] = &domain.VDisk{
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
		go c.scatterPlan()
	}
}

func (c *Core) scatterPlan() {
	c.scatterPlanStart()
	c.scatterPlanEnd()
}

func (c *Core) scatterPlanStart() {
	logger.Blue("Running scatter planner ...")

	// plan := c.state.Plan
	// c.state.Plan.Started = time.Now()

	// create two slices
	// one of source disks, the other of destinations disks
	// in scatter srcDisks will contain only one element
	srcDisk, dstDisks := getSourceAndDestinationDisks(c.state.Unraid.Disks, c.state.Plan)

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

	items, ownerIssue, groupIssue, folderIssue, fileIssue := c.getItemsAndIssues(c.state.Status, c.state.Unraid.BlockSize, reItems, reStat, []*domain.Disk{srcDisk}, c.state.Plan.ChosenFolders)

	toBeTransferred := make([]*domain.Item, 0)

	// // no items found, no sense going on, just end this planning
	// if len(items) == 0 {
	// 	p.endPlan(state.Status, plan, state.Unraid.Disks, items, toBeTransferred)
	// 	p.bus.Pub(&pubsub.Message{Payload: plan}, common.IntScatterPlanFinished)
	// 	return
	// }

	c.state.Plan.OwnerIssue = ownerIssue
	c.state.Plan.GroupIssue = groupIssue
	c.state.Plan.FolderIssue = folderIssue
	c.state.Plan.FileIssue = fileIssue

	logger.Blue("scatterPlan:items(%d)", len(items))

	for _, item := range items {
		logger.Blue("scatterPlan:found(%s):size(%d)", filepath.Join(item.Location, item.Path), item.Size)

		msg := fmt.Sprintf("Found %s (%s)", filepath.Join(item.Location, item.Path), lib.ByteSize(item.Size))
		packet = &domain.Packet{Topic: common.EventScatterPlanProgress, Payload: msg}
		c.ctx.Hub.Pub(packet, "socket:broadcast")
	}

	logger.Blue("scatterPlan:issues:owner(%d),group(%d),folder(%d),file(%d)", c.state.Plan.OwnerIssue, c.state.Plan.GroupIssue, c.state.Plan.FolderIssue, c.state.Plan.FileIssue)

	// Initialize fields
	c.state.Plan.BytesToTransfer = 0

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
			c.state.Plan.VDisks[disk.Path].Bin = bin
			c.state.Plan.VDisks[disk.Path].PlannedFree -= bin.Size
			c.state.Plan.VDisks[srcDisk.Path].PlannedFree += bin.Size

			c.state.Plan.BytesToTransfer += bin.Size

			toBeTransferred = append(toBeTransferred, bin.Items...)
			items = removeItems(items, bin.Items)
		}
	}

	c.endPlan(c.state.Status, c.state.Plan, c.state.Unraid.Disks, items, toBeTransferred)
	// p.bus.Pub(&pubsub.Message{Payload: plan}, common.IntScatterPlanFinished)
	c.scatterPlanEnd()
}

func (c *Core) scatterPlanEnd() {
	c.state.Status = common.OpScatterPlan

	packet := &domain.Packet{Topic: common.EventScatterPlanEnded, Payload: "Planning ended"}
	c.ctx.Hub.Pub(packet, "socket:broadcast")
}

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

// COMMON PLANNER
func getSourceAndDestinationDisks(disks []*domain.Disk, plan *domain.Plan) (*domain.Disk, []*domain.Disk) {
	var srcDisk *domain.Disk
	dstDisks := make([]*domain.Disk, 0)

	for _, disk := range disks {
		if plan.VDisks[disk.Path].Src {
			srcDisk = disk
		}

		if plan.VDisks[disk.Path].Dst {
			dstDisks = append(dstDisks, disk)
		}
	}

	return srcDisk, dstDisks
}

func getIssues(re *regexp.Regexp, disk *domain.Disk, path string) (int64, int64, int64, int64, error) {
	var ownerIssue, groupIssue, folderIssue, fileIssue int64

	folder := filepath.Join(disk.Path, path)

	if _, err := os.Stat(folder); os.IsNotExist(err) {
		return ownerIssue, groupIssue, folderIssue, fileIssue, err
	}

	scanFolder := folder + "/."
	cmd := fmt.Sprintf(`find "%s" -exec stat --format "%%A|%%U:%%G|%%F|%%n" {} \;`, scanFolder)

	err := lib.Shell2(cmd, func(line string) {
		result := re.FindStringSubmatch(line)
		if result == nil {
			return
		}

		u := result[1]
		g := result[2]
		o := result[3]
		user := result[4]
		group := result[5]
		kind := result[6]

		perms := u + g + o

		if user != "nobody" {
			ownerIssue++
		}

		if group != "users" {
			groupIssue++
		}

		if kind == "directory" {
			if perms != "rwxrwxrwx" {
				folderIssue++
			}
		} else {
			match := strings.Compare(perms, "r--r--r--") == 0 || strings.Compare(perms, "rw-rw-rw-") == 0
			if !match {
				fileIssue++
			}
		}
	})

	return ownerIssue, groupIssue, folderIssue, fileIssue, err
}

func getItems(blockSize uint64, re *regexp.Regexp, src, folder string) ([]*domain.Item, uint64, error) {
	var total, blocks uint64
	fBlockSize := float64(blockSize)
	srcFolder := filepath.Join(src, folder)

	var fi os.FileInfo
	var err error
	if fi, err = os.Stat(srcFolder); os.IsNotExist(err) {
		return nil, total, err
	}

	if !fi.IsDir() {
		return []*domain.Item{&domain.Item{Name: folder, Size: uint64(fi.Size()), Path: folder, Location: src}}, uint64(fi.Size()), nil
	}

	entries, err := ioutil.ReadDir(srcFolder)
	if err != nil {
		return nil, total, err
	}

	if len(entries) == 0 {
		// Size: 1 is a trick to allow natural processing of this empty folder: if set to zero, many comparison
		// would misinterpret this as a pending transfer and so on
		return []*domain.Item{&domain.Item{Name: srcFolder, Size: 1, Path: folder, Location: src}}, 0, nil
	}

	items := make([]*domain.Item, 0)

	cmd := fmt.Sprintf(`find "%s" ! -name . -prune -exec du -bs {} +`, srcFolder+"/.")

	err = lib.Shell2(cmd, func(line string) {
		result := re.FindStringSubmatch(line)

		size, _ := strconv.ParseInt(result[1], 10, 64)
		total += uint64(size)

		if blockSize > 0 {
			blocks = uint64(math.Ceil(float64(size) / fBlockSize))
		} else {
			blocks = 0
		}

		item := &domain.Item{Name: result[2], Size: uint64(size), Path: filepath.Join(folder, filepath.Base(result[2])), Location: src, BlocksUsed: uint64(blocks)}
		items = append(items, item)
	})

	if err != nil {
		return nil, total, err
	}

	return items, total, err
}

func (c *Core) getItemsAndIssues(status, blockSize uint64, reItems, reStat *regexp.Regexp, disks []*domain.Disk, folders []string) ([]*domain.Item, int64, int64, int64, int64) {
	var ownerIssue, groupIssue, folderIssue, fileIssue int64
	items := make([]*domain.Item, 0)

	// Get owner/permission issues
	// Get items to be transferred
	for _, disk := range disks {
		for _, path := range folders {
			// logging
			logger.Blue("scanning:disk(%s):folder(%s)", disk.Path, path)

			packet := &domain.Packet{Topic: getTopic(status), Payload: fmt.Sprintf("Scanning %s on %s", path, disk.Path)}
			c.ctx.Hub.Pub(packet, "socket:broadcast")

			// check owner and permissions issues for this folder/disk combination
			packet = &domain.Packet{Topic: getTopic(status), Payload: "Checking issues ..."}
			c.ctx.Hub.Pub(packet, "socket:broadcast")

			ownIssue, grpIssue, fldIssue, filIssue, err := getIssues(reStat, disk, path)
			if err != nil {
				logger.Yellow("issues:not-available:(%s)", err)
			} else {
				ownerIssue += ownIssue
				groupIssue += grpIssue
				folderIssue += fldIssue
				fileIssue += filIssue

				logger.Blue("issues:owner(%d):group(%d):folder(%d):file(%d)", ownIssue, grpIssue, fldIssue, filIssue)
			}

			// get children files/folders to be transferred
			packet = &domain.Packet{Topic: getTopic(status), Payload: "Getting items ..."}
			c.ctx.Hub.Pub(packet, "socket:broadcast")

			list, total, err := getItems(blockSize, reItems, disk.Path, path)
			if err != nil {
				logger.Yellow("items:not-available:(%s)", err)
			} else {
				logger.Blue("items:count(%d):size(%s)", len(list), lib.ByteSize(total))
				items = append(items, list...)
			}
		}
	}

	return items, ownerIssue, groupIssue, folderIssue, fileIssue
}

func (c *Core) sendTimeFeedbackToFrontend(topic, fended string, elapsed time.Duration) {
	packet := &domain.Packet{Topic: topic, Payload: fmt.Sprintf("Ended: %s", fended)}
	c.ctx.Hub.Pub(packet, "socket:broadcast")

	packet = &domain.Packet{Topic: topic, Payload: fmt.Sprintf("Elapsed: %s", elapsed)}
	c.ctx.Hub.Pub(packet, "socket:broadcast")
}

func (c *Core) sendMailFeedback(fstarted, ffinished string, elapsed time.Duration, plan *domain.Plan, notTransferred string) {
	subject := "unbalance - PLANNING completed"
	message := fmt.Sprintf("\n\nStarted: %s\nEnded: %s\n\nElapsed: %s", fstarted, ffinished, elapsed)
	if notTransferred != "" {
		switch c.ctx.Config.NotifyPlan {
		case 1:
			message += "\n\nSome folders will not be transferred because there's not enough space for them in any of the destination disks."
		case 2:
			message += "\n\nThe following folders will not be transferred because there's not enough space for them in any of the destination disks:\n\n" + notTransferred
		}
	}

	if plan.OwnerIssue > 0 || plan.GroupIssue > 0 || plan.FolderIssue > 0 || plan.FileIssue > 0 {
		message += fmt.Sprintf(`
			\n\nThere are some permission issues:
			\n\n%d file(s)/folder(s) with an owner other than 'nobody'
			\n%d file(s)/folder(s) with a group other than 'users'
			\n%d folder(s) with a permission other than 'drwxrwxrwx'
			\n%d files(s) with a permission other than '-rw-rw-rw-' or '-r--r--r--'
			\n\nCheck the log file (/boot/logs/unbalance.log) for additional information
			\n\nIt's strongly suggested to install the Fix Common Plugins and run the Docker Safe New Permissions command
		`, plan.OwnerIssue, plan.GroupIssue, plan.FolderIssue, plan.FileIssue)
	}

	if sendErr := lib.Sendmail(common.MailCmd, c.ctx.Config.NotifyPlan, subject, message, false); sendErr != nil {
		logger.Red("unable to send mail: %s", sendErr)
	}
}

func (c *Core) getReservedAmount(size uint64) uint64 {
	var reserved uint64

	switch c.ctx.Config.ReservedUnit {
	case "%":
		fcalc := size * c.ctx.Config.ReservedAmount / 100
		reserved = fcalc
	case "Mb":
		reserved = c.ctx.Config.ReservedAmount * 1024 * 1024
	case "Gb":
		reserved = c.ctx.Config.ReservedAmount * 1024 * 1024 * 1024
	default:
		reserved = common.ReservedSpace
	}

	return reserved
}

func (c *Core) endPlan(status uint64, plan *domain.Plan, disks []*domain.Disk, items, toBeTransferred []*domain.Item) {
	plan.Ended = time.Now()
	elapsed := lib.Round(time.Since(plan.Started), time.Millisecond)
	logger.Blue("%s", elapsed) // otherwise it won't send correct value to frontend 🤷‍♂️

	fstarted := plan.Started.Format(timeFormat)
	fended := plan.Ended.Format(timeFormat)

	// Send to frontend console started/ended/elapsed times
	c.sendTimeFeedbackToFrontend(getTopic(status), fended, time.Since(plan.Started))

	// send to frontend the items that will not be transferred, if any
	// notTransferred holds a string representation of all the items, separated by a '\n'
	notTransferred := ""

	if status == common.OpScatterPlan || status == common.OpScatterPlanning {
		// some logging
		if len(toBeTransferred) == 0 {
			logger.Blue("%s:No items can be transferred.", getName(status))
		} else {
			logger.Blue("%s:%d items will be transferred.", getName(status), len(toBeTransferred))
			for _, folder := range toBeTransferred {
				logger.Blue("%s:willBeTransferred(%s)", getName(status), folder.Path)
			}
		}

		if len(items) > 0 {
			packet := &domain.Packet{Topic: getTopic(status), Payload: "The following items will not be transferred, because there's not enough space in the target disks:\n"}
			c.ctx.Hub.Pub(packet, "socket:broadcast")

			logger.Blue("%s:%d items will NOT be transferred.", getName(status), len(items))
			for _, item := range items {
				notTransferred += item.Path + "\n"

				packet = &domain.Packet{Topic: getTopic(status), Payload: item.Path}
				c.ctx.Hub.Pub(packet, "socket:broadcast")
				logger.Blue("%s:notTransferred(%s)", getName(status), item.Path)
			}
		}
	}

	// send mail according to user preferences
	c.sendMailFeedback(fstarted, fended, elapsed, plan, notTransferred)

	// some local logging
	logger.Blue("%s:ItemsLeft(%d)", getName(status), len(items))
	logger.Blue("%s:Listing (%d) disks ...", getName(status), len(disks))
	for _, disk := range disks {
		if plan.VDisks[disk.Path].Bin != nil {
			logger.Blue("=========================================================")
			logger.Blue("disk(%s):items(%d)-(%s):currentFree(%s)-plannedFree(%s)", disk.Path, len(plan.VDisks[disk.Path].Bin.Items), lib.ByteSize(plan.VDisks[disk.Path].Bin.Size), lib.ByteSize(disk.Free), lib.ByteSize(plan.VDisks[disk.Path].PlannedFree))
			logger.Blue("---------------------------------------------------------")

			for _, item := range plan.VDisks[disk.Path].Bin.Items {
				logger.Blue("[%s] %s", lib.ByteSize(item.Size), item.Name)
			}

			logger.Blue("---------------------------------------------------------")
			logger.Blue("")
		} else {
			logger.Blue("=========================================================")
			logger.Blue("disk(%s):no-items:currentFree(%s)", disk.Path, lib.ByteSize(disk.Free))
			logger.Blue("---------------------------------------------------------")
			logger.Blue("---------------------------------------------------------")
			logger.Blue("")
		}
	}

	logger.Blue("=========================================================")
	logger.Blue("Bytes To Transfer: %s", lib.ByteSize(plan.BytesToTransfer))
	logger.Blue("---------------------------------------------------------")

	packet := &domain.Packet{Topic: getTopic(status), Payload: "Planning Ended"}
	c.ctx.Hub.Pub(packet, "socket:broadcast")
}

func (c *Core) printDisks(disks []*domain.Disk, blockSize uint64) {
	logger.Blue("planner:array(%d disks):blockSize(%d)", len(disks), blockSize)
	for _, disk := range disks {
		logger.Blue("disk(%s):fs(%s):size(%d):free(%d):blocksTotal(%d):blocksFree(%d)", disk.Path, disk.FsType, disk.Size, disk.Free, disk.BlocksTotal, disk.BlocksFree)
	}
}

// HELPER FUNCTIONS
func getName(status uint64) string {
	if status == common.OpScatterPlan || status == common.OpScatterPlanning {
		return "scatterPlan"
	}

	return "gatherPlan"
}

func getTopic(status uint64) string {
	if status == common.OpScatterPlan || status == common.OpScatterPlanning {
		return common.EventScatterPlanProgress
	}

	return common.EventGatherPlanProgress
}

func removeItems(items, list []*domain.Item) []*domain.Item {
	w := 0 // write index

loop:
	for _, item := range items {
		for _, itm := range list {
			if itm.Name == item.Name {
				continue loop
			}
		}
		items[w] = item
		w++
	}

	return items[:w]
}
