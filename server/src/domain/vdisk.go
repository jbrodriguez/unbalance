package domain

// VDisk -
type VDisk struct {
	Path        string `json:"path"`
	PlannedFree int64  `json:"plannedFree"`
	Bin         *Bin   `json:"bin"`
	Src         bool   `json:"src"`
	Dst         bool   `json:"dst"`
}
