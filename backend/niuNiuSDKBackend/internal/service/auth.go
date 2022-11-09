package service

import (
	"github.com/oklog/ulid/v2"
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

func Auth(clientClaims *jwt.ClientClaims) (bool, error) {
	sk := secretkey.SecretKey{
		SK: clientClaims.SK,
	}
	has, err := database.MEngine.Table(secretkey.SKTABLE).Exist(&sk)
	if err != nil {
		return false, err
	}
	if !has {
		return false, err
	}
	participant := models.Participant{
		Name: clientClaims.UserName,
	}
	has, err = database.MEngine.Table(models.ParticipantTable).Exist(&participant)
	if err != nil {
		return false, err
	}
	if has {
		return false, nil
	}
	participant.RoomName = clientClaims.RoomName
	participant.Permission = clientClaims.Permission
	participant.UUID = ulid.Make().String()
	room := models.Room{
		Name: clientClaims.RoomName,
	}
	has, err = database.MEngine.Table(models.RoomTable).Exist(&room)
	if err != nil {
		return false, err
	}
	if !has {
		if participant.Permission == models.PermissionHost {
			room.UUID = ulid.Make().String()
			room.HostUUID = participant.UUID
			if _, err := database.MEngine.Table(models.RoomTable).Insert(room); err != nil {
				log.Logger.Error("build room failed", log.Any("build room failed", err.Error()))
				return false, err
			}
		}
	}
	if _, err := database.MEngine.Table(models.ParticipantTable).Insert(&participant); err != nil {
		log.Logger.Error("participant enter failed", log.Any("participant enter failed", err.Error()))
		return false, err
	}
	return true, nil
}
