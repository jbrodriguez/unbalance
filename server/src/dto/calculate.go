package dto

// Calculate -
type Calculate struct {
	SourceDisk string          `json:"srcDisk"`
	Folders    []string        `json:"folders"`
	DestDisks  map[string]bool `json:"dstDisks"`
}
