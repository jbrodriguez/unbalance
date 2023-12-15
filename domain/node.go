package domain

// {title: 'films', key: '/mnt/films', isLeaf: false, children: [{title: 'bluray'},{title: 'blurip'}]},
type Node struct {
	Title  string `json:"title"`
	Key    string `json:"key"`
	IsLeaf bool   `json:"isLeaf"`
}
