package main

import (
	// "apertoire.net/unbalance/server/dto"
	// "apertoire.net/unbalance/server/helper"
	"apertoire.net/unbalance/server/lib"
	"apertoire.net/unbalance/server/model"
	"github.com/jbrodriguez/mlog"
	"github.com/stretchr/testify/assert"
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
