package domain

// {label: 'films', value: 'films', isLeaf: false, children: [{label: 'bluray'},{label: 'blurip'}]},
type Node struct {
	ID     string `json:"id"`
	Label  string `json:"label"`
	Parent string `json:"parent"`
	Leaf   bool   `json:"leaf"`
	Dir    bool   `json:"dir"`
}
