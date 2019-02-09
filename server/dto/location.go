package dto

import "unbalance/domain"

// Location -
type Location struct {
	Disks    map[string]*domain.Disk `json:"disks"`
	Presence map[string]string       `json:"presence"`
}
