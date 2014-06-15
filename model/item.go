package model

type Item struct {
	Name string
	Size uint64
	Path string
}

type BySize []*Item

func (s BySize) Len() int           { return len(s) }
func (s BySize) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }
func (s BySize) Less(i, j int) bool { return s[i].Size > s[j].Size }
