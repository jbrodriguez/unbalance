package services

import (
	"apertoire.net/unbalance/bus"
	"apertoire.net/unbalance/helper"
	"apertoire.net/unbalance/message"
	"code.google.com/p/go.net/websocket"
	"github.com/golang/glog"
	"net/http"
)

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

	event := &message.Status{make(chan *helper.Unraid)}
	self.Bus.GetStatus <- event
	unraid := <-event.Reply

	// b, err := json.Marshal(disks)
	// if err != nil {
	// 	glog.Info("errored out: ", err)
	// } else {
	// 	glog.Info(string(b))
	// }

	// m := json.RawMessage(b)
	data, err := helper.WriteJson(unraid)
	if err != nil {
		glog.Info("errored out: ", err)
	}

	reply := &message.Reply{Id: msg.Id, Result: &data}
	self.sockets[id].Write(reply)
	// self.sockets[id].Write(&model.Disk{Path: "/mnt/disk", Free: 434983434})

}

func (self *Server) Start() {
	glog.Info("starting Server service ...")

	self.sockets = make(map[int]*Socket)
	self.vtable = make(map[string]Handler)

	self.addCh = make(chan *Socket)
	self.delCh = make(chan *Socket)
	self.errCh = make(chan error)

	self.Handle("/api/v1/get/status", self.getStatus)

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
