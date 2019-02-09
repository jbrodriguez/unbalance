package main

import (
	// "github.com/jbrodriguez/mlog"
	// "github.com/jbrodriguez/pubsub"
	// "github.com/labstack/echo"

	// "net/http"
	// "net/http/httptest"

	// "regexp"
	// "strconv"

	// "io/ioutil"

	"fmt"
	"testing"

	"unbalance/algorithm"
	"unbalance/domain"

	"github.com/stretchr/testify/assert"
)

// var bus *pubsub.PubSub

// func TestMain(m *testing.M) {
// 	mlog.Start(mlog.LevelInfo, "")

// 	// home := os.Getenv("HOME")
// 	// path := filepath.Join(home, "tmp/mgtest")
// 	// os.RemoveAll(path)

// 	// home := os.Getenv("HOME")

// 	// // tearDown(home)
// 	// // assert.NoError(m, err)
// 	// locations := []string{
// 	// 	"/boot/config/plugins/unbalance",
// 	// 	home,
// 	// }

// 	// // settings, err := lib.NewSettings("unbalance.conf", version, locations)

// 	// // folders := []string{
// 	// // 	"/Files/Media/Videos/Movies",
// 	// // 	"/Backup",
// 	// // 	"/TVShows",
// 	// // 	"/films/blu rip",
// 	// // }

// 	// bus = pubsub.New(23)

// 	// settings, _ := lib.NewSettings("test", "vtest", locations)
// 	// settings.DryRun = false
// 	// settings.ReservedAmount = 450000000 / 1000 / 1000
// 	// settings.ReservedUnit = "Mb"
// 	// settings.APIFolders = []string{filepath.Join(home, "tmp/unbalance", "var/local/emhttp")}
// 	// settings.RsyncArgs = []string{"-avPRX"}

// 	// unraid := &domain.Unraid{NumDisks: 3, NumProtected: 3, Synced: time.Now(), SyncErrs: 0, Resync: 0, ResyncPos: 0, State: "STARTED", Size: 300, Free: 108}
// 	// disks := []*domain.Disk{
// 	// 	&domain.Disk{ID: 1, Name: "md1", Path: filepath.Join(home, "tmp/unbalance", "mnt/disk1"), Device: "sdc", Free: 1000000000000, Size: 4398046511104, Serial: "SAMSUNG_HD01", Status: "DISK_OK"},
// 	// 	&domain.Disk{ID: 2, Name: "md2", Path: filepath.Join(home, "tmp/unbalance", "mnt/disk2"), Device: "sdd", Free: 1000000000000, Size: 4398046511104, Serial: "SAMSUNG_HD02", Status: "DISK_OK"},
// 	// 	&domain.Disk{ID: 3, Name: "md3", Path: filepath.Join(home, "tmp/unbalance", "mnt/disk3"), Device: "sde", Free: 500000000000, Size: 4398046511104, Serial: "SAMSUNG_HD03", Status: "DISK_OK"},
// 	// }

// 	// // assert.Equal(m, 3, len(disks))

// 	// unraid.Disks = disks

// 	// fmt.Println("siete")
// 	// core := services.NewCore(bus, settings)
// 	// // core.SetStorage(unraid)

// 	// fmt.Println("ocho")
// 	// // mlog.Info("before start")
// 	// core.Start()
// 	// // require.Nil(m, err, "core.start error should be nil")

// 	// fmt.Println("nueve")
// 	// ret := m.Run()

// 	// // os.RemoveAll(path)

// 	// // mlog.Stop()

// 	// os.Exit(ret)
// }

// func TestOk(t *testing.T) {
// 	// mlog.Start(mlog.LevelInfo, "")

// 	disk := &domain.Disk{
// 		ID:     1,
// 		Name:   "md1",
// 		Path:   "/mnt/disk1",
// 		Device: "sdc",
// 		Free:   100,
// 		Size:   100,
// 		Serial: "SAMSUNG_HD15",
// 		Status: "DISK_OK",
// 	}

