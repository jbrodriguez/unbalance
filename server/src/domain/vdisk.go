package domain

// VDisk -
type VDisk struct {
	Path        string `json:"path"`
	PlannedFree uint64 `json:"plannedFree"`
	Bin         *Bin   `json:"bin"`
	Src         bool   `json:"src"`
	Dst         bool   `json:"dst"`
}
