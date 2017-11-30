package services

import (
	"fmt"
	"io/ioutil"
	"jbrodriguez/unbalance/server/src/algorithm"
	"jbrodriguez/unbalance/server/src/common"
	"jbrodriguez/unbalance/server/src/domain"
	"jbrodriguez/unbalance/server/src/dto"
	"jbrodriguez/unbalance/server/src/lib"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/jbrodriguez/actor"
	"github.com/jbrodriguez/mlog"
	"github.com/jbrodriguez/pubsub"
)

// Planner -
type Planner struct {
	bus      *pubsub.PubSub
	settings *lib.Settings
	actor    *actor.Actor

	reItems *regexp.Regexp
	reStat  *regexp.Regexp
}

// NewPlanner -
func NewPlanner(bus *pubsub.PubSub, settings *lib.Settings) *Planner {
	plan := &Planner{
		bus:      bus,
		settings: settings,
		actor:    actor.NewActor(bus),
	}

	plan.reItems = regexp.MustCompile(`(\d+)\s+(.*?)$`)
	plan.reStat = regexp.MustCompile(`[-dclpsbD]([-rwxsS]{3})([-rwxsS]{3})([-rwxtT]{3})\|(.*?)\:(.*?)\|(.*?)\|(.*)`)

	return plan
}

// Start -
func (p *Planner) Start() (err error) {
	mlog.Info("starting service Planner ...")

	p.actor.Register(common.IntScatterPlan, p.scatter)
	p.actor.Register(common.IntGatherPlan, p.gather)

	go p.actor.React()

	return nil
}

// Stop -
func (p *Planner) Stop() {
	mlog.Info("stopped service Planner ...")
}

func (p *Planner) scatter(msg *pubsub.Message) {
	state := msg.Payload.(*domain.State)

	mlog.Info("Running scatter planner ...")

	plan := state.Plan
	plan.Started = time.Now()

	// create two slices
	// one of source disks, the other of destinations disks
	// in scatter srcDisks will contain only one element
	srcDisk, dstDisks := getSourceAndDestinationDisks(state.Unraid.Disks, plan)

	// get dest disks with more free space to the top
	sort.Slice(dstDisks, func(i, j int) bool { return dstDisks[i].Free < dstDisks[j].Free })

	// some logging
	mlog.Info("scatterPlan:source:(%s)", srcDisk.Path)
	for _, disk := range dstDisks {
		mlog.Info("scatterPlan:dest:(%s)", disk.Path)
	}

	outbound := &dto.Packet{Topic: common.WsScatterPlanStarted, Payload: "Planning started"}
	p.bus.Pub(&pubsub.Message{Payload: outbound}, "socket:broadcast")

	items, ownerIssue, groupIssue, folderIssue, fileIssue := p.getItemsAndIssues(state.Status, p.reItems, p.reStat, []*domain.Disk{srcDisk}, plan.ChosenFolders)

	toBeTransferred := make([]*domain.Item, 0)

	// no items found, no sense going on, just end this planning
	if len(items) == 0 {
		p.endPlan(state.Status, plan, state.Unraid.Disks, items, toBeTransferred)
		p.bus.Pub(&pubsub.Message{Payload: plan}, common.IntScatterPlanFinished)
		return
	}

	plan.OwnerIssue = ownerIssue
	plan.GroupIssue = groupIssue
	plan.FolderIssue = folderIssue
	plan.FileIssue = fileIssue

	mlog.Info("scatterPlan:items(%d)", len(items))

	for _, item := range items {
		mlog.Info("scatterPlan:found(%s):size(%d)", filepath.Join(item.Location, item.Path), item.Size)

		msg := fmt.Sprintf("Found %s (%s)", filepath.Join(item.Location, item.Path), lib.ByteSize(item.Size))
		outbound = &dto.Packet{Topic: common.WsScatterPlanProgress, Payload: msg}
		p.bus.Pub(&pubsub.Message{Payload: outbound}, "socket:broadcast")
	}

	mlog.Info("scatterPlan:issues:owner(%d),group(%d),folder(%d),file(%d)", plan.OwnerIssue, plan.GroupIssue, plan.FolderIssue, plan.FileIssue)

	// Initialize fields
	plan.BytesToTransfer = 0

	for _, disk := range dstDisks {
		msg := fmt.Sprintf("Trying to allocate items to %s ...", disk.Name)
		outbound = &dto.Packet{Topic: common.WsScatterPlanProgress, Payload: msg}
		p.bus.Pub(&pubsub.Message{Payload: outbound}, "socket:broadcast")
		mlog.Info("scatterPlan:%s", msg)
		// time.Sleep(2 * time.Second)

		reserved := p.getReservedAmount(disk.Size)

		ceil := lib.Max(lib.ReservedSpace, reserved)
		mlog.Info("scatterPlan:ItemsLeft(%d):ReservedSpace(%d)", len(items), ceil)

		packer := algorithm.NewKnapsack(disk, items, ceil)
		bin := packer.BestFit()
		if bin != nil {
			plan.VDisks[disk.Path].Bin = bin
			plan.VDisks[disk.Path].PlannedFree -= bin.Size
			plan.VDisks[srcDisk.Path].PlannedFree += bin.Size

			plan.BytesToTransfer += bin.Size

			toBeTransferred = append(toBeTransferred, bin.Items...)
			items = removeItems(items, bin.Items)

			mlog.Info("scatterPlan:disk(%s):allocation=items(%d):currentFree(%s):plannedFree(%s)", disk.Path, len(bin.Items), lib.ByteSize(disk.Free), lib.ByteSize(plan.VDisks[disk.Path].PlannedFree))
		} else {
			mlog.Info("scatterPlan:disk(%s):no-allocation:currentFree(%s)", disk.Path, lib.ByteSize(disk.Free))
		}
	}

	p.endPlan(state.Status, plan, state.Unraid.Disks, items, toBeTransferred)
	p.bus.Pub(&pubsub.Message{Payload: plan}, common.IntScatterPlanFinished)
}

