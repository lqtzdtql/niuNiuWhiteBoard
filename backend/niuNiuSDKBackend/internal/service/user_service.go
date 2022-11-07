package service

import (
	"niuNiuSDKBackend/common/database"
	"niuNiuSDKBackend/internal/models"
)

type userService struct {
}

var UserService = new(userService)

func (u *userService) GetUserDetails(uuid string) models.User {
	var queryUser *models.User
	db := database.GetDB()
	db.Select("uuid", "username").First(&queryUser, "uuid = ?", uuid)
	return *queryUser
}
