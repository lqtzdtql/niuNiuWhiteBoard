package models

import (
	"gorm.io/plugin/soft_delete"
	"time"
)

type Participant struct {
	ID        int32                 `json:"id" gorm:"primarykey"`
	UserID    int32                 `json:"userId" gorm:"index;comment:'参与者ID'"`
	RoomID    int32                 `json:"groupId" gorm:"index;comment:'房间ID'"`
	Name      string                `json:"nickname" gorm:"type:varchar(350);comment:'参与者用户名 "`
	Mute      int16                 `json:"mute" gorm:"comment:'只读'"`
	CreatedAt time.Time             `json:"createAt"`
	UpdatedAt time.Time             `json:"updatedAt"`
	DeletedAt soft_delete.DeletedAt `json:"deletedAt"`
}
