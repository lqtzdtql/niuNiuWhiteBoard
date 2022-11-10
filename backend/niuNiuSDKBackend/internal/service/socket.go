package service

import (
	"net/http"
	"niuNiuSDKBackend/common/log"
	"niuNiuSDKBackend/internal/jwt"
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
	token := c.Query("Access-Token")
	clientClaims, err := jwt.ParseToken(token)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "token解析失败", "code": 501})
		log.Logger.Error("token parse failed", log.Any("token parse failed", err.Error()))
		return
	}
	//TODO：鉴权逻辑
	success, err := Auth(clientClaims)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "服务器错误", "code": 501})
		log.Logger.Error("server error", log.Any("server error", err.Error()))
		return
	}
	if !success {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "鉴权失败", "code": 401})
		log.Logger.Warn("auth failed", log.Any("auth failed", err.Error()))
		return
	}

	log.Logger.Info("newParticipant", zap.String("newParticipant", clientClaims.UserName))
	ws, err := upGrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusOK, gin.H{"message": "websocket升级失败", "code": 401})
		log.Logger.Error("Upgrade failed", log.Any("Upgrade failed", err.Error()))
		return
	}

	client := &server.Client{
		UUID: clientClaims.UserName,
		Conn: ws,
		Send: make(chan []byte),
	}

	server.MyServer.Register <- client
	go client.Read()
	go client.Write()
}
