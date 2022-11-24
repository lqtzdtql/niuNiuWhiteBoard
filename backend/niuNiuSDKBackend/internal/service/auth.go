package service

import (
	"net/http"
	"time"

	"niuNiuSDKBackend/common/database"
	"niuNiuSDKBackend/common/log"
	"niuNiuSDKBackend/internal/jwt"
	"niuNiuSDKBackend/internal/models"
	"niuNiuSDKBackend/secretkey"

	"github.com/gin-gonic/gin"
	"github.com/oklog/ulid/v2"
)

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

	participant := models.Participant{}
	has, err = database.MEngine.Table(models.ParticipantTable).Where("name=? and room_name=?", clientClaims.UserName, clientClaims.RoomName).Get(&participant)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "服务器错误", "code": 501})
		log.Logger.Error("database error", log.Any("database error", err.Error()))
		return
	}
	if has {
		c.JSON(http.StatusOK, gin.H{
			"message":   "auth success",
			"user_uuid": participant.UUID,
			"user_name": participant.Name,
			"room_uuid": participant.RoomUUID,
			"code":      200,
		})
		return
	}

	participant.Name = clientClaims.UserName
	participant.RoomName = clientClaims.RoomName
	participant.Permission = clientClaims.Permission
	participant.UUID = ulid.Make().String()

	room := models.Room{}
	have, err := database.MEngine.Table(models.RoomTable).Where("name=?", clientClaims.RoomName).Get(&room)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "服务器错误", "code": 501})
		log.Logger.Error("database error", log.Any("database error", err.Error()))
		return
	}
	if !have {
		//房间不存在，主持人权限的可以建房
		if participant.Permission == models.PermissionHost {
			room.Name = clientClaims.RoomName
			room.CreatedTime = time.Now()
			room.UpdatedTime = time.Now()
			room.UUID = ulid.Make().String()
			room.HostUUID = participant.UUID
			database.MEngine.Table(models.RoomTable).Insert(&room)
		} else {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "不是主持人，无法建房", "code": 401})
			log.Logger.Warn("can not build room", log.Any("can not build room", room.UUID))
			return
		}
	}
	participant.RoomUUID = room.UUID
	participant.CreatedTime = time.Now()
	participant.UpdatedTime = time.Now()
	if _, err := database.MEngine.Table(models.ParticipantTable).Insert(&participant); err != nil {
		log.Logger.Error("participant enter failed", log.Any("participant enter failed", err.Error()))
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "服务器错误", "code": 501})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":   "auth success",
		"user_uuid": participant.UUID,
		"user_name": participant.Name,
		"room_uuid": room.UUID,
		"code":      200,
	})
	return
}