// 	folders := make([]*domain.Item, 0)
// 	folders = append(folders,
// 		&domain.Item{Name: "movie1", Size: 1, Path: "/mnt/disk1/Movies/movie1"},
// 		&domain.Item{Name: "movie2", Size: 2, Path: "/mnt/disk1/Movies/movie2"},
// 		&domain.Item{Name: "movie3", Size: 3, Path: "/mnt/disk1/Movies/movie3"},
// 		&domain.Item{Name: "movie4", Size: 4, Path: "/mnt/disk1/Movies/movie4"},
// 		&domain.Item{Name: "movie5", Size: 5, Path: "/mnt/disk1/Movies/movie5"},
// 		&domain.Item{Name: "movie6", Size: 6, Path: "/mnt/disk1/Movies/movie6"},
// 		&domain.Item{Name: "movie7", Size: 7, Path: "/mnt/disk1/Movies/movie7"},
// 	)

// 	assert.Equal(t, 7, len(folders), "there should be 7 folders")

// 	packer := algorithm.NewKnapsack(disk, folders, 1, 1)
// 	bin := packer.BestFit()

// 	// if assert.NotNil(t, bin) {
// 	// 	// bin.Print()
// 	// }

// 	// if assert.NotNil(t, disk) {
// 	// 	// disk.Print()
// 	// }

// 	var size int64 = 28

// 	assert.Equal(t, size, bin.Size, "bin size should be 28")

// 	// mlog.Stop()
// }

// func TestFit1(t *testing.T) {
// 	// mlog.Start(mlog.LevelInfo, "")

// 	disk := &domain.Disk{
// 		ID:     1,
// 		Name:   "md1",
// 		Path:   "/mnt/disk1",
// 		Device: "sdc",
// 		Free:   100,
// 		Size:   100,
// 		Serial: "SAMSUNG_HD15",
// 		Status: "DISK_OK",
// 	}

// 	folders := make([]*domain.Item, 0)
// 	folders = append(folders,
// 		&domain.Item{Name: "movie1", Size: 100, Path: "/mnt/disk1/Movies/movie1"},
// 		&domain.Item{Name: "movie3", Size: 98, Path: "/mnt/disk1/Movies/movie3"},
// 		&domain.Item{Name: "movie2", Size: 99, Path: "/mnt/disk1/Movies/movie2"},
// 	)

// 	assert.Equal(t, 3, len(folders), "there should be 3 folders")

// 	packer := algorithm.NewKnapsack(disk, folders, 1, 1)
// 	bin := packer.BestFit()

// 	// if assert.NotNil(t, bin) {
// 	// 	// bin.Print()
// 	// }

// 	// if assert.NotNil(t, disk) {
// 	// 	// disk.Print()
// 	// }

// 	var size int64 = 99

// 	assert.Equal(t, size, bin.Size, "bin.size should be 99")

// 	// mlog.Stop()
// }

// func TestFit2(t *testing.T) {
// 	// mlog.Start(mlog.LevelInfo, "")

// 	disk := &domain.Disk{
// 		ID:     1,
// 		Name:   "md1",
// 		Path:   "/mnt/disk1",
// 		Device: "sdc",
// 		Free:   100,
// 		Size:   100,
// 		Serial: "SAMSUNG_HD15",
// 		Status: "DISK_OK",
// 	}

// 	folders := make([]*domain.Item, 0)
// 	folders = append(folders,
// 		&domain.Item{Name: "movie1", Size: 50, Path: "/mnt/disk1/Movies/movie1"},
// 		&domain.Item{Name: "movie2", Size: 49, Path: "/mnt/disk1/Movies/movie2"},
// 		&domain.Item{Name: "movie3", Size: 1, Path: "/mnt/disk1/Movies/movie3"},
// 	)

// 	assert.Equal(t, 3, len(folders), "there should be 3 folders")

// 	packer := algorithm.NewKnapsack(disk, folders, 1, 1)
// 	bin := packer.BestFit()

// 	// if assert.NotNil(t, bin) {
// 	// 	// bin.Print()
// 	// }

// 	// if assert.NotNil(t, disk) {
// 	// 	// disk.Print()
// 	// }

// 	var size int64 = 99

// 	assert.Equal(t, size, bin.Size, "bin.size should be 99")

// 	// mlog.Stop()
// }

// func createFile(home, folder, name string, size int64, mode os.FileMode) error {

// 	os.MkdirAll(filepath.Join(home, folder), mode)

// 	fd, err := os.OpenFile(filepath.Join(home, folder, name), os.O_RDWR|os.O_CREATE, 0666)
// 	if err != nil {
// 		return err
// 	}
// 	defer fd.Close()

// 	return fd.Truncate(size)
// }

// func TestFoldersNotMoved(t *testing.T) {
// 	// mlog.Start(mlog.LevelInfo, "")

// 	home := os.Getenv("HOME")

