package dto

import "jbrodriguez/unbalance/server/src/domain"

// Location -
type Location struct {
	Disks    map[string]*domain.Disk `json:"disks"`
	Presence map[string]string       `json:"presence"`
}
