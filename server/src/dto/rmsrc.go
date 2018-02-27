package dto

import "jbrodriguez/unbalance/server/src/domain"

// RmSrc -
type RmSrc struct {
	Operation *domain.Operation `json:"operation"`
	ID        string            `json:"id"`
}
