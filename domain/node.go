package domain

// {label: 'films', checkbox: true, collapsed: true, collapsible: true, children: [{label: 'bluray'},{label: 'blurip'}]},
type Node struct {
	Label     string `json:"label"`
	Checkbox  bool   `json:"checkbox"`
	Collapsed bool   `json:"collapsed"`
	Children  []Node `json:"children"`
	Path      string `json:"path"`
}
