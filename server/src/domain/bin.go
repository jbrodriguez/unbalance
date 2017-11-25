package domain

// import (
// 	// "github.com/jbrodriguez/mlog"
// 	// "jbrodriguez/unbalance/server/src/lib"
// )

// Bin -
type Bin struct {
	Size  int64
	Items []*Item
}

// Add -
func (b *Bin) Add(item *Item) {
	b.Items = append(b.Items, item)
	b.Size += item.Size
}

// // Print -
// func (b *Bin) Print() {
// 	for _, item := range b.Items {
// 		mlog.Info("[%s] %s\n", lib.ByteSize(item.Size), item.Name)
// 	}
// }

// // ByFilled -
// type ByFilled []*Bin

// func (s ByFilled) Len() int           { return len(s) }
// func (s ByFilled) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }
// func (s ByFilled) Less(i, j int) bool { return s[i].Size > s[j].Size }
