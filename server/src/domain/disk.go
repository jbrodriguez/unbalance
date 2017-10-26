package domain

type Disk struct {
	ID     int64  `json:"id"`
	Name   string `json:"name"`
	Path   string `json:"path"`
	Device string `json:"device"`
	Type   string `json:"type"`
	FsType string `json:"fsType"`
	Free   int64  `json:"free"`
	Size   int64  `json:"size"`
	Serial string `json:"serial"`
	Status string `json:"status"`
}
