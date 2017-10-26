package domain

type State struct {
	Unraid    *Unraid    `json:"unraid"`
	Operation *Operation `json:"operation"`
}
