package dto

type BestFit struct {
	SourceDisk string          `json:"sourceDisk"`
	DestDisks  map[string]bool `json:"destDisk"`
}
