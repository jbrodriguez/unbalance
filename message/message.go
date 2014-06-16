package message

import (
	"apertoire.net/unbalance/lib"
	"encoding/json"
)

type BestFit struct {
	SourceDisk string           `json:"fromDisk"`
	TargetDisk string           `json:"toDisk"`
	Reply      chan *lib.Unraid `json:"-"`
}

type Status struct {
	Reply chan *lib.Unraid
}

type Request struct {
	Id     int              `json:"id"`
	Method string           `json:"method"`
	Params *json.RawMessage `json:"params"`
}

type Reply struct {
	Id     int              `json:"id"`
	Result *json.RawMessage `json:"result"`
	Error  *json.RawMessage `json:"error"`
}

type Progress {
	TotalSize uint64
	TotalCopied uint64
	CurrentFile string
	CurrentSize uint64
	CurrentCopied uint64
}

// type DiskReply struct {
// 	Id     int               `json: "id"`
// 	Result []*model.Disk     `json: "result"`
// 	Error  map[string]string `json: "error"`
// }
