package domain

import (
	"time"
)

// Operation - represents the different types of operationns in the app
// OpKind can be
// OpScatterMove
// OpScatterCopy
// OpScatterValidate
// OpGatherMove
type Operation struct {
	ID       string    `json:"id"`
	OpKind   uint64    `json:"opKind"`
	Started  time.Time `json:"started"`
	Finished time.Time `json:"finished"`

	// transfer section
	BytesToTransfer  uint64     `json:"bytesToTransfer"`
	BytesTransferred uint64     `json:"bytesTransferred"`
	DryRun           bool       `json:"dryRun"`
	RsyncArgs        []string   `json:"rsyncArgs"`
	RsyncStrArgs     string     `json:"rsyncStrArgs"`
	Commands         []*Command `json:"commands"`

	// progress section
	Completed     float64 `json:"completed"`
	Speed         float64 `json:"speed"`
	Remaining     string  `json:"remaining"`
	DeltaTransfer uint64  `json:"deltaTransfer"`
	Line          string  `json:"line"`

	// 60 samples of transferred bytes to calculate speed
	Samples     [60]uint64 `json:"-"`
	SampleIndex int        `json:"-"`
	PrevSample  uint64     `json:"-"`
}
