package service

import (
	"net/http"
	"niuNiuSDKBackend/common/log"
	"niuNiuSDKBackend/internal/models"
	"niuNiuSDKBackend/internal/server"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var upGrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func RunSocket(c *gin.Context) {
	room := c.MustGet("room").(*models.Room)
	participant := c.MustGet("participant").(*models.Participant)

	ws, err := upGrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusOK, gin.H{"message": "websocket upgrade failed", "code": 401})
		log.Logger.Error("Upgrade failed", log.Any("Upgrade failed", err.Error()))
		return
	}
	log.Logger.Info("websocket build success", log.Any("websocket build success", ws.RemoteAddr()))

	client := &server.Client{
		UUID:     participant.UUID,
		Conn:     ws,
		RoomUUID: room.UUID,
		Send:     make(chan []byte),
	}
	server.MyServer.Register <- client
	go client.Read()
	log.Logger.Info("start to read", log.Any("start to read", ws.RemoteAddr()))
	go client.Write()
	log.Logger.Info("start to write", log.Any("start to write", ws.RemoteAddr()))
}