// 	os.RemoveAll(filepath.Join(home, "tmp/unbalance/mnt"))

// 	err := createFile(home, "tmp/unbalance/mnt/disk1/movies/Interstellar (2014)", "Interstellar.mkv", 800) // 902
// 	assert.NoError(t, err)

// 	err = createFile(home, "tmp/unbalance/mnt/disk1/movies/Avatar (2009)", "Avatar.mkv", 750) // 852
// 	assert.NoError(t, err)

// 	err = createFile(home, "tmp/unbalance/mnt/disk1/movies", "Blade (1998).mkv", 450) // 450
// 	assert.NoError(t, err)

// 	folders := []string{
// 		"movies",
// 		"tvshows",
// 	}

// 	bus := pubsub.New(23)

// 	settings, _ := lib.NewSettings("test")
// 	settings.Folders = folders
// 	settings.DryRun = true
// 	settings.ReservedAmount = 5
// 	settings.ReservedUnit = "%"
// 	settings.ApiFolders = []string{filepath.Join(home, "tmp/unbalance", "var/local/emhttp")}

// 	condition := &model.Condition{NumDisks: 3, NumProtected: 3, Synced: time.Now(), SyncErrs: 0, Resync: 0, ResyncPos: 0, State: "STARTED", Size: 300, Free: 108, NewFree: 0}
// 	disks := []*model.Disk{
// 		&model.Disk{ID: 1, Name: "md1", Path: filepath.Join(home, "tmp/unbalance", "mnt/disk1"), Device: "sdc", Free: 1000000000000, NewFree: 0, Size: 4398046511104, Serial: "SAMSUNG_HD01", Status: "DISK_OK"},
// 		&model.Disk{ID: 2, Name: "md2", Path: filepath.Join(home, "tmp/unbalance", "mnt/disk2"), Device: "sdd", Free: 1000000000000, NewFree: 0, Size: 4398046511104, Serial: "SAMSUNG_HD02", Status: "DISK_OK"},
// 		&model.Disk{ID: 3, Name: "md3", Path: filepath.Join(home, "tmp/unbalance", "mnt/disk3"), Device: "sde", Free: 500000000000, NewFree: 0, Size: 4398046511104, Serial: "SAMSUNG_HD03", Status: "DISK_OK"},
// 	}

// 	assert.Equal(t, 3, len(disks))

// 	unraid := &model.Unraid{}
// 	unraid.Condition = condition
// 	unraid.Disks = disks
// 	unraid.SourceDiskName = ""
// 	unraid.BytesToMove = 0

// 	core := services.NewCore(bus, settings)
// 	core.SetStorage(unraid)

// 	mlog.Info("before start")
// 	err = core.Start()
// 	require.Nil(t, err, "core.start error should be nil")

// 	var packet dto.Packet
// 	// calcJSON := `{"topic":"calculate","payload":"{\"srcDisk\":\"/mnt/disk1\",\"dstDisks\":{\"/mnt/disk1\":false,\"/mnt/disk2\":true,\"/mnt/disk3\":true}}"}`
// 	template := `{"topic":"calculate","payload":"{\"srcDisk\":\"%s\",\"dstDisks\":{\"%s\":false,\"%s\":true,\"%s\":true}}"}`
// 	calcJSON := fmt.Sprintf(template,
// 		filepath.Join(home, "tmp/unbalance", "mnt/disk1"),
// 		filepath.Join(home, "tmp/unbalance", "mnt/disk1"),
// 		filepath.Join(home, "tmp/unbalance", "mnt/disk2"),
// 		filepath.Join(home, "tmp/unbalance", "mnt/disk3"),
// 	)

// 	mlog.Info("json: %s", calcJSON)
// 	err = json.NewDecoder(strings.NewReader(calcJSON)).Decode(&packet)
// 	mlog.Info("error: %s", err)
// 	mlog.Info("packet: %+v", packet)
// 	mlog.Info("payload: %s", packet.Payload)
// 	require.Nil(t, err)

// 	// destDisks := make(map[string]bool, 2)
// 	// destDisks[filepath.Join(home, "tmp/unbalance", "mnt/disk2")] = true
// 	// destDisks[filepath.Join(home, "tmp/unbalance", "mnt/disk3")] = true

// 	// args := &dto.Calculate{
// 	// 	SourceDisk: filepath.Join(home, "tmp/unbalance", "mnt/disk1"),
// 	// 	DestDisks:  destDisks,
// 	// }

