package main

import (
	"apertoire.net/unbalance/lib"
	"apertoire.net/unbalance/model"
	"fmt"
	"log"
	"sort"
	"testing"
)

func TestKnapsack(t *testing.T) {
	// disks := []*model.Disk{
	// 	&model.Disk{Id:1 Name:"md1", Path:"/mnt/disk1", Device:"sdn", Free:2798870528, FreeMinusBuffer:0, NewFree:2798870528, Size:3000501350400, Serial:"WDC_WD30EZRX-00DC0B0_WD-WMC1T0345089", Status:"DISK_OK"},
	// 	&model.Disk{Id:2 Name:"md2", Path:"/mnt/disk2", Device:"sdm", Free:12431601664, FreeMinusBuffer:0, NewFree:12431601664, Size:3000501350400, Serial:"WDC_WD30EZRX-00DC0B0_WD-WMC1T0373550", Status:"DISK_OK"}
	// 	})

	disks := []*model.Disk{
		(&model.Disk{Id: 1, Name: "md1", Path: "/mnt/disk1", Device: "sdn", Free: 2798870528, FreeMinusBuffer: 0, NewFree: 2798870528, Size: 3000501350400, Serial: "WDC_WD30EZRX-00DC0B0_WD-WMC1T0345089", Status: "DISK_OK"}),
		(&model.Disk{Id: 2, Name: "md2", Path: "/mnt/disk2", Device: "sdm", Free: 12431601664, FreeMinusBuffer: 0, NewFree: 12431601664, Size: 3000501350400, Serial: "WDC_WD30EZRX-00DC0B0_WD-WMC1T0373550", Status: "DISK_OK"}),
		(&model.Disk{Id: 3, Name: "md3", Path: "/mnt/disk3", Device: "sdk", Free: 8654426112, FreeMinusBuffer: 0, NewFree: 8654426112, Size: 3000501350400, Serial: "ST3000DM001-9YN166_W1F181AR", Status: "DISK_OK"}),
		(&model.Disk{Id: 4, Name: "md4", Path: "/mnt/disk4", Device: "sdl", Free: 110264877056, FreeMinusBuffer: 0, NewFree: 110264877056, Size: 3000501350400, Serial: "ST3000DM001-9YN166_Z1F1546H", Status: "DISK_OK"}),
		(&model.Disk{Id: 5, Name: "md5", Path: "/mnt/disk5", Device: "sdi", Free: 7675904, FreeMinusBuffer: 0, NewFree: 7675904, Size: 3000501350400, Serial: "TOSHIBA_DT01ACA300_23CEUGZWS", Status: "DISK_OK"}),
		(&model.Disk{Id: 6, Name: "md6", Path: "/mnt/disk6", Device: "sdj", Free: 13362188288, FreeMinusBuffer: 0, NewFree: 13362188288, Size: 3000501350400, Serial: "TOSHIBA_DT01ACA300_23CENSPWS", Status: "DISK_OK"}),
		(&model.Disk{Id: 7, Name: "md7", Path: "/mnt/disk7", Device: "sdh", Free: 10317832192, FreeMinusBuffer: 0, NewFree: 10317832192, Size: 3000501350400, Serial: "TOSHIBA_DT01ACA300_23DG6Z7WS", Status: "DISK_OK"}),
		(&model.Disk{Id: 8, Name: "md8", Path: "/mnt/disk8", Device: "sdb", Free: 116319207424, FreeMinusBuffer: 0, NewFree: 116319207424, Size: 3000501350400, Serial: "ST3000DM001-1CH166_W1F45LE8", Status: "DISK_OK"}),
		(&model.Disk{Id: 9, Name: "md9", Path: "/mnt/disk9", Device: "sdg", Free: 25462644736, FreeMinusBuffer: 0, NewFree: 25462644736, Size: 3000501350400, Serial: "TOSHIBA_DT01ACA300_Y3UEB7GGS", Status: "DISK_OK"}),
		(&model.Disk{Id: 10, Name: "md10", Path: "/mnt/disk10", Device: "sdf", Free: 380406677504, FreeMinusBuffer: 0, NewFree: 380406677504, Size: 3000501350400, Serial: "TOSHIBA_DT01ACA300_X3V9V7TGS", Status: "DISK_OK"}),
		(&model.Disk{Id: 11, Name: "md11", Path: "/mnt/disk11", Device: "sde", Free: 0, FreeMinusBuffer: 0, NewFree: 0, Size: 3000501350400, Serial: "WDC_WD30EFRX-68AX9N0_WD-WMC1T0571629", Status: "DISK_OK"}),
		(&model.Disk{Id: 12, Name: "md12", Path: "/mnt/disk12", Device: "sdd", Free: 6960766976, FreeMinusBuffer: 0, NewFree: 6960766976, Size: 4000664875008, Serial: "ST4000DM000-1F2168_Z301LVKC", Status: "DISK_OK"}),
		(&model.Disk{Id: 13, Name: "md13", Path: "/mnt/disk13", Device: "sdc", Free: 67401682944, FreeMinusBuffer: 0, NewFree: 67401682944, Size: 4000664875008, Serial: "WDC_WD40EZRX-00SPEB0_WD-WCC4EM0WN2RE", Status: "DISK_OK"}),
	}

	folders := []*model.Item{
		&model.Item{Name: "/mnt/disk10/films/bluray/12 Years A Slave (2013)", Size: 47203411422, Path: "films/bluray/12 Years A Slave (2013)"},
		&model.Item{Name: "/mnt/disk10/films/bluray/47 Ronin (2013)", Size: 45129611679, Path: "films/bluray/47 Ronin (2013)"},
		&model.Item{Name: "/mnt/disk10/films/bluray/A Clockwork Orange (1971)", Size: 37421608444, Path: "films/bluray/A Clockwork Orange (1971)"},
		&model.Item{Name: "/mnt/disk10/films/bluray/Adventures In Babysitting (1987)", Size: 24674732196, Path: "films/bluray/Adventures In Babysitting (1987)"},
		&model.Item{Name: "/mnt/disk10/films/bluray/Air America (1990)", Size: 22267598888, Path: "films/bluray/Air America (1990)"},
		&model.Item{Name: "/mnt/disk10/films/bluray/Anchorman 2 The Legend Continues (2013)", Size: 48809853015, Path: "films/bluray/Anchorman 2 The Legend Continues (2013)"},
		&model.Item{Name: "/mnt/disk10/films/bluray/Atlantis The Lost Empire (2001)", Size: 29041292089, Path: "films/bluray/Atlantis The Lost Empire (2001)"},
		&model.Item{Name: "/mnt/disk10/films/bluray/Bandidas (2006)", Size: 24494396452, Path: "films/bluray/Bandidas (2006)"},
		&model.Item{Name: "/mnt/disk10/films/bluray/Basic Instinct (1992)", Size: 30436608336, Path: "films/bluray/Basic Instinct (1992)"},
		&model.Item{Name: "/mnt/disk10/films/bluray/Battle of the Year (2013)", Size: 35826748569, Path: "films/bluray/Battle of the Year (2013)"},
		&model.Item{Name: "/mnt/disk10/films/bluray/Beyond Outrage (2012)", Size: 23621005231, Path: "films/bluray/Beyond Outrage (2012)"},
		&model.Item{Name: "/mnt/disk10/films/bluray/Bull Durham (1988)", Size: 21270313400, Path: "films/bluray/Bull Durham (1988)"},
		&model.Item{Name: "/mnt/disk10/films/bluray/Cannibal Holocaust (1980)", Size: 24136941139, Path: "films/bluray/Cannibal Holocaust (1980)"},
		&model.Item{Name: "/mnt/disk10/films/bluray/Clear and Present Danger (1994)", Size: 48885583465, Path: "films/bluray/Clear and Present Danger (1994)"},
		&model.Item{Name: "/mnt/disk10/films/bluray/Conspiracy Theory (1997)", Size: 38441334783, Path: "films/bluray/Conspiracy Theory (1997)"},
		&model.Item{Name: "/mnt/disk10/films/bluray/Delivery Man (2013)", Size: 35777554713, Path: "films/bluray/Delivery Man (2013)"},
		&model.Item{Name: "/mnt/disk10/films/bluray/Die Hard 2 (1990)", Size: 41159696026, Path: "films/bluray/Die Hard 2 (1990)"},
		&model.Item{Name: "/mnt/disk10/films/bluray/Dom Hemingway (2013)", Size: 27498101849, Path: "films/bluray/Dom Hemingway (2013)"},
		&model.Item{Name: "/mnt/disk10/films/bluray/Dumbo (1941)", Size: 38924355171, Path: "films/bluray/Dumbo (1941)"},
		&model.Item{Name: "/mnt/disk10/films/bluray/Edison (2005)", Size: 22278484492, Path: "films/bluray/Edison (2005)"},
		&model.Item{Name: "/mnt/disk10/films/bluray/Fantastic Mr. Fox (2009)", Size: 47178630394, Path: "films/bluray/Fantastic Mr. Fox (2009)"},
		&model.Item{Name: "/mnt/disk10/films/bluray/Ghostbusters (1984)", Size: 34459588025, Path: "films/bluray/Ghostbusters (1984)"},
		&model.Item{Name: "/mnt/disk10/films/bluray/Grudge Match (2013)", Size: 34255545678, Path: "films/bluray/Grudge Match (2013)"},
		&model.Item{Name: "/mnt/disk10/films/bluray/Harry And The Hendersons (1987)", Size: 24793805056, Path: "films/bluray/Harry And The Hendersons (1987)"},
		&model.Item{Name: "/mnt/disk10/films/bluray/Home of the Brave (2006)", Size: 24842064623, Path: "films/bluray/Home of the Brave (2006)"},
		&model.Item{Name: "/mnt/disk10/films/bluray/In The Name Of The Father (1993)", Size: 41672138942, Path: "films/bluray/In The Name Of The Father (1993)"},
		&model.Item{Name: "/mnt/disk10/films/bluray/Inside Llewyn Davis (2013)", Size: 23434092035, Path: "films/bluray/Inside Llewyn Davis (2013)"},
		&model.Item{Name: "/mnt/disk10/films/bluray/Johnny English Reborn (2011)", Size: 43870275979, Path: "films/bluray/Johnny English Reborn (2011)"},
		&model.Item{Name: "/mnt/disk10/films/bluray/Julie & Julia (2009)", Size: 48432197133, Path: "films/bluray/Julie & Julia (2009)"},
		&model.Item{Name: "/mnt/disk10/films/bluray/Kiss of the Spider Woman (1985)", Size: 44837864869, Path: "films/bluray/Kiss of the Spider Woman (1985)"},
		&model.Item{Name: "/mnt/disk10/films/bluray/Le Petit Nicolas (2009)", Size: 45820947171, Path: "films/bluray/Le Petit Nicolas (2009)"},
		&model.Item{Name: "/mnt/disk10/films/bluray/Led Zeppelin The Song Remains The Same (1976)", Size: 31167960530, Path: "films/bluray/Led Zeppelin The Song Remains The Same (1976)"},
		&model.Item{Name: "/mnt/disk10/films/bluray/Les Femmes de L'Ombre (2008)", Size: 34799253865, Path: "films/bluray/Les Femmes de L'Ombre (2008)"},
		&model.Item{Name: "/mnt/disk10/films/bluray/Love Camp (1977)", Size: 26489110149, Path: "films/bluray/Love Camp (1977)"},
		&model.Item{Name: "/mnt/disk10/films/bluray/Max Payne (2008)", Size: 34512237271, Path: "films/bluray/Max Payne (2008)"},
		&model.Item{Name: "/mnt/disk10/films/bluray/On Her Majesty's Secret Service (1969)", Size: 45208614548, Path: "films/bluray/On Her Majesty's Secret Service (1969)"},
		&model.Item{Name: "/mnt/disk10/films/bluray/Philomena (2013)", Size: 45561597750, Path: "films/bluray/Philomena (2013)"},
		&model.Item{Name: "/mnt/disk10/films/bluray/Pushing Hands (1992)", Size: 22658171885, Path: "films/bluray/Pushing Hands (1992)"},
		&model.Item{Name: "/mnt/disk10/films/bluray/Ride Along (2014)", Size: 39683893054, Path: "films/bluray/Ride Along (2014)"},
		&model.Item{Name: "/mnt/disk10/films/bluray/Rocky II (1979)", Size: 39661529553, Path: "films/bluray/Rocky II (1979)"},
		&model.Item{Name: "/mnt/disk10/films/bluray/Rocky III (1982)", Size: 31679530317, Path: "films/bluray/Rocky III (1982)"},
		&model.Item{Name: "/mnt/disk10/films/bluray/Rocky IV (1985)", Size: 30884579867, Path: "films/bluray/Rocky IV (1985)"},
		&model.Item{Name: "/mnt/disk10/films/bluray/Rocky V (1990)", Size: 36051697059, Path: "films/bluray/Rocky V (1990)"},
		&model.Item{Name: "/mnt/disk10/films/bluray/Run Lola Run (1998)", Size: 29867813190, Path: "films/bluray/Run Lola Run (1998)"},
		&model.Item{Name: "/mnt/disk10/films/bluray/Shame the Devil (2013)", Size: 24183700438, Path: "films/bluray/Shame the Devil (2013)"},
		&model.Item{Name: "/mnt/disk10/films/bluray/Silent Night, Zombie Night (2011)", Size: 16805964909, Path: "films/bluray/Silent Night, Zombie Night (2011)"},
		&model.Item{Name: "/mnt/disk10/films/bluray/Small Soldiers (1998)", Size: 39396615482, Path: "films/bluray/Small Soldiers (1998)"},
		&model.Item{Name: "/mnt/disk10/films/bluray/Smokin' Aces (2006)", Size: 33476281141, Path: "films/bluray/Smokin' Aces (2006)"},
		&model.Item{Name: "/mnt/disk10/films/bluray/Snitch (2013)", Size: 45977458647, Path: "films/bluray/Snitch (2013)"},
		&model.Item{Name: "/mnt/disk10/films/bluray/Stand By Me (1986)", Size: 30073157086, Path: "films/bluray/Stand By Me (1986)"},
		&model.Item{Name: "/mnt/disk10/films/bluray/The Best Man Holiday (2013)", Size: 47827610635, Path: "films/bluray/The Best Man Holiday (2013)"},
		&model.Item{Name: "/mnt/disk10/films/bluray/The Black Stallion (1979)", Size: 38666906415, Path: "films/bluray/The Black Stallion (1979)"},
		&model.Item{Name: "/mnt/disk10/films/bluray/The Contractor (2007)", Size: 30579884276, Path: "films/bluray/The Contractor (2007)"},
		&model.Item{Name: "/mnt/disk10/films/bluray/The Fierce Wife (2012)", Size: 39742941525, Path: "films/bluray/The Fierce Wife (2012)"},
		&model.Item{Name: "/mnt/disk10/films/bluray/The Four Feathers (1939)", Size: 23707329338, Path: "films/bluray/The Four Feathers (1939)"},
		&model.Item{Name: "/mnt/disk10/films/bluray/The Messenger The Story of Joan of Arc (1999)", Size: 46051230224, Path: "films/bluray/The Messenger The Story of Joan of Arc (1999)"},
		&model.Item{Name: "/mnt/disk10/films/bluray/The Muppets (2011)", Size: 36206402863, Path: "films/bluray/The Muppets (2011)"},
		&model.Item{Name: "/mnt/disk10/films/bluray/The Nut Job (2014)", Size: 75877121127, Path: "films/bluray/The Nut Job (2014)"},
		&model.Item{Name: "/mnt/disk10/films/bluray/The Secret Life of Walter Mitty (2013)", Size: 42118484958, Path: "films/bluray/The Secret Life of Walter Mitty (2013)"},
		&model.Item{Name: "/mnt/disk10/films/bluray/The Wolf of Wall Street (2013)", Size: 46757209164, Path: "films/bluray/The Wolf of Wall Street (2013)"},
		&model.Item{Name: "/mnt/disk10/films/bluray/Thick as Thieves (2009)", Size: 23486681428, Path: "films/bluray/Thick as Thieves (2009)"},
		&model.Item{Name: "/mnt/disk10/films/bluray/TimeCop (1994)", Size: 18624669685, Path: "films/bluray/TimeCop (1994)"},
		&model.Item{Name: "/mnt/disk10/films/bluray/Tremors 3 Back to Perfection  (2002)", Size: 34886345149, Path: "films/bluray/Tremors 3 Back to Perfection  (2002)"},
		&model.Item{Name: "/mnt/disk10/films/bluray/Tremors 4 The Legend Begins (2004)", Size: 33743155599, Path: "films/bluray/Tremors 4 The Legend Begins (2004)"},
		&model.Item{Name: "/mnt/disk10/films/bluray/We Were Soldiers (2002)", Size: 42824026365, Path: "films/bluray/We Were Soldiers (2002)"},
		&model.Item{Name: "/mnt/disk10/films/bluray/What's Eating Gilbert Grape (1993)", Size: 23118653830, Path: "films/bluray/What's Eating Gilbert Grape (1993)"},
		&model.Item{Name: "/mnt/disk10/films/bluray/When a Man Falls (2007)", Size: 20326300507, Path: "films/bluray/When a Man Falls (2007)"},
		&model.Item{Name: "/mnt/disk10/films/bluray/Winter's Bone (2010)", Size: 24272638508, Path: "films/bluray/Winter's Bone (2010)"},
		&model.Item{Name: "/mnt/disk10/films/bluray/You Only Live Twice (1967)", Size: 47879483403, Path: "films/bluray/You Only Live Twice (1967)"},
		&model.Item{Name: "/mnt/disk10/films/bluray/Young and Beautiful (2013)", Size: 36128559084, Path: "films/bluray/Young and Beautiful (2013)"},
		&model.Item{Name: "/mnt/disk10/films/bluray/Zaytoun (2012)", Size: 34500054124, Path: "films/bluray/Zaytoun (2012)"},
	}

	srcDisk := disks[9]

	sort.Sort(model.ByFree(disks))

	srcDisk.NewFree = srcDisk.Free

	var bytesToMove uint64

	for _, disk := range disks {
		disk.NewFree = disk.Free
		if disk.Path != srcDisk.Path {
			log.Println("--------------------------- new search ------------------------")
			packer := lib.NewKnapsack(disk, folders)
			bin := packer.BestFit()
			if bin != nil {
				// srcDiskSizeFreeFinal += bin.Size
				srcDisk.NewFree += bin.Size
				disk.NewFree -= bin.Size
				bytesToMove += bin.Size

				removeFolders(folders, bin.Items)
			}
		}
	}

	for _, disk := range disks {
		disk.Print()
	}

	fmt.Println("=========================================================")
	fmt.Println(fmt.Sprintf("Results for %s", srcDisk.Path))
	fmt.Println(fmt.Sprintf("Original Free Space: %s", lib.ByteSize(srcDisk.Free)))
	fmt.Println(fmt.Sprintf("Final Free Space: %s", lib.ByteSize(srcDisk.NewFree)))
	fmt.Println(fmt.Sprintf("Gained Space: %s", lib.ByteSize(srcDisk.NewFree-srcDisk.Free)))
	fmt.Println("---------------------------------------------------------")

	// for _, disk := range disks {
	// 	disk.Print()
	// }
}

func removeFolders(folders []*model.Item, list []*model.Item) []*model.Item {
	w := 0 // write index

loop:
	for _, fld := range folders {
		for _, itm := range list {
			if itm.Name == fld.Name {
				continue loop
			}
		}
		folders[w] = fld
		w++
	}

	return folders[:w]
}
