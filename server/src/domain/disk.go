package domain

// Disk -
type Disk struct {
	ID          int64  `json:"id"`
	Name        string `json:"name"`
	Path        string `json:"path"`
	Device      string `json:"device"`
	Type        string `json:"type"`
	FsType      string `json:"fsType"`
	Free        uint64 `json:"free"`
	Size        uint64 `json:"size"`
	Serial      string `json:"serial"`
	Status      string `json:"status"`
	BlocksTotal uint64 `json:"-"`
	BlocksFree  uint64 `json:"-"`
}