func (p *Planner) gather(msg *pubsub.Message) {
	state := msg.Payload.(*domain.State)

	mlog.Info("Running gather planner ...")

	plan := state.Plan
	plan.Started = time.Now()

	outbound := &dto.Packet{Topic: common.WsGatherPlanStarted, Payload: "Planning Started"}
	p.bus.Pub(&pubsub.Message{Payload: outbound}, "socket:broadcast")

	items, ownerIssue, groupIssue, folderIssue, fileIssue := p.getItemsAndIssues(state.Status, p.reItems, p.reStat, state.Unraid.Disks, plan.ChosenFolders)

	// no items found, no sense going on, just end this planning
	if len(items) == 0 {
		p.endPlan(state.Status, plan, state.Unraid.Disks, items, make([]*domain.Item, 0))
		p.bus.Pub(&pubsub.Message{Payload: plan}, common.IntScatterPlanFinished)
		return
	}

	plan.OwnerIssue = ownerIssue
	plan.GroupIssue = groupIssue
	plan.FolderIssue = folderIssue
	plan.FileIssue = fileIssue

	mlog.Info("gatherPlan:items(%d)", len(items))

	for _, item := range items {
		mlog.Info("gatherPlan:found(%s):size(%d)", filepath.Join(item.Location, item.Path), item.Size)

		msg := fmt.Sprintf("Found %s (%s)", filepath.Join(item.Location, item.Path), lib.ByteSize(item.Size))
		outbound = &dto.Packet{Topic: common.WsGatherPlanProgress, Payload: msg}
		p.bus.Pub(&pubsub.Message{Payload: outbound}, "socket:broadcast")
	}

	mlog.Info("gatherPlan:issues:owner(%d),group(%d),folder(%d),file(%d)", plan.OwnerIssue, plan.GroupIssue, plan.FolderIssue, plan.FileIssue)

	// Initialize fields
	plan.BytesToTransfer = 0

	for _, disk := range state.Unraid.Disks {
		msg := fmt.Sprintf("Trying to allocate items to %s ...", disk.Name)
		outbound = &dto.Packet{Topic: common.WsGatherPlanProgress, Payload: msg}
		p.bus.Pub(&pubsub.Message{Payload: outbound}, "socket:broadcast")
		mlog.Info("gatherPlan:%s", msg)
		// time.Sleep(2 * time.Second)

		reserved := p.getReservedAmount(disk.Size)

		ceil := lib.Max(lib.ReservedSpace, reserved)
		mlog.Info("gatherPlan:ItemsLeft(%d):ReservedSpace(%d)", len(items), ceil)

		packer := algorithm.NewGreedy(disk, items, ceil)
		bin := packer.FitAll()
		if bin != nil {
			plan.VDisks[disk.Path].Bin = bin
			plan.VDisks[disk.Path].PlannedFree -= bin.Size

			plan.BytesToTransfer += bin.Size

			mlog.Info("gatherPlan:disk(%s):allocation=items(%d):currentFree(%s):plannedFree(%s)", disk.Path, len(bin.Items), lib.ByteSize(disk.Free), lib.ByteSize(plan.VDisks[disk.Path].PlannedFree))
		} else {
			mlog.Info("gatherPlan:disk(%s):no-allocation:currentFree(%s)", disk.Path, lib.ByteSize(disk.Free))
		}
	}

	p.endPlan(state.Status, plan, state.Unraid.Disks, items, make([]*domain.Item, 0))
	p.bus.Pub(&pubsub.Message{Payload: plan}, common.IntGatherPlanFinished)
}

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

