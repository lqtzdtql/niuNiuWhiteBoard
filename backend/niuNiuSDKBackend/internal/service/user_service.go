package service

import (
	"niuNiuSDKBackend/internal/models"
)

type userService struct {
}

var UserService = new(userService)

func (u *userService) GetUserDetails(uuid string) models.User {
	var queryUser *models.User
	return *queryUser
}
