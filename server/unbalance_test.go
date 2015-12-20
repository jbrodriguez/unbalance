package main

import (
	// "jbrodriguez/unbalance/server/dto"
	// "jbrodriguez/unbalance/server/helper"
	"github.com/jbrodriguez/mlog"
	"github.com/jbrodriguez/pubsub"
	"github.com/stretchr/testify/assert"
	"jbrodriguez/unbalance/server/dto"
	"jbrodriguez/unbalance/server/lib"
	"jbrodriguez/unbalance/server/model"
	"jbrodriguez/unbalance/server/services"
	"os"
	"path/filepath"
	// "regexp"
	// "strconv"
	"testing"
	"time"
)

func TestOk(t *testing.T) {
	mlog.Start(mlog.LevelInfo, "")

	disk := &model.Disk{
		Id:      1,
		Name:    "md1",
		Path:    "/mnt/disk1",
		Device:  "sdc",
		Free:    100,
		NewFree: 0,
		Size:    100,
		Serial:  "SAMSUNG_HD15",
		Status:  "DISK_OK",
	}

	folders := make([]*model.Item, 0)
	folders = append(folders,
		&model.Item{Name: "movie1", Size: 1, Path: "/mnt/disk1/Movies/movie1"},
		&model.Item{Name: "movie2", Size: 2, Path: "/mnt/disk1/Movies/movie2"},
		&model.Item{Name: "movie3", Size: 3, Path: "/mnt/disk1/Movies/movie3"},
		&model.Item{Name: "movie4", Size: 4, Path: "/mnt/disk1/Movies/movie4"},
		&model.Item{Name: "movie5", Size: 5, Path: "/mnt/disk1/Movies/movie5"},
		&model.Item{Name: "movie6", Size: 6, Path: "/mnt/disk1/Movies/movie6"},
		&model.Item{Name: "movie7", Size: 7, Path: "/mnt/disk1/Movies/movie7"},
	)

	assert.Equal(t, 7, len(folders))

	packer := lib.NewKnapsack(disk, folders, 1)
	bin := packer.BestFit()

	if assert.NotNil(t, bin) {
		bin.Print()
	}

	if assert.NotNil(t, disk) {
		disk.Print()
	}

	var size int64
	size = 28

	assert.Equal(t, size, bin.Size)

	// mlog.Stop()
}

func TestFit1(t *testing.T) {
	mlog.Start(mlog.LevelInfo, "")

	disk := &model.Disk{
		Id:      1,
		Name:    "md1",
		Path:    "/mnt/disk1",
		Device:  "sdc",
		Free:    100,
		NewFree: 0,
		Size:    100,
		Serial:  "SAMSUNG_HD15",
		Status:  "DISK_OK",
	}

	folders := make([]*model.Item, 0)
	folders = append(folders,
		&model.Item{Name: "movie1", Size: 100, Path: "/mnt/disk1/Movies/movie1"},
		&model.Item{Name: "movie2", Size: 99, Path: "/mnt/disk1/Movies/movie2"},
		&model.Item{Name: "movie3", Size: 98, Path: "/mnt/disk1/Movies/movie3"},
	)

	assert.Equal(t, 3, len(folders))

	packer := lib.NewKnapsack(disk, folders, 1)
	bin := packer.BestFit()

	if assert.NotNil(t, bin) {
		bin.Print()
	}

	if assert.NotNil(t, disk) {
		disk.Print()
	}

	var size int64
	size = 99
	assert.Equal(t, size, bin.Size)

	// mlog.Stop()
}

func TestFit2(t *testing.T) {
	mlog.Start(mlog.LevelInfo, "")

	disk := &model.Disk{
		Id:      1,
		Name:    "md1",
		Path:    "/mnt/disk1",
		Device:  "sdc",
		Free:    100,
		NewFree: 0,
		Size:    100,
		Serial:  "SAMSUNG_HD15",
		Status:  "DISK_OK",
	}

	folders := make([]*model.Item, 0)
	folders = append(folders,
		&model.Item{Name: "movie1", Size: 50, Path: "/mnt/disk1/Movies/movie1"},
		&model.Item{Name: "movie2", Size: 49, Path: "/mnt/disk1/Movies/movie2"},
		&model.Item{Name: "movie3", Size: 1, Path: "/mnt/disk1/Movies/movie3"},
	)

	assert.Equal(t, 3, len(folders))

	packer := lib.NewKnapsack(disk, folders, 1)
	bin := packer.BestFit()

	if assert.NotNil(t, bin) {
		bin.Print()
	}

	if assert.NotNil(t, disk) {
		disk.Print()
	}

	var size int64
	size = 99
	assert.Equal(t, size, bin.Size)

	// mlog.Stop()
}

// func TestFolders(t *testing.T) {
// 	re, _ := regexp.Compile(`(\d+)\s+(.*?)$`)

// 	samples := []string{
// 		"1670755474  /mnt/disk2/tv shows/./High Def",
// 		"6 /mnt/disk2/tv shows/./empty",
// 		"99297588  /mnt/disk2/tv shows/./Interstellar.mkv",
// 		"6 /mnt/disk2/tv shows/./.TemporaryFiles",
// 	}

