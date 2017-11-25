package domain

import (
	"time"
)

// Plan - represents the different types of operation planning in the app
// OpKind can be
// OP_SCATTER_CALC
// OP_GATHER_CALC
type Plan struct {
	Started  time.Time `json:"started"`
	Finished time.Time `json:"finished"`

	// calculate section
	ChosenFolders         []string `json:"chosenFolders"`
	FoldersNotTransferred []string `json:"-"`
	OwnerIssue            int64    `json:"ownerIssue"`
	GroupIssue            int64    `json:"groupIssue"`
	FolderIssue           int64    `json:"folderIssue"`
	FileIssue             int64    `json:"fileIssue"`

	// transfer section
	BytesToTransfer int64             `json:"bytesToTransfer"`
	VDisks          map[string]*VDisk `json:"vdisks"`
}
