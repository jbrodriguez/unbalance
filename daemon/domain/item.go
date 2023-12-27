package domain

// Item -
type Item struct {
	Name       string `json:"name"`
	Size       uint64 `json:"size"`
	Path       string `json:"path"`
	Location   string `json:"location"`
	BlocksUsed uint64 `json:"blocksUsed"`
}

// // BySize -
// type BySize []*Item

// func (s BySize) Len() int           { return len(s) }
// func (s BySize) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }
// func (s BySize) Less(i, j int) bool { return s[i].Size > s[j].Size }
