package dto

type Entry struct {
	Path  string  `json:"path"`
	Nodes []*Node `json:"nodes"`
}

type Node struct {
	Type string `json:"type"`
	Path string `json:"path"`
}