// 	msg := &pubsub.Message{Payload: packet.Payload, Reply: make(chan interface{})}
// 	bus.Pub(msg, "calculate")

// 	reply := <-msg.Reply
// 	resp := reply.(*model.Unraid)

// 	mlog.Info("Unraid: %+v", resp)

// 	// cmd := &pubsub.Message{Reply: make(chan interface{})}
// 	// bus.Pub(cmd, "storage:move")

// 	// reply = <-msg.Reply

// 	core.Stop()
// }

// func tearDown(home string) error {
// 	os.RemoveAll(filepath.Join(home, "tmp/unbalance/mnt"))

// 	os.MkdirAll(filepath.Join(home, "tmp/unbalance/mnt", "disk2"), 0777)
// 	os.MkdirAll(filepath.Join(home, "tmp/unbalance/mnt", "disk3"), 0777)

// 	err := createFile(home, "tmp/unbalance/mnt/disk1/Files/Media/Videos/Movies/The Fast & Furious Series", "The Fast & The Furious.mkv", 800, 0777)
// 	if err != nil {
// 		return err
// 	}

// 	err = createFile(home, "tmp/unbalance/mnt/disk1/Files/Media/Videos/Movies/The Fast & Furious Series", "Faster & Furiousest.mkv", 1200, 0777)
// 	if err != nil {
// 		return err
// 	}

// 	err = createFile(home, "tmp/unbalance/mnt/disk1/Files/Media/Videos/Movies/", "Synchronicity [2015].mkv", 1500, 0777)
// 	if err != nil {
// 		return err
// 	}

// 	err = createFile(home, "tmp/unbalance/mnt/disk1/Backup/", "data.txt", 700, 0777)
// 	if err != nil {
// 		return err
// 	}

// 	err = createFile(home, "tmp/unbalance/mnt/disk1/TVShows/NCIS/", "NCIS 04x17 - Skeletons.avi", 2700, 0777)
// 	if err != nil {
// 		return err
// 	}

// 	err = createFile(home, "tmp/unbalance/mnt/disk1/films/blu rip/Air (2014)", "air.mkv", 1600, 0777)
// 	if err != nil {
// 		return err
// 	}

// 	err = createFile(home, "tmp/unbalance/mnt/disk1/films/blu rip/", "Interstellar.mkv", 1900, 0777)
// 	if err != nil {
// 		return err
// 	}

// 	return nil
// }

// func TestRsyncError(t *testing.T) {
// 	mlog.Info("TestRsyncError")

// 	home := os.Getenv("HOME")

// 	var packet dto.Packet
// 	// calcJSON := `{"topic":"calculate","payload":"{\"srcDisk\":\"/mnt/disk1\",\"dstDisks\":{\"/mnt/disk1\":false,\"/mnt/disk2\":true,\"/mnt/disk3\":true}}"}`
// 	template := `{"topic":"calculate","payload":"{\"srcDisk\":\"%s\",\"dstDisks\":{\"%s\":false,\"%s\":true,\"%s\":false}}"}`
// 	calcJSON := fmt.Sprintf(template,
// 		filepath.Join(home, "tmp/unbalance", "mnt/disk1"),
// 		filepath.Join(home, "tmp/unbalance", "mnt/disk1"),
// 		filepath.Join(home, "tmp/unbalance", "mnt/disk2"),
// 		filepath.Join(home, "tmp/unbalance", "mnt/disk3"),
// 	)

// 	// errChown := os.Chown(filepath.Join(home, "tmp/unbalance", "/mnt/disk1/films/blu rip/Air (2014)/air.mkv"), 0, 0)
// 	// if errChown != nil {
// 	// 	mlog.Warning("error chowning: %s", errChown)
// 	// 	return
// 	// }

// 	// errChown := os.Chown(filepath.Join(home, "tmp/unbalance", "/mnt/disk1/TVShows/NCIS/NCIS 04x17 - Skeletons.avi"), 0, 0)
// 	// if errChown != nil {
// 	// 	mlog.Warning("error chowning: %s", errChown)
// 	// 	return
// 	// }

// 	mlog.Info("json: %s", calcJSON)
// 	err := json.NewDecoder(strings.NewReader(calcJSON)).Decode(&packet)
// 	mlog.Info("error: %s", err)
// 	mlog.Info("packet: %+v", packet)
// 	mlog.Info("payload: %s", packet.Payload)
// 	require.Nil(t, err)

