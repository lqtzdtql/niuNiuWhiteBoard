package service

import (
	"net/http"
	"niuNiuSDKBackend/common/log"
	"niuNiuSDKBackend/internal/server"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"go.uber.org/zap"
)

var upGrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func RunSocket(c *gin.Context) {
	user := c.Query("user")
	if user == "" {
		return
	}
	//TODO：鉴权逻辑

	log.Logger.Info("newUser", zap.String("newUser", user))
	ws, err := upGrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusOK, gin.H{"message": "websocket升级失败", "code": 401})
		log.Logger.Error("Upgrade failed", log.Any("Upgrade failed", err.Error()))
		return
	}

	client := &server.Client{
		Name: user,
		Conn: ws,
		Send: make(chan []byte),
	}

	server.MyServer.Register <- client
	go client.Read()
	go client.Write()
}
