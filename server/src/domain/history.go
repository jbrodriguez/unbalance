package domain

import (
	"time"
)

// History -
type History struct {
	LastChecked time.Time             `json:"lastChecked"`
	Items       map[string]*Operation `json:"items"`
	Order       []string              `json:"order"`
}
