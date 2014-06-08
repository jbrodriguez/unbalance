package message

import (
	"apertoire.net/unbalance/model"
	"encoding/json"
)

type FitData struct {
	SourceDisk string
	TargetDisk string
	Reply      chan string
}

type Disks struct {
	Reply chan []*model.Disk
}

type Request struct {
	Id     int              `json: "id"`
	Method string           `json: "method"`
	Params *json.RawMessage `json: "params"`
}

type Reply struct {
	Id     int              `json: "id"`
	Result *json.RawMessage `json: "result"`
	Error  *json.RawMessage `json: "error"`
}

// type DiskReply struct {
// 	Id     int               `json: "id"`
// 	Result []*model.Disk     `json: "result"`
// 	Error  map[string]string `json: "error"`
// }
