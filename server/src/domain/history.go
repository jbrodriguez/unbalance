package domain

import (
	"time"
)

// History -
type History struct {
	Version     int                   `json:"version"`
	LastChecked time.Time             `json:"lastChecked"`
	Items       map[string]*Operation `json:"items"`
	Order       []string              `json:"order"`
}
