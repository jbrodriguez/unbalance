package services

import (
	"apertoire.net/unbalance/server/dto"
	"github.com/gorilla/websocket"
	"github.com/jbrodriguez/mlog"
	"github.com/jbrodriguez/pubsub"
	"time"
)

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer.
	maxMessageSize = 512
)

type Connection struct {
	id   string
	ws   *websocket.Conn
	send chan *dto.MessageOut
	hub  *Socket
}

// write writes a message with the given message type and payload.
func (c *Connection) write(mt int, payload []byte) error {
	c.ws.SetWriteDeadline(time.Now().Add(writeWait))
	return c.ws.WriteMessage(mt, payload)
}

func (c *Connection) writer() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.ws.Close()
	}()

	// mlog.Info("before write loop")

	for {
		select {
		case message, ok := <-c.send:
			if !ok {
				c.write(websocket.CloseMessage, []byte{})
				mlog.Warning("Closing socket ...")
				return
			}

			c.ws.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.ws.WriteJSON(message); err != nil {
				mlog.Error(err)
				return
			}
		case <-ticker.C:
			if err := c.write(websocket.PingMessage, []byte{}); err != nil {
				mlog.Info("error with ping: %s", err.Error())
				return
			}
		}
	}

}

func (c *Connection) reader() {
	defer func() {
		c.hub.unregister <- c
		c.ws.Close()
	}()

	c.ws.SetReadLimit(maxMessageSize)
	c.ws.SetReadDeadline(time.Now().Add(pongWait))
	c.ws.SetPongHandler(func(string) error {
		c.ws.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

	// mlog.Info("before read loop")

	for {
		var msgIn dto.MessageIn
		err := c.ws.ReadJSON(&msgIn)
		if err != nil {
			mlog.Info("error reading socket: %s", err.Error())
			break
		}

		// if DEBUG {
		mlog.Info("client type is: %s", msgIn)
		// }

		//		c.client = msgIn

		msg := &pubsub.Message{}
		c.hub.bus.Pub(msg, msgIn.Topic)

		// switch msgIn.Topic {
		// case "storage:move":
		// 	msg := &pubsub.Message{}
		// 	c.hub.bus.Pub(msg, "cmd.storageMove")

		// case "storage:update":
		// 	msg := &pubsub.Message{}
		// 	c.hub.bus.Pub(msg, "cmd.storageUpdate")

		// default:
		// 	mlog.Info("Unexpected Topic: " + msgIn.Topic)
		// }
	}

}
