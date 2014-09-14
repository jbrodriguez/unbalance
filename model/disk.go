package model

import (
	"fmt"
	"strings"
)

const (
	BYTE     = 1.0
	KILOBYTE = 1000 * BYTE
	MEGABYTE = 1000 * KILOBYTE
	GIGABYTE = 1000 * MEGABYTE
	TERABYTE = 1000 * GIGABYTE
)

type Disk struct {
	Id              int    `json:"id"`
	Name            string `json:"name"`
	Path            string `json:"path"`
	Device          string `json:"device"`
	Free            uint64 `json:"free"`
	FreeMinusBuffer uint64 `json:"freeMinusBuffer"`
	NewFree         uint64 `json:"newFree"`
	Size            uint64 `json:"size"`
	Serial          string `json:"serial"`
	Status          string `json:"status"`
	Bin             *Bin   `json:"-"`
}

// func (self *Disk) Print() {
// 	// this disk was not assigned to a bin
// 	if self.Bin != nil {
// 		fmt.Println("=========================================================")
// 		fmt.Println(fmt.Sprintf("[%d/%d] %2.2f%% (%s)", self.Bin.Size, self.Free, (float64(self.Bin.Size)/float64(self.Free))*100, self.Path))
// 		fmt.Println("---------------------------------------------------------")
// 		self.Bin.Print()
// 		fmt.Println("---------------------------------------------------------")
// 		fmt.Println("")
// 	} else {
// 		fmt.Println("=========================================================")
// 		fmt.Println(fmt.Sprintf("[0/%d] 0%% (%s)", self.Free, self.Path))
// 		fmt.Println("---------------------------------------------------------")
// 		fmt.Println("---------------------------------------------------------")
// 		fmt.Println("")
// 	}
// }

func (self *Disk) Print() {
	if self.Bin != nil {
		fmt.Printf("Disk %s: %s\n", self.Path, HumanBytes(self.Free-self.Bin.Size))
	} else {
		fmt.Printf("Disk %s: no diggity (%s)\n", self.Path, HumanBytes(self.Free))
	}
}

type ByFree []*Disk

func (s ByFree) Len() int           { return len(s) }
func (s ByFree) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }
func (s ByFree) Less(i, j int) bool { return s[i].Free > s[j].Free }

func HumanBytes(bytes uint64) string {
	unit := ""
	value := float32(bytes)

	switch {
	case bytes >= TERABYTE:
		unit = "T"
		value = value / TERABYTE
	case bytes >= GIGABYTE:
		unit = "G"
		value = value / GIGABYTE
	case bytes >= MEGABYTE:
		unit = "M"
		value = value / MEGABYTE
	case bytes >= KILOBYTE:
		unit = "K"
		value = value / KILOBYTE
	case bytes == 0:
		return "0"
	}

	stringValue := fmt.Sprintf("%.1f", value)
	stringValue = strings.TrimSuffix(stringValue, ".0")
	return fmt.Sprintf("%s%s", stringValue, unit)
}