// 	// destDisks := make(map[string]bool, 2)
// 	// destDisks[filepath.Join(home, "tmp/unbalance", "mnt/disk2")] = true
// 	// destDisks[filepath.Join(home, "tmp/unbalance", "mnt/disk3")] = true

// 	// args := &dto.Calculate{
// 	// 	SourceDisk: filepath.Join(home, "tmp/unbalance", "mnt/disk1"),
// 	// 	DestDisks:  destDisks,
// 	// }

// 	msg := &pubsub.Message{Payload: packet.Payload, Reply: make(chan interface{})}
// 	bus.Pub(msg, "calculate")

// 	usr, _ := user.Current()
// 	mlog.Info("user %+v", usr)

// 	time.Sleep(5 * time.Second)

// 	errRemove := os.RemoveAll(filepath.Join(home, "tmp/unbalance", "/mnt/disk1/TVShows"))
// 	if errRemove != nil {
// 		mlog.Warning("error removing: %s", errRemove)
// 	}

// 	// mlog.Info("Unraid (after calc): %+v", resp)

// 	cmd := &pubsub.Message{Reply: make(chan interface{})}
// 	bus.Pub(cmd, "move")

// 	time.Sleep(10 * time.Second)

// 	// err = tearDown(home)
// 	// assert.NoError(t, err)

// 	// mlog.Info("TestRsyncCustom")

// 	// payload := `{"rsyncFlags":["-avX", "--partial"]}`
// 	// msg = &pubsub.Message{Payload: payload, Reply: make(chan interface{})}
// 	// bus.Pub(msg, "/config/set/rsyncFlags")

// 	// mlog.Info("flags set")

// 	// time.Sleep(5 * time.Second)

// 	// mlog.Info("packet %+v", packet.Payload)

// 	// calc := &pubsub.Message{Payload: packet.Payload}
// 	// bus.Pub(calc, "calculate")

// 	// mlog.Info("calculate 2 sent")

// 	// // mlog.Info("Unraid (after calc): %+v", resp)

// 	// cmd2 := &pubsub.Message{Reply: make(chan interface{})}
// 	// bus.Pub(cmd2, "move")

// 	// time.Sleep(10 * time.Second)

// 	// core.Stop()
// }

// func TestRsync(t *testing.T) {
// 	mlog.Info("TestRsyncDefault")

// 	home := os.Getenv("HOME")

// 	var packet dto.Packet
// 	// calcJSON := `{"topic":"calculate","payload":"{\"srcDisk\":\"/mnt/disk1\",\"dstDisks\":{\"/mnt/disk1\":false,\"/mnt/disk2\":true,\"/mnt/disk3\":true}}"}`
// 	template := `{"topic":"calculate","payload":"{\"srcDisk\":\"%s\",\"dstDisks\":{\"%s\":false,\"%s\":true,\"%s\":false}}"}`
// 	calcJSON := fmt.Sprintf(template,
// 		filepath.Join(home, "tmp/unbalance", "mnt/disk1"),
// 		filepath.Join(home, "tmp/unbalance", "mnt/disk1"),
// 		filepath.Join(home, "tmp/unbalance", "mnt/disk2"),
// 		filepath.Join(home, "tmp/unbalance", "mnt/disk3"),
// 	)

// 	mlog.Info("json: %s", calcJSON)
// 	err := json.NewDecoder(strings.NewReader(calcJSON)).Decode(&packet)
// 	mlog.Info("error: %s", err)
// 	mlog.Info("packet: %+v", packet)
// 	mlog.Info("payload: %s", packet.Payload)
// 	require.Nil(t, err)

// 	// destDisks := make(map[string]bool, 2)
// 	// destDisks[filepath.Join(home, "tmp/unbalance", "mnt/disk2")] = true
// 	// destDisks[filepath.Join(home, "tmp/unbalance", "mnt/disk3")] = true

// 	// args := &dto.Calculate{
// 	// 	SourceDisk: filepath.Join(home, "tmp/unbalance", "mnt/disk1"),
// 	// 	DestDisks:  destDisks,
// 	// }

// 	msg := &pubsub.Message{Payload: packet.Payload, Reply: make(chan interface{})}
// 	bus.Pub(msg, "calculate")

// 	time.Sleep(5 * time.Second)

// 	// mlog.Info("Unraid (after calc): %+v", resp)

// 	cmd := &pubsub.Message{Reply: make(chan interface{})}
// 	bus.Pub(cmd, "move")

