package model

import (
	"time"
)

type Condition struct {
	NumDisks     int64     `json:"numDisks"`
	NumProtected int64     `json:"numProtected"`
	Synced       time.Time `json:"synced"`
	SyncErrs     int64     `json:"syncErrs"`
	Resync       int64     `json:"resync"`
	ResyncPos    int64     `json:"resyncPos"`
	State        string    `json:"state"`
	Size         int64     `json:"size"`
	Free         int64     `json:"free"`
	NewFree      int64     `json:"newFree"`
}
