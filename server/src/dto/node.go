package dto

// Entry -
type Entry struct {
	Path  string  `json:"path"`
	Nodes []*Node `json:"nodes"`
}

// Node -
type Node struct {
	Type string `json:"type"`
	Path string `json:"path"`
}
