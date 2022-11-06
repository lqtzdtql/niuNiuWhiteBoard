package service

import (
	"niuNiuSDKBackend/common/database"
	"niuNiuSDKBackend/internal/models"
)

type roomService struct {
}

var RoomService = new(roomService)

func (g *roomService) GetUserIdByRoomUuid(groupUuid string) []models.User {
	var room models.Room
	db := database.GetDB()
	db.First(&room, "uuid = ?", groupUuid)
	if room.ID <= 0 {
		return nil
	}

	var users []models.User
	db.Raw("SELECT u.uuid, u.avatar, u.username FROM `rooms` AS g JOIN group_members AS gm ON gm.room_id = g.id JOIN users AS u ON u.id = gm.user_id WHERE g.id = ?",
		room.ID).Scan(&users)
	return users
}
