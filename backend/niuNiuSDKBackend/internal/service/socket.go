package service

import (
	"net/http"
	"time"

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
	has, err = database.MEngine.Table(models.ParticipantTable).Where("name=? and room_name=?", clientClaims.UserName, clientClaims.RoomName).Get(&participant)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "服务器错误", "code": 501})
		log.Logger.Error("database error", log.Any("database error", err.Error()))
		return
	}
	if !has {
		database.MEngine.Table(models.RoomTable).Delete(&models.Room{Name: clientClaims.RoomName})
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "用户不存在", "code": 401})
		log.Logger.Warn("participant not exit", log.Any("participant not exit", clientClaims.UserName))
		return
	}

	room := models.Room{}
	has, err = database.MEngine.Table(models.RoomTable).Where("name=?", clientClaims.RoomName).Get(&room)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "服务器错误", "code": 501})
		log.Logger.Error("database error", log.Any("database error", err.Error()))
		return
	}
	if !has {
		database.MEngine.Table(models.ParticipantTable).Delete(&participant)
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "房间不存在", "code": 401})
		log.Logger.Warn("participant not exit", log.Any("room not exit", clientClaims.RoomName))
		return
	}

	log.Logger.Info("participant info", log.Any("participant info", participant))

	log.Logger.Info("room info", log.Any("room info", room))
	if _, ok := server.MyServer.Clients[participant.UUID]; ok {
		log.Logger.Warn("websocket has build, not build again", log.Any("websocket has build, not build again", participant.UUID))
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
		UUID:          participant.UUID,
		Conn:          ws,
		RoomUUID:      room.UUID,
		Send:          make(chan []byte),
		HeartbeatTime: time.Now().Unix(),
	}
	log.Logger.Debug("client", log.Any("client", client.UUID))

	server.MyServer.Register <- client
	go client.Read()
	log.Logger.Info("start to read", log.Any("start to read", ws.RemoteAddr()))
	go client.Write()
	log.Logger.Info("start to write", log.Any("start to write", ws.RemoteAddr()))
}
