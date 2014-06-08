package services

import (
	"apertoire.net/unbalance/message"
	// "apertoire.net/unbalance/model"
	"code.google.com/p/go.net/websocket"
	"fmt"
	"github.com/golang/glog"
	"io"
)

const channelBufSize = 100

var maxId = 0

type Socket struct {
	id     int
	ws     *websocket.Conn
	server *Server
	ch     chan interface{}
	doneCh chan bool
}

func NewSocket(ws *websocket.Conn, server *Server) *Socket {
	if ws == nil {
		panic("ws cannot be nil")
	}

	if server == nil {
		panic("server cannot be nil")
	}

	maxId++
	ch := make(chan interface{}, channelBufSize)
	doneCh := make(chan bool)

	return &Socket{maxId, ws, server, ch, doneCh}
}

func (self *Socket) Write(msg interface{}) {
	select {
	case self.ch <- msg:
	default:
		self.server.Del(self)
		err := fmt.Errorf("socket %d is disconnected", self.id)
		self.server.Err(err)
	}
}

func (self *Socket) Listen() {
	glog.Info("are we on the listen")
	go self.listenWrite()
	self.listenRead()
}

func (self *Socket) listenWrite() {
	glog.Info("listening write to socket")
	for {
		select {
		case msg := <-self.ch:
			err := websocket.JSON.Send(self.ws, msg)
			if err != nil {
				glog.Warning("errored out: ", err)
			}
			glog.Info("Sent: ", msg)

		case <-self.doneCh:
			self.server.Del(self)
			self.doneCh <- true
			return
		}
	}
}

func (self *Socket) listenRead() {
	glog.Info("listening read from socket")
	for {
		select {
		case <-self.doneCh:
			self.server.Del(self)
			self.doneCh <- true
			return

		default:
			var msg message.Request
			err := websocket.JSON.Receive(self.ws, &msg)
			glog.Info("is there anybody out there?: ", err)
			if err == io.EOF {
				self.doneCh <- true
			} else if err != nil {
				self.server.Err(err)
			} else {
				self.server.Dispatch(self.id, &msg)
				// websocket.JSON.Send(self.ws, &model.Disk{Path: "/mnt/disk1", Free: 48939348})
			}
		}
	}
}
