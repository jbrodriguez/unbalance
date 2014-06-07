package message

type FitData struct {
	SourceDisk string
	TargetDisk string
	Reply      chan string
}

type Message struct {
	Id     int      `json: "id"`
	Method string   `json: "method"`
	Params []string `json: "params"`
	Data   []string `json: "data"`
}
