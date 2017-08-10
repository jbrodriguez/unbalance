package model

import (
	"time"
)

const (
	StateIdle        = 0
	StateCalc        = 1
	StateMove        = 2
	StateCopy        = 3
	StateValidate    = 4
	StateFindTargets = 5
	StateGather      = 6
)

type Operation struct {
	PrevState             uint64
	OpState               uint64
	SourceDiskName        string
	RsyncFlags            []string
	RsyncStrFlags         string
	Commands              []Command
	BytesToTransfer       int64
	BytesTransferred      int64
	Started               time.Time
	Finished              time.Time
	FoldersNotTransferred []string
	OwnerIssue            int64
	GroupIssue            int64
	FolderIssue           int64
	FileIssue             int64
	DryRun                bool
	Target                string
}
