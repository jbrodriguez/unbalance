package domain

type ScatterSetup struct {
	Source   string   `json:"source"`
	Targets  []string `json:"targets"`
	Selected []string `json:"selected"`
}

type GatherSetup struct {
	Selected []string `json:"selected"`
}
