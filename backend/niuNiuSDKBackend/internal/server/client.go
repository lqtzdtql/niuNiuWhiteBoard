package server

import (
	"encoding/json"
	"time"

	"niuNiuSDKBackend/common/log"
	"niuNiuSDKBackend/internal/models"

	"github.com/gorilla/websocket"
)

type Client struct {
	Conn          *websocket.Conn //一个账号，一个连接
	UUID          string          //参与者uuid
	RoomUUID      string          //参与者所在房间uuid
	Send          chan []byte
	HeartbeatTime int64 // 前一次心跳时间
}

func (c *Client) Read() {
	defer func() {
		user := ExitRoom(c)
		leave := &models.LeaveEnterRoomRes{
			ContentType: models.LEAVE_ROOM,
			UserName:    user.Name,
		}
		message, _ := json.Marshal(leave)
		MyServer.mutex.Lock()
		for _, c := range MyServer.Clients {
			if c.UUID != c.UUID {
				c.Send <- message
			}
		}
		MyServer.mutex.Unlock()
		c.Conn.Close()
		delete(MyServer.Clients, c.UUID)
	}()

	for {
		_, message, err := c.Conn.ReadMessage()
		if err != nil {
			log.Logger.Error("client read message error", log.Any("client read message error", err.Error()))
			MyServer.UnRegister <- c
			c.Conn.Close()
			break
		}
		msg := &models.Message{}
		json.Unmarshal(message, msg)
		log.Logger.Debug("receive message", log.Any("receive message", msg))
		// pong
		if msg.ContentType == models.HEAT_BEAT {
			c.HeartbeatTime = time.Now().Unix()
			pong := &models.HeatBeatRes{
				ContentType: models.HEAT_BEAT,
			}
			respong, _ := json.Marshal(pong)
			c.Conn.WriteMessage(websocket.TextMessage, respong)
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
		c.Conn.WriteMessage(websocket.TextMessage, message)
	}
}
