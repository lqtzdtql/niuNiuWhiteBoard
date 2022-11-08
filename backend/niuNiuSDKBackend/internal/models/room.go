package models

import (
	"gorm.io/plugin/soft_delete"
	"time"
)

type Room struct {
	ID        int32                 `json:"id" gorm:"primarykey"`
	Uuid      string                `json:"uuid" gorm:"type:varchar(150);not null;unique_index:idx_uuid;comment:'uuid'"`
	Name      string                `json:"name" gorm:"type:varchar(150);comment:'房间名称"`
	CreatedAt time.Time             `json:"createAt"`
	UpdatedAt time.Time             `json:"updatedAt"`
	DeletedAt soft_delete.DeletedAt `json:"deletedAt"`
	HostID    int32                 `json:"userId" gorm:"index;comment:'主持人标识'"`
}
