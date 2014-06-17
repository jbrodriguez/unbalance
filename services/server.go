package services

import (
	"apertoire.net/unbalance/bus"
	"apertoire.net/unbalance/lib"
	"apertoire.net/unbalance/message"
	"apertoire.net/unbalance/model"
	"code.google.com/p/go.net/websocket"
	"encoding/json"
	"fmt"
	"github.com/golang/glog"
	"io"
	"net/http"
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

type Handler func(id int, msg *message.Request)

type Server struct {
	Bus *bus.Bus

	sockets map[int]*Socket
	vtable  map[string]Handler

	addCh chan *Socket
	delCh chan *Socket
	errCh chan error
}

func (self *Server) getStatus(id int, msg *message.Request) {
	// fire this message onto the bus, wait for the reply
	glog.Info("Omaha !!!")

	event := &message.StorageStatus{make(chan *model.Unraid)}
	self.Bus.GetStatus <- event
	unraid := <-event.Reply

	// b, err := json.Marshal(disks)
	// if err != nil {
	// 	glog.Info("errored out: ", err)
	// } else {
	// 	glog.Info(string(b))
	// }

	// m := json.RawMessage(b)
	data, err := lib.WriteJson(unraid)
	if err != nil {
		glog.Info("errored out: ", err)
	}

	reply := &message.Reply{Id: msg.Id, Result: &data}
	self.sockets[id].Write(reply)
	// self.sockets[id].Write(&model.Disk{Path: "/mnt/disk", Free: 434983434})

}

func (self *Server) getBestFit(id int, msg *message.Request) {
	params := new(message.BestFit)
	err := json.Unmarshal(*msg.Params, params)
	if err != nil {
		glog.Fatal("motherfucker: ", err)
	}
	glog.Infof("this is all you: %+v", params)

	params.Reply = make(chan *model.Unraid)

	self.Bus.GetBestFit <- params
	unraid := <-params.Reply

	data, err := lib.WriteJson(unraid)
	if err != nil {
		glog.Info("errored out: ", err)
	}

	reply := &message.Reply{Id: msg.Id, Result: &data}
	self.sockets[id].Write(reply)
}

func (self *Server) Start() {
	glog.Info("starting Server service ...")

	self.sockets = make(map[int]*Socket)
	self.vtable = make(map[string]Handler)

	self.addCh = make(chan *Socket)
	self.delCh = make(chan *Socket)
	self.errCh = make(chan error)

	self.Handle("/api/v1/get/status", self.getStatus)
	self.Handle("/api/v1/get/bestFit", self.getBestFit)

	// start the websocket listener, and handles incoming websocket connections
	go self.react()

	http.Handle("/", http.FileServer(http.Dir("ui")))

	go func() {
		glog.Fatal(http.ListenAndServe(":6237", nil))
	}()

	glog.Info("Server service listening on :6237")
}

func (self *Server) Stop() {
	glog.Info("Server service stopped")
}

func (self *Server) Add(socket *Socket) {
	// go func() {
	self.addCh <- socket
	// }()
}

func (self *Server) Del(socket *Socket) {
	self.delCh <- socket
}

func (self *Server) Err(err error) {
	self.errCh <- err
}

func (self *Server) Handle(pattern string, handler Handler) {
	self.vtable[pattern] = handler
}

func (self *Server) Dispatch(id int, msg *message.Request) {
	pattern := msg.Method
	handler := self.vtable[pattern]

	glog.Info("amma dispatch you: ", pattern)

	handler(id, msg)
}

func (self *Server) react() {
	onConnected := func(ws *websocket.Conn) {
		glog.Info("socket connected")

		defer func() {
			err := ws.Close()
			if err != nil {
				self.errCh <- err
			}
		}()

		socket := NewSocket(ws, self)
		self.Add(socket)
		socket.Listen()
	}

	http.Handle("/api", websocket.Handler(onConnected))
	glog.Info("created Handler")

	for {
		select {
		// case msg := <-self.Bus.NewConnection:

		case socket := <-self.addCh:
			self.sockets[socket.id] = socket
		}
	}
}
