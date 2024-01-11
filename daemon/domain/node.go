package domain

// {label: 'films', value: 'films', isLeaf: false, children: [{label: 'bluray'},{label: 'blurip'}]},
type Node struct {
	ID     string `json:"id"`
	Label  string `json:"label"`
	Leaf   bool   `json:"leaf"`
	Parent string `json:"parent"`
}
