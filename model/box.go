package model

import (
	"time"
)

type Box struct {
	NumDisks     uint64    `json:"numDisks"`
	NumProtected uint64    `json:"numProtected"`
	Synced       time.Time `json:"synced"`
	SyncErrs     uint64    `json:"syncErrs"`
	Resync       uint64    `json:"resync"`
	ResyncPrcnt  uint64    `json:"resyncPrct"`
	ResyncPos    uint64    `json:"resyncPos"`
	State        string    `json:"state"`
	Size         uint64    `json:"size"`
	Free         uint64    `json:"free"`
}
