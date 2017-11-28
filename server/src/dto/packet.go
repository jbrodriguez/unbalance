package dto

// Packet -
type Packet struct {
	ID      string      `json:"-"`
	Topic   string      `json:"topic"`
	Payload interface{} `json:"payload"`
}

// Chosen -
type Chosen struct {
	Payload []string `json:"payload"`
}
