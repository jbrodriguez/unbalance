package domain

type Entry struct {
	Path  string `json:"path"`
	Nodes []Node `json:"nodes"`
}
