package service

import (
	"niuNiuSDKBackend/common/database"
	"niuNiuSDKBackend/common/log"
	"niuNiuSDKBackend/internal/models"
)

type messageService struct {
}

var MessageService = new(messageService)

const NULL_ID int32 = 0

func (m *messageService) SaveMessage(message models.Message) {
	db := database.GetDB()
	var fromUser models.User
	db.Find(&fromUser, "uuid = ?", message.From)
	if NULL_ID == fromUser.Id {
		log.Logger.Error("SaveMessage not find from user", log.Any("SaveMessage not find from user", fromUser.Id))
		return
	}

	var toUserId int32 = 0
	var group models.Room
	db.Find(&group, "uuid = ?", message.To)
	if NULL_ID == group.ID {
		return
	}
	toUserId = group.ID

	saveMessage := models.Message{
		FromUserId:  fromUser.Id,
		ToUserId:    toUserId,
		Content:     message.Content,
		ContentType: message.ContentType,
		MessageType: message.MessageType,
	}
	db.Save(&saveMessage)
}
