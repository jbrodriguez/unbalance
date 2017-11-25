package domain

type Command struct {
	Src         string `json:"src"`
	Dst         string `json:"dst"`
	Entry       string `json:"entry"`
	Size        int64  `json:"size"`
	Transferred int64  `json:"transferred"`
}
