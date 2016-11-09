package dto

// Calculate -
type Calculate struct {
	SourceDisk string          `json:"srcDisk"`
	DestDisks  map[string]bool `json:"dstDisks"`
}
