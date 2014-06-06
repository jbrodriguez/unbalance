package message

type FitData struct {
	SourceDisk string
	TargetDisk string
	Reply      chan string
}

type Message struct {
	Id     int
	Method string
	Params string
	Data   string
}
