package services

import (
	"apertoire.net/unbalance/message"
	"code.google.com/p/go.net/websocket"
	"io"
	"log"
)

const channelBufSize = 100

var maxId = 0

type Socket struct {
	id     int
	ws     *websocket.Conn
	server *Server
	ch     chan *message.Message
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
	ch := make(chan *message.Message, channelBufSize)
	doneCh := make(chan bool)

	return &Socket{maxId, ws, server, ch, doneCh}
}

func (self *Socket) Listen() {
	go self.listenWrite()
	self.listenRead()
}

func (self *Socket) listenWrite() {
	log.Println("listening write to socket")
	for {
		select {
		case msg := <-self.ch:
			log.Println("Send: ", msg)
			websocket.JSON.Send(self.ws, msg)

		case <-self.doneCh:
			self.server.Del(self)
			self.doneCh <- true
			return
		}
	}
}

func (self *Socket) listenRead() {
	log.Println("listening read from socket")
	for {
		select {
		case <-self.doneCh:
			self.server.Del(self)
			self.doneCh <- true
			return

		default:
			var msg message.Message
			err := websocket.JSON.Receive(self.ws, &msg)
			if err == io.EOF {
				self.doneCh <- true
			} else if err != nil {
				self.server.Err(err)
			} else {
				self.server.Dispatch(self.id, &msg)
			}
		}
	}
}
