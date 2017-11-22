package dto

type Progress struct {
	Completed     float64 `json:"completed"`
	Speed         float64 `json:"speed"`
	Remaining     string  `json:"remaining"`
	DeltaTransfer int64   `json:"deltaTransfer"`
	Line          string  `json:"line"`
}
