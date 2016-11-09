package net

import (
	"golang.org/x/net/websocket"
	"jbrodriguez/unbalance/server/src/dto"
)

// MessageFunc -
type MessageFunc func(message *dto.Packet)

// CloseFunc -
type CloseFunc func(conn *Connection, err error)

// Connection -
type Connection struct {
	id        string
	ws        *websocket.Conn
	onMessage MessageFunc
	onClose   CloseFunc
}

// NewConnection -
func NewConnection(ws *websocket.Conn, onMessage MessageFunc, onClose CloseFunc) *Connection {
	return &Connection{
		ws:        ws,
		onMessage: onMessage,
		onClose:   onClose,
	}
}

func (c *Connection) Read() (err error) {
	for {
		var packet dto.Packet
		err = websocket.JSON.Receive(c.ws, &packet)
		if err != nil {
			go c.onClose(c, err)
			return
		}

		go c.onMessage(&packet)
	}
}

func (c *Connection) Write(packet *dto.Packet) (err error) {
	err = websocket.JSON.Send(c.ws, packet)
	return
}
