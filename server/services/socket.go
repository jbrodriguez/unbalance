package services

import (
	"github.com/gorilla/websocket"
	"github.com/jbrodriguez/mlog"
	"github.com/jbrodriguez/pubsub"
	"jbrodriguez/unbalance/server/dto"
	"jbrodriguez/unbalance/server/lib"
	"net/http"
)

const (
	//	pongWait       = 60 * time.Second
	bufferSize = 8192
)

type Socket struct {
	bus      *pubsub.PubSub
	settings *lib.Settings

	// registered connections
	connections map[*Connection]bool

	// inbound messages from connections
	broadcast chan *pubsub.Message

	// inbound messages from connections
	emit chan *pubsub.Message

	// register requests from connection
	register chan *Connection

	// unregister request from connection
	unregister chan *Connection
}

func NewSocket(bus *pubsub.PubSub, settings *lib.Settings) *Socket {
	return &Socket{
		bus:      bus,
		settings: settings,

		connections: make(map[*Connection]bool),
		register:    make(chan *Connection),
		unregister:  make(chan *Connection),

		// broadcast: bus.Sub("socket:broadcast"),
		// emit:      bus.Sub("socket:emit"),
	}
}

func (s *Socket) handler(w http.ResponseWriter, r *http.Request) {
	ws, err := websocket.Upgrade(w, r, nil, bufferSize, bufferSize)
	if _, ok := err.(websocket.HandshakeError); ok {
		http.Error(w, "Not a websocket handshake", 400)
		return
	} else if err != nil {
		mlog.Error(err)
		return
	}

	// mlog.Info("opening connection ... ")
	c := &Connection{send: make(chan *dto.Packet, 50), ws: ws, hub: s}
	s.register <- c
	// defer func() { s.unregister <- c }()

	go c.writer()
	c.reader()
}

func (s *Socket) Start() {
	mlog.Info("starting service Socket ...")
	go s.react()
}

func (s *Socket) Stop() {
	mlog.Info("stopped service Socket ...")
}

func (s *Socket) react() {
	for {
		select {
		case c := <-s.register:
			s.connections[c] = true
		case c := <-s.unregister:
			if _, ok := s.connections[c]; ok {
				delete(s.connections, c)
				close(c.send)
			}
		case m := <-s.broadcast:
			// mlog.Info("broadcasting %v [%v]", m, m.Payload)

			for c := range s.connections {
				c.send <- m.Payload.(*dto.Packet)
				// select {
				// case c.send <- m.Payload.(*dto.Packet):
				// 	mlog.Info("after c.send")
				// 	default:
				// 		mlog.Info("default.close")
				// 		close(c.send)
				// 		delete(s.connections, c)
				// 	// go c.ws.Close()
				// }
			}
			//		case _ := <-s.emit:
			//go c.calculateBestFit(msg)
		}
	}
}
