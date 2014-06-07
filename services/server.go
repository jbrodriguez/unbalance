package services

import (
	"apertoire.net/unbalance/bus"
	"apertoire.net/unbalance/message"
	"code.google.com/p/go.net/websocket"
	"log"
	"net/http"
)

type Handler func(id int, msg *message.Message)

type Server struct {
	Bus *bus.Bus

	sockets map[int]*Socket

	vtable map[string]Handler

	addCh chan *Socket
	delCh chan *Socket
	errCh chan error
}

func (self *Server) getDisks(id int, msg *message.Message) {
	// fire this message onto the bus, fire and forget
	log.Println("Omaha !!!")
}

func (self *Server) Start() {
	log.Printf("starting Server service ...")

	self.sockets = make(map[int]*Socket)
	self.vtable = make(map[string]Handler)

	self.addCh = make(chan *Socket)
	self.delCh = make(chan *Socket)
	self.errCh = make(chan error)

	self.Handle("/api/v1/get/disks", self.getDisks)

	// start the websocket listener, and handles incoming websocket connections
	go self.react()

	http.Handle("/", http.FileServer(http.Dir("ui")))

	go func() {
		log.Fatal(http.ListenAndServe(":6237", nil))
	}()

	log.Printf("Server service listening on :6237")
}

func (self *Server) Stop() {
	log.Printf("Server service stopped")
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

func (self *Server) Dispatch(id int, msg *message.Message) {
	pattern := msg.Method
	handler := self.vtable[pattern]

	log.Println("amma dispatch you: ", pattern)

	handler(id, msg)
}

func (self *Server) react() {
	onConnected := func(ws *websocket.Conn) {
		log.Println("socket connected")

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
	log.Println("created Handler")

	for {
		select {
		// case msg := <-self.Bus.NewConnection:

		case socket := <-self.addCh:
			self.sockets[socket.id] = socket
		}
	}
}
