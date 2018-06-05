package domain

import (
	"time"
)

// Unraid -
type Unraid struct {
	NumDisks     int64     `json:"numDisks"`
	NumProtected int64     `json:"numProtected"`
	Synced       time.Time `json:"synced"`
	SyncErrs     int64     `json:"syncErrs"`
	Resync       int64     `json:"resync"`
	ResyncPos    int64     `json:"resyncPos"`
	State        string    `json:"state"`
	Size         int64     `json:"size"`
	Free         int64     `json:"free"`
	Disks        []*Disk   `json:"disks"`
	BlockSize    int64     `json:"-"`
}
