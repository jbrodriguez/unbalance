package domain

type Sizes struct {
	Disks map[string]uint64 `json:"disks"`
	Total uint64            `json:"total"`
}
