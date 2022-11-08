package server

import (
	"github.com/gorilla/websocket"
	"niuNiuSDKBackend/common/log"
	"niuNiuSDKBackend/internal/models"
)

type Client struct {
	Conn *websocket.Conn
	Name string
	Send chan []byte
}

func (c *Client) Read() {
	defer func() {
		MyServer.UnRegister <- c
		c.Conn.Close()
	}()

	for {
		c.Conn.PongHandler()
		_, message, err := c.Conn.ReadMessage()
		if err != nil {
			log.Logger.Error("client read message error", log.Any("client read message error", err.Error()))
			MyServer.UnRegister <- c
			c.Conn.Close()
			break
		}

		msg := &models.Message{}

		// pong
		if msg.ContentType == models.HEAT_BEAT {
			pong := &models.Message{
				Content:     models.PONG,
				ContentType: models.HEAT_BEAT,
			}
			c.Conn.WriteJSON(pong)
		} else {
			MyServer.Broadcast <- message
		}
	}
}

func (c *Client) Write() {
	defer func() {
		c.Conn.Close()
	}()

	for message := range c.Send {
		c.Conn.WriteJSON(message)
	}
}
