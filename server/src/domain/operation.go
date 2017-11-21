package domain

import (
	"time"
)

// Operation - represents the different types of operationns in the app
// OpKind can be
// ??? scatterCalc
// ??? gatherFindTarget
// ??? - move
// ??? - copy
// ??? - validate
type Operation struct {
	OpKind                int64             `json:"opKind"`
	Started               time.Time         `json:"started"`
	Finished              time.Time         `json:"finished"`
	ChosenFolders         []string          `json:"chosenFolders"`
	FoldersNotTransferred []string          `json:"-"`
	OwnerIssue            int64             `json:"ownerIssue"`
	GroupIssue            int64             `json:"groupIssue"`
	FolderIssue           int64             `json:"folderIssue"`
	FileIssue             int64             `json:"fileIssue"`
	BytesToTransfer       int64             `json:"bytesToTransfer"`
	DryRun                bool              `json:"-"`
	RsyncFlags            []string          `json:"-"`
	RsyncStrFlags         string            `json:"-"`
	Commands              []Command         `json:"commands"`
	BytesTransferred      int64             `json:"bytesTransferred"`
	VDisks                map[string]*VDisk `json:"vdisks"`
}
