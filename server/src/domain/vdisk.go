package domain

type VDisk struct {
	Path        string `json:"path"`
	PlannedFree int64  `json:"plannedFree"`
	Bin         *Bin   `json:"-"`
	Src         bool   `json:"src"`
	Dst         bool   `json:"dst"`
}