// 	for _, sample := range samples {
// 		result := re.FindStringSubmatch(sample)
// 		// mlog.Info("[%s] %s", result[1], result[2])

// 		size, _ := strconv.ParseInt(result[1], 10, 64)

// 		name := result[2]
// 		path := filepath.Join("tv shows", filepath.Base(result[2]))

// 		mlog.Info("name(%s); path(%s); size(%d)", name, path, size)
// 	}

// 	samples2 := []string{
// 		"1670755474	/mnt/disk2/tv shows/High Def",
// 		"99297588	/mnt/disk2/tv shows/Interstellar.mkv",
// 		"6	/mnt/disk2/tv shows/empty",
// 	}

// 	for _, sample := range samples2 {
// 		result := re.FindStringSubmatch(sample)
// 		// mlog.Info("[%s] %s", result[1], result[2])

// 		size, _ := strconv.ParseInt(result[1], 10, 64)

// 		name := result[2]
// 		path := filepath.Join("tv shows", filepath.Base(result[2]))

// 		mlog.Info("name(%s); path(%s); size(%d)", name, path, size)
// 	}

// }

func createFile(home, folder, name string, size int64) error {

	os.MkdirAll(filepath.Join(home, folder), 0777)

	fd, err := os.OpenFile(filepath.Join(home, folder, name), os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		return err
	}
	defer fd.Close()

	if err := fd.Truncate(size); err != nil {
		return err
	}

	return nil
}

func TestFoldersNotMoved(t *testing.T) {
	mlog.Start(mlog.LevelInfo, "")

	home := os.Getenv("HOME")

	os.RemoveAll(filepath.Join(home, "tmp/unbalance"))

	err := createFile(home, "tmp/unbalance/mnt/disk1/movies/Interstellar (2014)", "Interstellar.mkv", 800) // 902
	assert.NoError(t, err)

	err = createFile(home, "tmp/unbalance/mnt/disk1/movies/Avatar (2009)", "Avatar.mkv", 750) // 852
	assert.NoError(t, err)

	err = createFile(home, "tmp/unbalance/mnt/disk1/movies", "Blade (1998).mkv", 450) // 450
	assert.NoError(t, err)

	folders := []string{
		"movies",
		"tvshows",
	}

	bus := pubsub.New(1)

	settings := &model.Settings{}
	settings.Folders = folders
	settings.DryRun = true
	settings.ConfigDir = filepath.Join(home, "tmp/unbalance")
	settings.LogDir = ""
	settings.Save()

	settings.ReservedSpace = 200

	condition := &model.Condition{NumDisks: 3, NumProtected: 3, Synced: time.Now(), SyncErrs: 0, Resync: 0, ResyncPrcnt: 0, ResyncPos: 0, State: "STARTED", Size: 300, Free: 108, NewFree: 0}
	disks := []*model.Disk{
		&model.Disk{Id: 1, Name: "md1", Path: filepath.Join(home, "tmp/unbalance", "mnt/disk1"), Device: "sdc", Free: 2000, NewFree: 0, Size: 2500, Serial: "SAMSUNG_HD01", Status: "DISK_OK"},
		&model.Disk{Id: 2, Name: "md2", Path: filepath.Join(home, "tmp/unbalance", "mnt/disk2"), Device: "sdd", Free: 2000, NewFree: 0, Size: 2500, Serial: "SAMSUNG_HD02", Status: "DISK_OK"},
		&model.Disk{Id: 3, Name: "md3", Path: filepath.Join(home, "tmp/unbalance", "mnt/disk3"), Device: "sde", Free: 300, NewFree: 0, Size: 2500, Serial: "SAMSUNG_HD03", Status: "DISK_OK"},
	}

	assert.Equal(t, 3, len(disks))

	unraid := &model.Unraid{}
	unraid.Condition = condition
	unraid.Disks = disks
	unraid.SourceDiskName = ""
	unraid.BytesToMove = 0

	core := services.NewCore(bus, settings)
	core.SetStorage(unraid)

	core.Start()

	destDisks := make(map[string]bool, 2)
	destDisks[filepath.Join(home, "tmp/unbalance", "mnt/disk2")] = true
	destDisks[filepath.Join(home, "tmp/unbalance", "mnt/disk3")] = true

	bestFit := &dto.BestFit{
		SourceDisk: filepath.Join(home, "tmp/unbalance", "mnt/disk1"),
		DestDisks:  destDisks,
	}

	msg := &pubsub.Message{Payload: bestFit, Reply: make(chan interface{})}
	bus.Pub(msg, "cmd.calculateBestFit")

	reply := <-msg.Reply
	resp := reply.(*model.Unraid)

	mlog.Info("Unraid: %+v", resp)

	cmd := &pubsub.Message{Reply: make(chan interface{})}
	bus.Pub(cmd, "storage:move")

	reply = <-msg.Reply

	core.Stop()
}
