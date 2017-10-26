package domain

import (
	"time"
)

type Operation struct {
	PrevState             uint64            `json:"prevState"`
	OpState               uint64            `json:"opState"`
	RsyncFlags            []string          `json:"-"`
	RsyncStrFlags         string            `json:"-"`
	Commands              []Command         `json:"commands"`
	BytesToTransfer       int64             `json:"bytesToTransfer"`
	BytesTransferred      int64             `json:"bytesTransferred"`
	Started               time.Time         `json:"started"`
	Finished              time.Time         `json:"finished"`
	FoldersNotTransferred []string          `json:"-"`
	OwnerIssue            int64             `json:"ownerIssue"`
	GroupIssue            int64             `json:"groupIssue"`
	FolderIssue           int64             `json:"folderIssue"`
	FileIssue             int64             `json:"fileIssue"`
	DryRun                bool              `json:"-"`
	VDisks                map[string]*VDisk `json:"vdisks"`
}