func getItems(status int64, re *regexp.Regexp, src string, folder string) ([]*domain.Item, int64, error) {
	var total int64
	srcFolder := filepath.Join(src, folder)

	var fi os.FileInfo
	var err error
	if fi, err = os.Stat(srcFolder); os.IsNotExist(err) {
		return nil, total, err
	}

	if !fi.IsDir() {
		return []*domain.Item{&domain.Item{Name: folder, Size: fi.Size(), Path: folder, Location: src}}, fi.Size(), nil
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
		total += size

		item := &domain.Item{Name: result[2], Size: size, Path: filepath.Join(folder, filepath.Base(result[2])), Location: src}
		items = append(items, item)
	})

	return items, total, err
}

func (p *Planner) getItemsAndIssues(status int64, reItems, reStat *regexp.Regexp, disks []*domain.Disk, folders []string) ([]*domain.Item, int64, int64, int64, int64) {
	var ownerIssue, groupIssue, folderIssue, fileIssue int64
	items := make([]*domain.Item, 0)

	// Get owner/permission issues
	// Get items to be transferred
	for _, disk := range disks {
		for _, path := range folders {
			// logging
			mlog.Info("scanning:disk(%s):folder(%s)", disk.Path, path)

			outbound := &dto.Packet{Topic: getTopic(status), Payload: fmt.Sprintf("Scanning %s on %s", path, disk.Path)}
			p.bus.Pub(&pubsub.Message{Payload: outbound}, "socket:broadcast")

			// check owner and permissions issues for this folder/disk combination
			outbound = &dto.Packet{Topic: getTopic(status), Payload: "Checking issues ..."}
			p.bus.Pub(&pubsub.Message{Payload: outbound}, "socket:broadcast")

			ownIssue, grpIssue, fldIssue, filIssue, err := getIssues(reStat, disk, path)
			if err != nil {
				mlog.Warning("issues:not-available:(%s)", err)
			} else {
				ownerIssue += ownIssue
				groupIssue += grpIssue
				folderIssue += fldIssue
				fileIssue += filIssue

				mlog.Info("issues:owner(%d):group(%d):folder(%d):file(%d)", ownIssue, grpIssue, fldIssue, filIssue)
			}

			// get children files/folders to be transferred
			outbound = &dto.Packet{Topic: getTopic(status), Payload: "Getting items ..."}
			p.bus.Pub(&pubsub.Message{Payload: outbound}, "socket:broadcast")

			list, total, err := getItems(status, reItems, disk.Path, path)
			if err != nil {
				mlog.Warning("items:not-available:(%s)", err)
			} else {
				mlog.Info("items:count(%d):size(%s)", len(list), lib.ByteSize(total))
				items = append(items, list...)
			}
		}
	}

	return items, ownerIssue, groupIssue, folderIssue, fileIssue
}

func (p *Planner) sendTimeFeedbackToFrontend(topic string, fstarted, ffinished string, elapsed time.Duration) {
	outbound := &dto.Packet{Topic: topic, Payload: fmt.Sprintf("Ended: %s", ffinished)}
	p.bus.Pub(&pubsub.Message{Payload: outbound}, "socket:broadcast")

	outbound = &dto.Packet{Topic: topic, Payload: fmt.Sprintf("Elapsed: %s", elapsed)}
	p.bus.Pub(&pubsub.Message{Payload: outbound}, "socket:broadcast")
}

