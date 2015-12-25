package net

import (
	"golang.org/x/net/websocket"
	"jbrodriguez/unbalance/server/dto"
)

type MessageFunc func(message *dto.Packet)
type CloseFunc func(conn *Connection, err error)

type Connection struct {
	id        string
	ws        *websocket.Conn
	onMessage MessageFunc
	onClose   CloseFunc
}

func NewConnection(ws *websocket.Conn, onMessage MessageFunc, onClose CloseFunc) *Connection {
	return &Connection{
		ws:        ws,
		onMessage: onMessage,
		onClose:   onClose,
	}
}

func (c *Connection) Read() {
	for {
		var packet dto.Packet
		err := websocket.Message.Receive(c.ws, &packet)
		if err != nil {
			go c.onClose(c, err)
			break
		} else {
			go c.onMessage(&packet)
		}
	}
}

func (c *Connection) Write(packet *dto.Packet) (err error) {
	err = websocket.Message.Send(c.ws, packet)
	return
}