// 	time.Sleep(10 * time.Second)

// 	// err = tearDown(home)
// 	// assert.NoError(t, err)

// 	// mlog.Info("TestRsyncCustom")

// 	// payload := `{"rsyncFlags":["-avX", "--partial"]}`
// 	// msg = &pubsub.Message{Payload: payload, Reply: make(chan interface{})}
// 	// bus.Pub(msg, "/config/set/rsyncFlags")

// 	// mlog.Info("flags set")

// 	// time.Sleep(5 * time.Second)

// 	// mlog.Info("packet %+v", packet.Payload)

// 	// calc := &pubsub.Message{Payload: packet.Payload}
// 	// bus.Pub(calc, "calculate")

// 	// mlog.Info("calculate 2 sent")

// 	// // mlog.Info("Unraid (after calc): %+v", resp)

// 	// cmd2 := &pubsub.Message{Reply: make(chan interface{})}
// 	// bus.Pub(cmd2, "move")

// 	// time.Sleep(10 * time.Second)

// 	// core.Stop()
// }

// func TestRsyncCustom(t *testing.T) {
// 	mlog.Info("TestRsyncCustom")

// 	home := os.Getenv("HOME")

// 	payload := `{"rsyncFlags":["-avX", "--partial"]}`
// 	msg := &pubsub.Message{Payload: payload, Reply: make(chan interface{})}
// 	bus.Pub(msg, "/config/set/rsyncFlags")

// 	mlog.Info("flags set")

// 	var packet dto.Packet
// 	// calcJSON := `{"topic":"calculate","payload":"{\"srcDisk\":\"/mnt/disk1\",\"dstDisks\":{\"/mnt/disk1\":false,\"/mnt/disk2\":true,\"/mnt/disk3\":true}}"}`
// 	template := `{"topic":"calculate","payload":"{\"srcDisk\":\"%s\",\"dstDisks\":{\"%s\":false,\"%s\":true,\"%s\":false}}"}`
// 	calcJSON := fmt.Sprintf(template,
// 		filepath.Join(home, "tmp/unbalance", "mnt/disk1"),
// 		filepath.Join(home, "tmp/unbalance", "mnt/disk1"),
// 		filepath.Join(home, "tmp/unbalance", "mnt/disk2"),
// 		filepath.Join(home, "tmp/unbalance", "mnt/disk3"),
// 	)

// 	mlog.Info("json: %s", calcJSON)
// 	err := json.NewDecoder(strings.NewReader(calcJSON)).Decode(&packet)
// 	mlog.Info("error: %s", err)
// 	mlog.Info("packet: %+v", packet)
// 	mlog.Info("payload: %s", packet.Payload)
// 	require.Nil(t, err)

// 	// destDisks := make(map[string]bool, 2)
// 	// destDisks[filepath.Join(home, "tmp/unbalance", "mnt/disk2")] = true
// 	// destDisks[filepath.Join(home, "tmp/unbalance", "mnt/disk3")] = true

// 	// args := &dto.Calculate{
// 	// 	SourceDisk: filepath.Join(home, "tmp/unbalance", "mnt/disk1"),
// 	// 	DestDisks:  destDisks,
// 	// }

// 	msg = &pubsub.Message{Payload: packet.Payload, Reply: make(chan interface{})}
// 	bus.Pub(msg, "calculate")

// 	time.Sleep(5 * time.Second)

// 	// mlog.Info("Unraid (after calc): %+v", resp)

// 	cmd := &pubsub.Message{Reply: make(chan interface{})}
// 	bus.Pub(cmd, "move")

// 	// time.Sleep(10 * time.Second)

// 	// err = tearDown(home)
// 	// assert.NoError(t, err)

// 	// mlog.Info("TestRsyncCustom")

// 	// payload := `{"rsyncFlags":["-avX", "--partial"]}`
// 	// msg = &pubsub.Message{Payload: payload, Reply: make(chan interface{})}
// 	// bus.Pub(msg, "/config/set/rsyncFlags")

// 	// mlog.Info("flags set")

// 	// time.Sleep(5 * time.Second)

// 	// mlog.Info("packet %+v", packet.Payload)

// 	// calc := &pubsub.Message{Payload: packet.Payload}
// 	// bus.Pub(calc, "calculate")

// 	// mlog.Info("calculate 2 sent")

// 	// // mlog.Info("Unraid (after calc): %+v", resp)

