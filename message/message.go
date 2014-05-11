package message

type FitData struct {
	SourceDisk string
	TargetDisk string
	Reply      chan string
}