func (p *Planner) sendMailFeeback(fstarted, ffinished string, elapsed time.Duration, plan *domain.Plan, notTransferred string) {
	subject := "unBALANCE - PLANNING completed"
	message := fmt.Sprintf("\n\nStarted: %s\nEnded: %s\n\nElapsed: %s", fstarted, ffinished, elapsed)
	if notTransferred != "" {
		switch p.settings.NotifyPlan {
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

	if sendErr := lib.Sendmail(common.MailCmd, p.settings.NotifyPlan, subject, message, false); sendErr != nil {
		mlog.Error(sendErr)
	}
}

func (p *Planner) getReservedAmount(size int64) int64 {
	var reserved int64

	switch p.settings.ReservedUnit {
	case "%":
		fcalc := size * p.settings.ReservedAmount / 100
		reserved = fcalc
	case "Mb":
		reserved = p.settings.ReservedAmount * 1000 * 1000
	case "Gb":
		reserved = p.settings.ReservedAmount * 1000 * 1000 * 1000
	default:
		reserved = lib.ReservedSpace
	}

	return reserved
}

func (p *Planner) endPlan(status int64, plan *domain.Plan, disks []*domain.Disk, items []*domain.Item, toBeTransferred []*domain.Item) {
	plan.Finished = time.Now()
	elapsed := lib.Round(time.Since(plan.Started), time.Millisecond)

	fstarted := plan.Started.Format(timeFormat)
	ffinished := plan.Finished.Format(timeFormat)

	// Send to frontend console started/ended/elapsed times
	p.sendTimeFeedbackToFrontend(getTopic(status), fstarted, ffinished, elapsed)

	// send to frontend the items that will not be transferred, if any
	// notTransferred holds a string representation of all the items, separated by a '\n'
	notTransferred := ""

	if status == common.OpScatterPlan {
		// some logging
		if len(toBeTransferred) == 0 {
			mlog.Info("%s:No items can be transferred.", getName(status))
		} else {
			mlog.Info("%s:%d items will be transferred.", getName(status), len(toBeTransferred))
			for _, folder := range toBeTransferred {
				mlog.Info("%s:willBeTransferred(%s)", getName(status), folder.Path)
			}
		}

		if len(items) > 0 {
			outbound := &dto.Packet{Topic: getTopic(status), Payload: "The following items will not be transferred, because there's not enough space in the target disks:\n"}
			p.bus.Pub(&pubsub.Message{Payload: outbound}, "socket:broadcast")

			mlog.Info("%s:%d items will NOT be transferred.", getName(status), len(items))
			for _, item := range items {
				notTransferred += item.Path + "\n"

				outbound = &dto.Packet{Topic: getTopic(status), Payload: item.Path}
				p.bus.Pub(&pubsub.Message{Payload: outbound}, "socket:broadcast")
				mlog.Info("%s:notTransferred(%s)", getName(status), item.Path)
			}
		}
	}

	// send mail according to user preferences
	p.sendMailFeeback(fstarted, ffinished, elapsed, plan, notTransferred)

	// some local logging
	mlog.Info("%s:ItemsLeft(%d)", getName(status), len(items))
	mlog.Info("%s:Listing (%d) disks ...", getName(status), len(disks))
	for _, disk := range disks {
		if plan.VDisks[disk.Path].Bin != nil {
			mlog.Info("=========================================================")
			mlog.Info("disk(%s):items(%d)-(%s):currentFree(%s)-plannedFree(%s)", disk.Path, len(plan.VDisks[disk.Path].Bin.Items), lib.ByteSize(plan.VDisks[disk.Path].Bin.Size), lib.ByteSize(disk.Free), lib.ByteSize(plan.VDisks[disk.Path].PlannedFree))
			mlog.Info("---------------------------------------------------------")

			for _, item := range plan.VDisks[disk.Path].Bin.Items {
				mlog.Info("[%s] %s", lib.ByteSize(item.Size), item.Name)
			}

			mlog.Info("---------------------------------------------------------")
			mlog.Info("")
		} else {
			mlog.Info("=========================================================")
			mlog.Info("disk(%s):no-items:currentFree(%s)", disk.Path, lib.ByteSize(disk.Free))
			mlog.Info("---------------------------------------------------------")
			mlog.Info("---------------------------------------------------------")
			mlog.Info("")
		}
	}

	mlog.Info("=========================================================")
	mlog.Info("Bytes To Transfer: %s", lib.ByteSize(plan.BytesToTransfer))
	mlog.Info("---------------------------------------------------------")

	outbound := &dto.Packet{Topic: getTopic(status), Payload: "Planning Finished"}
	p.bus.Pub(&pubsub.Message{Payload: outbound}, "socket:broadcast")
}

// HELPER FUNCTIONS
func getName(status int64) string {
	if status == common.OpScatterPlan {
		return "scatterPlan"
	}

	return "gatherPlan"
}

func getTopic(status int64) string {
	if status == common.OpScatterPlan {
		return common.WsScatterPlanProgress
	}

	return common.WsGatherPlanProgress
}

func removeItems(items []*domain.Item, list []*domain.Item) []*domain.Item {
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