// 	// cmd2 := &pubsub.Message{Reply: make(chan interface{})}
// 	// bus.Pub(cmd2, "move")

// 	// time.Sleep(10 * time.Second)

// 	// core.Stop()
// }

// func TestBind(t *testing.T) {
// 	// userJSON := `{"topic":"calculate","payload":"{\"srcDisk\":\"/mnt/disk1\",\"dstDisks\":{\"/mnt/disk1\":false,\"/mnt/disk2\":true,\"/mnt/disk3\":true}}"}`
// 	userJSON := `{"srcDisk":"/mnt/disk1","dstDisks":{"/mnt/disk1":false,"/mnt/disk2":true,"/mnt/disk3":true}}`
//
// 	e := echo.New()
//
// 	req, _ := http.NewRequest(echo.POST, "/", strings.NewReader(userJSON))
// 	rec := httptest.NewRecorder()
// 	c := echo.NewContext(req, echo.NewResponse(rec, e), e)
//
// 	testBind(t, c, "application/json")
// }

// {\"srcDisk\":\"/mnt/disk1\",\"dstDisks\":{\"/mnt/disk1\":false,\"/mnt/disk2\":true,\"/mnt/disk3\":true}}
// {\"srcDisk\":\"/mnt/disk1\",\"dstDisks\":{\"/mnt/disk1\":false,\"/mnt/disk2\":true,\"/mnt/disk3\":true}}

// func testBind(t *testing.T, c *echo.Context, ct string) {
// 	c.Request().Header.Set(echo.ContentType, ct)
// 	var args dto.Calculate
// 	err := c.Bind(&args)
// 	if ct == "" {
// 		assert.Error(t, echo.UnsupportedMediaType)
// 	} else if assert.NoError(t, err) {
// 		assert.Equal(t, "/mnt/disk1", args.SourceDisk)
// 		assert.Equal(t, `{"/mnt/disk1":false,"/mnt/disk2":true,"/mnt/disk3":true}`, args.DestDisks)
//
// 		// assert.Equal(t, "calculate", args.Topic)
// 		// assert.Equal(t, `{"srcDisk":"/mnt/disk1","dstDisks":{"/mnt/disk1":false,"/mnt/disk2":true,"/mnt/disk3":true}}`, args.Payload)
//
// 		// var param dto.Calculate
// 		// err := c.Bind(&param)
// 		// if assert.NoError(t, err) {
// 		// 	assert.Equal(t, "/mnt/disk1", param.SourceDisk)
// 		// 	assert.Equal(t, `{"/mnt/disk1":false,"/mnt/disk2":true,"/mnt/disk3":true}`, param.DestDisks)
// 		// }
// 	}
// }

// func TestPercentProgress(t *testing.T) {
// 	started := time.Now()

// 	var bytesToMove int64 = 1299623666930

// 	var bytesMoved int64 = 19515085445
// 	delta := time.Since(started) + (time.Minute * 23)

// 	speed := float64(bytesMoved) / delta.Seconds()
// 	mbs := speed / 1024 / 1024

// 	left := float64(bytesToMove-bytesMoved) / speed
// 	duration := time.Duration(left) * time.Second

// 	mlog.Info("left(%s) | mbs(%.2f MB/s) | (delta=%d)", duration, mbs, delta)
// }

// func TestCommandCreation(t *testing.T) {
// 	rsyncArgs := []string{
// 		"-avX",
// 		"--partial",
// 	}

// 	diskName := "/mnt/disk2"
// 	itemPath := "TVShows/NCIS/NCIS 04x17 - Skeletons.avi"
// 	diskPath := "/mnt/disk3"

// 	cmd := fmt.Sprintf("rsync %s \"%s\" \"%s/\"", strings.Join(rsyncArgs, " "), filepath.Join(diskName, itemPath), filepath.Join(diskPath, filepath.Dir(itemPath)))
// 	mlog.Info("cmd(%s)", cmd)
// 	assert.Equal(t, `rsync -avX --partial "/mnt/disk2/TVShows/NCIS/NCIS 04x17 - Skeletons.avi" "/mnt/disk3/TVShows/NCIS/"`, cmd)

// 	diskName = "/mnt/disk3"
// 	itemPath = "blurip/Air (2014)"
// 	diskPath = "/mnt/disk2"

