package domain

// VDisk -
type VDisk struct {
	Path        string `json:"path"`
	CurrentFree uint64 `json:"currentFree"`
	PlannedFree uint64 `json:"plannedFree"`
	Bin         *Bin   `json:"bin"`
	Src         bool   `json:"src"`
	Dst         bool   `json:"dst"`
}
