package service

import (
	"niuNiuSDKBackend/internal/models"
)

type messageService struct {
}

var MessageService = new(messageService)

const NULL_ID int32 = 0

func (m *messageService) SaveMessage(message models.Message) {
}
