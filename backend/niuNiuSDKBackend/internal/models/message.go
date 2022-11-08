package models

import (
	"gorm.io/plugin/soft_delete"
	"time"
)

const (
	PONG = "pong"

	// 广播类消息
	MESSAGE_TYPE_BROADCAST = 1

	// TODO:消息内容类型
	SIGNALING = 1
	POINT     = 2
	OBJECT    = 3
	REPAINT   = 4
	HEAT_BEAT = 5

	// 消息队列类型
	GO_CHANNEL = "gochannel"
	KAFKA      = "kafka"
)

type Message struct {
	From        string                `json:"from,omitempty"`
	To          string                `json:"to,omitempty"`
	Content     string                `json:"content,omitempty"`
	ContentType int32                 `json:"contentType,omitempty"`
	CreatedAt   time.Time             `json:"createAt"`
	UpdatedAt   time.Time             `json:"updatedAt"`
	DeletedAt   soft_delete.DeletedAt `json:"deletedAt"`
}
