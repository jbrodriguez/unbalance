package dto

import (
	"jbrodriguez/unbalance/server/src/domain"
)

type Transfer struct {
	Operation *domain.Operation `json:"operation"`
	Progress  *Progress         `json:"progress"`
}