// 	cmd = fmt.Sprintf("rsync %s \"%s\" \"%s/\"", strings.Join(rsyncArgs, " "), filepath.Join(diskName, itemPath), filepath.Join(diskPath, filepath.Dir(itemPath)))
// 	mlog.Info("cmd(%s)", cmd)
// 	assert.Equal(t, `rsync -avX --partial "/mnt/disk3/blurip/Air (2014)" "/mnt/disk2/blurip/"`, cmd)
// }

// func TestComparisons(t *testing.T) {
// 	perms := "rw-rw-rw-"
// 	match := strings.Compare(perms, "r--r--r--") == 0 || strings.Compare(perms, "rw-rw-rw-") == 0
// 	assert.Equal(t, true, match, fmt.Sprintf("perms(%s) should be fine", perms))

// 	perms = "r--r--r--"
// 	match = strings.Compare(perms, "r--r--r--") == 0 || strings.Compare(perms, "rw-rw-rw-") == 0
// 	assert.Equal(t, true, match, fmt.Sprintf("perms(%s) should be fine", perms))

// 	perms = "r-xr--r--"
// 	match = strings.Compare(perms, "r--r--r--") == 0 || strings.Compare(perms, "rw-rw-rw-") == 0
// 	assert.Equal(t, false, match, fmt.Sprintf("perms(%s) should be fine", perms))

// 	perms = "rw-rw-rwx"
// 	match = strings.Compare(perms, "r--r--r--") == 0 || strings.Compare(perms, "rw-rw-rw-") == 0
// 	assert.Equal(t, false, match, fmt.Sprintf("perms(%s) should be fine", perms))
// }

// func TestPermissionRegex(t *testing.T) {
// 	re, _ := regexp.Compile(`[-dclpsbD]([-rwxsS]{3})([-rwxsS]{3})([-rwxtT]{3})\|(.*?)\:(.*?)\|(.*?)\|(.*)`)

// 	line1 := "drwxrwxrwx|nobody:users|directory|/mnt/disk1/Media/."
// 	line2 := "drwxrwsrwx|nobody:users|directory|/mnt/disk1/Media/./Movies"
// 	line3 := "drwxrwsrwx|nobody:users|directory|/mnt/disk1/Media/./Movies/A Hologram for the King (2016)"

// 	result := re.FindStringSubmatch(line1)
// 	assert.NotNil(t, result, fmt.Sprintf("line(%s) couldn't be parsed", line1))

// 	result = re.FindStringSubmatch(line2)
// 	assert.NotNil(t, result, fmt.Sprintf("line(%s) couldn't be parsed", line2))

// 	result = re.FindStringSubmatch(line3)
// 	assert.NotNil(t, result, fmt.Sprintf("line(%s) couldn't be parsed", line3))
// }

// func TestStat(t *testing.T) {
// 	home := os.Getenv("HOME")

// 	folder := filepath.Join(home, "tmp/unbalance/mnt")

// 	os.RemoveAll(folder)

// 	err := createFile(home, "tmp/unbalance/mnt/disk1/movies/Interstellar (2014)", "Interstellar.mkv", 800, 0222) // 902
// 	assert.NoError(t, err)

// 	dirs, err := ioutil.ReadDir(folder)
// 	if err != nil {
// 		mlog.Warning("getFolders:Unable to readdir: %s", err)
// 	}

// 	if len(dirs) == 0 {
// 		mlog.Info("getFolders:No subdirectories under %s", folder)
// 	}
// }

func TestKnap(t *testing.T) {
	// mlog.Start(mlog.LevelInfo, "")
	// mlog.Info("first blood")
	fmt.Printf("first blood\n")

	disk := &domain.Disk{
		ID:     1,
		Name:   "md1",
		Path:   "/mnt/disk1",
		Device: "sdc",
		Free:   195170000,
		Size:   45015000000,
		Serial: "SAMSUNG_HD15",
		Status: "DISK_OK",
	}

	folders := make([]*domain.Item, 0)
	folders = append(folders,
		&domain.Item{Name: "k19", Size: 40059033087, Path: "/mnt/disk1/Movies/k19"},
		&domain.Item{Name: "almostf", Size: 36752144890, Path: "/mnt/disk1/Movies/almostf"},
		&domain.Item{Name: "asgood", Size: 36570146300, Path: "/mnt/disk1/Movies/asgood"},
	)

	// assert.Equal(t, 3, len(folders), "there should be 3 folders")

	packer := algorithm.NewKnapsack(disk, folders, 512*1024*1024, 0)
	bin := packer.BestFit()

	assert.Nil(t, bin)

	// mlog.Stop()
}
