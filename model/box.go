package model

import (
	"time"
)

type Box struct {
	NumDisks     uint64
	NumProtected uint64
	Synced       time.Time
	SyncErrs     uint64
	Resync       uint64
	ResyncPrcnt  uint64
	ResyncPos    uint64
	State        string
}
