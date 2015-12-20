package dto

type Calculate struct {
	SourceDisk string          `json:"sourceDisk"`
	DestDisks  map[string]bool `json:"destDisk"`
}
