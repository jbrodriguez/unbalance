package domain

import (
	"fmt"
	"path/filepath"
	"time"

	"unbalance/daemon/logger"
)

// Plan - represents the different types of operation planning in the app
// OpKind can be
// OpScatterPlan
// OpGatherPlan
type Plan struct {
	Started time.Time `json:"started"`
	Ended   time.Time `json:"ended"`

	// plan section
	ChosenFolders []string          `json:"chosenFolders"`
	OwnerIssue    int64             `json:"ownerIssue"`
	GroupIssue    int64             `json:"groupIssue"`
	FolderIssue   int64             `json:"folderIssue"`
	FileIssue     int64             `json:"fileIssue"`
	VDisks        map[string]*VDisk `json:"vdisks"`

	// transfer section
	BytesToTransfer uint64 `json:"bytesToTransfer"`
}

// Print -
func (p *Plan) Print(order []string) {
	var vdisks string

	for _, path := range order {
		vdisk := p.VDisks[path]

		var items string

		if vdisk.Bin != nil {
			items = "\n"
			for _, item := range vdisk.Bin.Items {
				items += fmt.Sprintf("item(%s):size(%d)\n", filepath.Join(item.Location, item.Path), item.Size)
			}
		}

		vdisks += fmt.Sprintf("vdisk:path(%s):plannedFree(%d):src(%t):dst(%t):bin(%t)%s\n", vdisk.Path, vdisk.PlannedFree, vdisk.Src, vdisk.Dst, vdisk.Bin != nil, items)
	}

	logger.Blue("\nPlan\nStarted: %s\nEnded: %s\nChosenFolders: %v\nOwnerIssues: %d\nGroupIssues: %d\nFolderIssues: %d\nFileIssues: %d\nBytesToTransfer: %d\n%s",
		p.Started, p.Ended, p.ChosenFolders, p.OwnerIssue, p.GroupIssue, p.FolderIssue, p.FileIssue, p.BytesToTransfer, vdisks)
}
