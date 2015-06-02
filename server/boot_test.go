package main

import (
	// "apertoire.net/unbalance/server/dto"
	// "apertoire.net/unbalance/server/helper"
	"apertoire.net/unbalance/server/lib"
	"apertoire.net/unbalance/server/model"
	// "apertoire.net/unbalance/server/services"
	"github.com/jbrodriguez/mlog"
	// "github.com/jbrodriguez/pubsub"
	"github.com/stretchr/testify/assert"
	// "os"
	"path/filepath"
	"regexp"
	"strconv"
	"testing"
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

func TestFolders(t *testing.T) {
	re, _ := regexp.Compile(`(\d+)\s+(.*?)$`)

	samples := []string{
		"1670755474  /mnt/disk2/tv shows/./High Def",
		"6 /mnt/disk2/tv shows/./empty",
		"99297588  /mnt/disk2/tv shows/./Interstellar.mkv",
		"6 /mnt/disk2/tv shows/./.TemporaryFiles",
	}

	for _, sample := range samples {
		result := re.FindStringSubmatch(sample)
		// mlog.Info("[%s] %s", result[1], result[2])

		size, _ := strconv.ParseInt(result[1], 10, 64)

		name := result[2]
		path := filepath.Join("tv shows", filepath.Base(result[2]))

		mlog.Info("name(%s); path(%s); size(%d)", name, path, size)
	}

	samples2 := []string{
		"1670755474	/mnt/disk2/tv shows/High Def",
		"99297588	/mnt/disk2/tv shows/Interstellar.mkv",
		"6	/mnt/disk2/tv shows/empty",
	}

	for _, sample := range samples2 {
		result := re.FindStringSubmatch(sample)
		// mlog.Info("[%s] %s", result[1], result[2])

		size, _ := strconv.ParseInt(result[1], 10, 64)

		name := result[2]
		path := filepath.Join("tv shows", filepath.Base(result[2]))

		mlog.Info("name(%s); path(%s); size(%d)", name, path, size)
	}

}

// func TestGetFolders(t *testing.T) {
// 	home := os.Getenv("HOME")

// 	mlog.Start(mlog.LevelInfo, "")

// 	bus := pubsub.New(1)
// 	core := services.NewCore(bus, nil)

// 	locations := []string{
// 		filepath.Join(home, "tmp/unbalance/.mediagui"),
// 		filepath.Join(home, "tmp/unbalance/empty"),
// 	}

// 	for _, location := range locations {
// 		os.MkdirAll(location, 0777)
// 	}
// 	// defer os.RemoveAll(filepath.Join(home, "tmp/unbalance"))

// 	items := core.TestGetFolders(home, "tmp/unbalance")

// 	assert.Equal(t, 2, len(items))

// 	for _, item := range items {
// 		mlog.Info("name(%s); path(%s); size(%d)", item.Name, item.Path, item.Size)

// 	}
// 	mlog.Info("items: %+v", items)

// }
