package service

import (
	"net/http"
	"niuNiuSDKBackend/common/database"
	"niuNiuSDKBackend/common/log"
	"niuNiuSDKBackend/internal/jwt"
	"niuNiuSDKBackend/internal/models"
	"niuNiuSDKBackend/internal/server"
	"niuNiuSDKBackend/secretkey"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var upGrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func RunSocket(c *gin.Context) {
	token := c.Query("token")
	clientClaims, err := jwt.ParseToken(token)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "token解析失败", "code": 501})
		log.Logger.Error("token parse failed", log.Any("token parse failed", err.Error()))
		return
	}
	sk := secretkey.SecretKey{
		SK: clientClaims.SK,
	}
	has, err := database.MEngine.Table(secretkey.SKTABLE).Exist(&sk)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "服务器错误", "code": 501})
		log.Logger.Error("database error", log.Any("database error", err.Error()))
		return
	}
	if !has {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "sk错误", "code": 401})
		log.Logger.Warn("sk not exist", log.Any("sk not exist", sk))
		return
	}

	participant := models.Participant{}
	_, err = database.MEngine.Table(models.ParticipantTable).Where("name=?", clientClaims.UserName).Get(&participant)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "服务器错误", "code": 501})
		log.Logger.Error("database error", log.Any("database error", err.Error()))
		return
	}
	room := models.Room{}
	has, err = database.MEngine.Table(models.RoomTable).Where("name=?", clientClaims.RoomName).Get(&room)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "服务器错误", "code": 501})
		log.Logger.Error("database error", log.Any("database error", err.Error()))
		return
	}

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
