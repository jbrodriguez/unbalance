package domain

type Packet struct {
	ID      string `json:"-"`
	Topic   string `json:"topic"`
	Payload any    `json:"payload"`
}
