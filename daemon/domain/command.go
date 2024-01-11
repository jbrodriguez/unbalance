package domain

// Command -
type Command struct {
	ID          string `json:"id"`
	Src         string `json:"src"`
	Dst         string `json:"dst"`
	Entry       string `json:"entry"`
	Size        uint64 `json:"size"`
	Transferred uint64 `json:"transferred"`
	Status      int    `json:"status"`
}
