package domain

type State struct {
	Status    int64        `json:"status"`
	Unraid    *Unraid      `json:"unraid"`
	Operation *Operation   `json:"operation"`
	History   []*Operation `json:"history"`
}
