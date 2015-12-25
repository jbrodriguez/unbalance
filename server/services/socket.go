package services

import (
	"github.com/jbrodriguez/mlog"
	"github.com/jbrodriguez/pubsub"
	"golang.org/x/net/websocket"
	"jbrodriguez/unbalance/server/dto"
	"jbrodriguez/unbalance/server/lib"
	"jbrodriguez/unbalance/server/net"
	// "net/http"
)

// const (
// 	//	pongWait       = 60 * time.Second
// 	bufferSize = 8192
// )

type Socket struct {
	Service

	bus      *pubsub.PubSub
	settings *lib.Settings

	mailbox chan *pubsub.Mailbox

	// registered connections
	// connections map[*Connection]bool
	pool map[*net.Connection]bool

	// inbound messages from connections
	// broadcast chan *pubsub.Message

	// // inbound messages from connections
	// emit chan *pubsub.Message

	// // register requests from connection
	// register chan *Connection

	// // unregister request from connection
	// unregister chan *Connection
}

func NewSocket(bus *pubsub.PubSub, settings *lib.Settings) *Socket {
	socket := &Socket{
		bus:      bus,
		settings: settings,
		pool:     make(map[*net.Connection]bool),
		// register:    make(chan *Connection),
		// unregister:  make(chan *Connection),

		// broadcast: bus.Sub("socket:broadcast"),
		// emit:      bus.Sub("socket:emit"),
	}

	socket.init()
	return socket
}

// func (s *Socket) handler(w http.ResponseWriter, r *http.Request) {
// 	ws, err := websocket.Upgrade(w, r, nil, bufferSize, bufferSize)
// 	if _, ok := err.(websocket.HandshakeError); ok {
// 		http.Error(w, "Not a websocket handshake", 400)
// 		return
// 	} else if err != nil {
// 		mlog.Error(err)
// 		return
// 	}

// 	// mlog.Info("opening connection ... ")
// 	c := &Connection{send: make(chan *dto.Packet, 50), ws: ws, hub: s}
// 	s.register <- c
// 	// defer func() { s.unregister <- c }()

// 	go c.writer()
// 	c.reader()
// }

func (s *Socket) Start() {
	mlog.Info("starting service Socket ...")

	s.mailbox = s.register(s.bus, "/add/connection", s.addConnection)
	s.registerAdditional(s.bus, "storage:broadcast", s.broadcast, s.mailbox)

	go s.react()
}

func (s *Socket) Stop() {
	mlog.Info("stopped service Socket ...")
}

func (s *Socket) react() {
	for mbox := range s.mailbox {
		// mlog.Info("Core:Topic: %s", mbox.Topic)
		s.dispatch(mbox.Topic, mbox.Content)
	}
}

func (s *Socket) addConnection(msg *pubsub.Message) {
	ws := msg.Payload.(*websocket.Conn)
	conn := net.NewConnection(ws, s.onMessage, s.onClose)
	s.pool[conn] = true
	go conn.Read()
}

func (s *Socket) onMessage(packet *dto.Packet) {
	s.bus.Pub(&pubsub.Message{Payload: packet.Payload}, packet.Topic)
}

func (s *Socket) onClose(c *net.Connection, err error) {
	mlog.Warning("closing socket (%+v): %s", c, err)
	if _, ok := s.pool[c]; ok {
		delete(s.pool, c)
	}
}

func (s *Socket) broadcast(msg *pubsub.Message) {
	packet := msg.Payload.(*dto.Packet)
	for conn, _ := range s.pool {
		conn.Write(packet)
	}
}
