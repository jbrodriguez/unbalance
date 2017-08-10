package model

type Command struct {
	Src     string
	Dst     string
	Path    string
	Size    int64
	WorkDir string
}
