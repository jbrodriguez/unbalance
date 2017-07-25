package model

// Item -
type Item struct {
	Name     string
	Size     int64
	Path     string
	Location string
}

// BySize -
type BySize []*Item

func (s BySize) Len() int           { return len(s) }
func (s BySize) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }
func (s BySize) Less(i, j int) bool { return s[i].Size > s[j].Size }
