package message

import (
	"apertoire.net/unbalance/helper"
	"encoding/json"
)

type FitData struct {
	SourceDisk string
	TargetDisk string
	Reply      chan string
}

type Status struct {
	Reply chan *helper.Unraid
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
