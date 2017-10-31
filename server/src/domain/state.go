package domain

type State struct {
	Status     int64        `json:"status"`
	Unraid     *Unraid      `json:"unraid"`
	Operations []*Operation `json:"-"`
	Operation  *Operation   `json:"operation"`
}
