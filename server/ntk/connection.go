package ntk

import (
	"encoding/json"
	"fmt"

	"unbalance/dto"

	"github.com/gorilla/websocket"
)

// MessageFunc -
type MessageFunc func(message *dto.Packet)

// CloseFunc -
type CloseFunc func(conn *Connection, err error)

// Connection -
type Connection struct {
	// id        string
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
		var p []byte

		_, p, err = c.ws.ReadMessage()
		if err != nil {
			errm := fmt.Errorf("unable to ReadMessage: (%s)", err)
			go c.onClose(c, errm)
			return err
		}

		err = json.Unmarshal(p, &packet)

		// err = c.ws.ReadJSON(&packet)
		if err != nil {
			errm := fmt.Errorf("unable to Unmarshal: content(%s): err(%s)", p, err)
			go c.onClose(c, errm)
			return err
		}

		go c.onMessage(&packet)
	}
}

func (c *Connection) Write(packet *dto.Packet) (err error) {
	err = c.ws.WriteJSON(packet)
	return
}
