package dto

type MessageIn struct {
	Id      string `json:"-"`
	Topic   string `json:"topic"`
	Payload string `json:"payload"`
}

type MessageOut struct {
	Id      string      `json:"-"`
	Topic   string      `json:"topic"`
	Payload interface{} `json:"payload"`
}
