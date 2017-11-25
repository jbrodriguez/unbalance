package domain

import (
	"time"
)

// Operation - represents the different types of operationns in the app
// OpKind can be
// OP_SCATTER_MOVE
// OP_SCATTER_COPY
// OP_SCATTER_VALIDATE
// OP_GATHER_MOVE
type Operation struct {
	ID       string    `json:"id"`
	OpKind   int64     `json:"opKind"`
	Started  time.Time `json:"started"`
	Finished time.Time `json:"finished"`

	// transfer section
	BytesToTransfer  int64      `json:"bytesToTransfer"`
	BytesTransferred int64      `json:"bytesTransferred"`
	DryRun           bool       `json:"dryRun"`
	RsyncFlags       []string   `json:"rsyncFlags"`
	RsyncStrFlags    string     `json:"rsyncStrFlags"`
	Commands         []*Command `json:"commands"`

	// progress section
	Completed     float64 `json:"completed"`
	Speed         float64 `json:"speed"`
	Remaining     string  `json:"remaining"`
	DeltaTransfer int64   `json:"deltaTransfer"`
	Line          string  `json:"line"`
}
