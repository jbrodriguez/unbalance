package domain

import (
	"time"
)

// Unraid -
type Unraid struct {
	NumDisks     uint64    `json:"numDisks"`
	NumProtected uint64    `json:"numProtected"`
	Synced       time.Time `json:"synced"`
	SyncErrs     uint64    `json:"syncErrs"`
	Resync       uint64    `json:"resync"`
	ResyncPos    uint64    `json:"resyncPos"`
	State        string    `json:"state"`
	Size         uint64    `json:"size"`
	Free         uint64    `json:"free"`
	Disks        []*Disk   `json:"disks"`
	BlockSize    uint64    `json:"-"`
}
