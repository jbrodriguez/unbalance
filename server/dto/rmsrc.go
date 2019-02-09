package dto

import "unbalance/domain"

// RmSrc -
type RmSrc struct {
	Operation *domain.Operation `json:"operation"`
	ID        string            `json:"id"`
}
