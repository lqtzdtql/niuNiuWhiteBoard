package server

import (
	"github.com/gorilla/websocket"
	"niuNiuSDKBackend/common/kafka"
	"niuNiuSDKBackend/common/log"
	"niuNiuSDKBackend/config"
	"niuNiuSDKBackend/internal/models"
)

type Client struct {
	Conn *websocket.Conn
	Name string
	Send chan []byte
}

func (c *Client) Read() {
	defer func() {
		MyServer.Ungister <- c
		c.Conn.Close()
	}()

	for {
		c.Conn.PongHandler()
		_, message, err := c.Conn.ReadMessage()
		if err != nil {
			log.Logger.Error("client read message error", log.Any("client read message error", err.Error()))
			MyServer.Ungister <- c
			c.Conn.Close()
			break
		}

		msg := &models.Message{}

		// pong
		if msg.Type == models.HEAT_BEAT {
			pong := &models.Message{
				Content: models.PONG,
				Type:    models.HEAT_BEAT,
			}
			c.Conn.WriteJSON(pong)
		} else {
			if config.GetConfig().MsgChannelType.ChannelType == models.KAFKA {
				kafka.Send(message)
			} else {
				MyServer.Broadcast <- message
			}
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
