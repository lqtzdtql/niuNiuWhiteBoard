package service

import (
	"niuNiuSDKBackend/internal/models"
)

type roomService struct {
}

var RoomService = new(roomService)

func (g *roomService) GetUserIdByRoomUuid(groupUuid string) []models.User {
	var users []models.User
	return users
}
