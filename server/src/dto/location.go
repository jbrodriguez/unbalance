package dto

import "jbrodriguez/unbalance/server/src/domain"

type Location struct {
	Disks    map[string]*domain.Disk `json:"disks"`
	Presence map[string]string       `json:"presence"`
}
