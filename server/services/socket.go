package services

import (
	"apertoire.net/unbalance/server/dto"
	"apertoire.net/unbalance/server/model"
	"github.com/gorilla/websocket"
	"github.com/jbrodriguez/mlog"
	"github.com/jbrodriguez/pubsub"
	"net/http"
)

type Socket struct {
	bus    *pubsub.PubSub
	config *model.Config

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

func NewSocket(bus *pubsub.PubSub, config *model.Config) *Socket {
	return &Socket{
		bus:    bus,
		config: config,

		connections: make(map[*Connection]bool),
		register:    make(chan *Connection),
		unregister:  make(chan *Connection),

		broadcast: bus.Sub("socket:broadcast"),
		emit:      bus.Sub("socket:emit"),
	}
}

func (s *Socket) handler(w http.ResponseWriter, r *http.Request) {
	ws, err := websocket.Upgrade(w, r, nil, 1024, 1024)

	if _, ok := err.(websocket.HandshakeError); ok {
		http.Error(w, "Not a websocket handshake", 400)
		return
	} else if err != nil {
		mlog.Error(err)
		return
	}

	c := &Connection{send: make(chan *dto.MessageOut), ws: ws, hub: s}
	s.register <- c
	defer func() { s.unregister <- c }()

	go c.writer()
	c.reader()
}

func (s *Socket) Start() {
	mlog.Info("starting service Socket ...")
	go s.react()
}

func (s *Socket) Stop() {
	mlog.Info("stopped service Core ...")
}

func (s *Socket) react() {
	for {
		select {
		case c := <-s.register:
			s.connections[c] = true
		case c := <-s.unregister:
			delete(s.connections, c)
			close(c.send)
		case m := <-s.broadcast:
			for c := range s.connections {
				select {
				case c.send <- m.Payload.(*dto.MessageOut):
				default:
					delete(s.connections, c)
					close(c.send)
					go c.ws.Close()
				}
			}
			//		case _ := <-s.emit:
			//go c.calculateBestFit(msg)
		}
	}
}
