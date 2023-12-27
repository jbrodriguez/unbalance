package domain

type Branch struct {
	Nodes map[string]Node `json:"nodes"`
	Order []string        `json:"order"`
}
