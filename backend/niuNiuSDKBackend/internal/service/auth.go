package service

import (
	"github.com/gin-gonic/gin"
	"github.com/oklog/ulid/v2"
	"net/http"
	"niuNiuSDKBackend/common/database"
	"niuNiuSDKBackend/common/log"
	"niuNiuSDKBackend/internal/jwt"
	"niuNiuSDKBackend/internal/models"
	"niuNiuSDKBackend/secretkey"
)

// 先判断有没有该用户uuid
// 如果有该用户，返回错误
// 如果没有该用户
// --判断有没有该房间
// ----如果是房主，有则进入，没有则建立room，同时participants中加入该用户。
// ----如果是普通用户,没有则返回错误，有则participants中加入该用户。

func Auth(c *gin.Context) {
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
	participant := models.Participant{
		Name: clientClaims.UserName,
	}
	has, err = database.MEngine.Table(models.ParticipantTable).Exist(&participant)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "服务器错误", "code": 501})
		log.Logger.Error("database error", log.Any("database error", err.Error()))
		return
	}
	if has {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "user has in the room", "code": 401})
		log.Logger.Warn("user has in the room", log.Any("user has in the room", participant))
		return
	}
	participant.Permission = clientClaims.Permission
	participant.UUID = ulid.Make().String()
	room := models.Room{
		Name: clientClaims.RoomName,
	}
	has, err = database.MEngine.Table(models.RoomTable).Exist(&room)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "服务器错误", "code": 501})
		log.Logger.Error("database error", log.Any("database error", err.Error()))
		return
	}
	if !has {
		if participant.Permission == models.PermissionHost {
			room.UUID = ulid.Make().String()
			room.HostUUID = participant.UUID
			if _, err := database.MEngine.Table(models.RoomTable).Insert(&room); err != nil {
				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "服务器错误", "code": 501})
				log.Logger.Error("build room failed", log.Any("build room failed", err.Error()))
				return
			}
		}
	}
	participant.RoomUUID = room.UUID
	if _, err := database.MEngine.Table(models.ParticipantTable).Insert(&participant); err != nil {
		log.Logger.Error("participant enter failed", log.Any("participant enter failed", err.Error()))
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "服务器错误", "code": 501})
		return
	}
	c.Set("room", &room)
	c.Set("participant", &participant)
	c.JSON(http.StatusOK, gin.H{
		"message":   "auth success",
		"user_uuid": participant.UUID,
		"room_uuid": room.UUID,
		"code":      200,
	})
	c.Next()
	return
}
