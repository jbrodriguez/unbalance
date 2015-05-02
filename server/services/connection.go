package services

import (
	"apertoire.net/unbalance/server/dto"
	"github.com/gorilla/websocket"
	"github.com/jbrodriguez/mlog"
	"github.com/jbrodriguez/pubsub"
)

type Connection struct {
	id   string
	ws   *websocket.Conn
	send chan *dto.MessageOut
	hub  *Socket
}

func (c *Connection) writer() {
	for message := range c.send {
		err := c.ws.WriteJSON(message)
		if err != nil {
			mlog.Error(err)
			c.hub.unregister <- c
			break
		}
	}
	c.ws.Close()
}

func (c *Connection) reader() {
	for {
		var msgIn dto.MessageIn
		err := c.ws.ReadJSON(&msgIn)
		if err != nil {
			mlog.Info("error is: %s", err.Error())
			break
		}

		// if DEBUG {
		mlog.Info("client type is: %s", msgIn)
		// }

		//		c.client = msgIn

		switch msgIn.Topic {
		case "storage:move":
			msg := &pubsub.Message{}
			c.hub.bus.Pub(msg, "cmd.storageMove")

		case "storage:update":
			msg := &pubsub.Message{}
			c.hub.bus.Pub(msg, "cmd.storageUpdate")

		default:
			mlog.Info("Unexpected Topic: " + msgIn.Topic)
		}
	}

	c.ws.Close()
}
