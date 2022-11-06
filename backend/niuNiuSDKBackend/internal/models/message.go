package models

import (
	"gorm.io/plugin/soft_delete"
	"time"
)

const (
	HEAT_BEAT = "heatbeat"
	PONG      = "pong"

	// 广播类消息
	MESSAGE_TYPE_BROADCAST = 1

	// TODO:消息内容类型
	SIGNALING = 1
	POINT     = 2
	OBJECT    = 5

	// 消息队列类型
	GO_CHANNEL = "gochannel"
	KAFKA      = "kafka"
)

type Message struct {
	FromUsername string                `json:"fromUsername,omitempty"`
	From         string                ` json:"from,omitempty"`
	To           string                `json:"to,omitempty"`
	FromUserId   int32                 `json:"fromUserId" gorm:"index"`
	ToUserId     int32                 `json:"toUserId" gorm:"index;comment:'发送给端的id，可为用户id或者房间id'"`
	Content      string                `json:"content,omitempty"`
	ContentType  int32                 `json:"contentType,omitempty"`
	Type         string                `json:"type,omitempty"`
	MessageType  int32                 `json:"messageType,omitempty"`
	CreatedAt    time.Time             `json:"createAt"`
	UpdatedAt    time.Time             `json:"updatedAt"`
	DeletedAt    soft_delete.DeletedAt `json:"deletedAt"`
}
