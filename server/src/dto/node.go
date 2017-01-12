package dto

// Entry -
type Entry struct {
	Path  string `json:"path"`
	Nodes []Node `json:"nodes"`
}

// Node -
type Node struct {
	Label     string `json:"label"`
	Checkbox  bool   `json:"checkbox"`
	Collapsed bool   `json:"collapsed"`
	Children  []Node `json:"children"`
	Path      string `json:"path"`
}

// 	{label: 'films', checkbox: true, collapsed: true, collapsible: true, children: [{label: 'bluray'},{label: 'blurip'}]},
