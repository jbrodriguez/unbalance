package domain

type Command struct {
	Src         string `json:"src"`
	Dst         string `json:"dst"`
	WorkDir     string `json:"workdir"`
	Size        int64  `json:"size"`
	Transferred int64  `json:"transferred"`
}
