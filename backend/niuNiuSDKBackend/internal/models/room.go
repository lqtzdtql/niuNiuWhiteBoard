package models

import (
	"gorm.io/plugin/soft_delete"
	"time"
)

// 连接建立时，根据roomuuid查询room数据库，查看有没有该房间
// 如果没有，则建立房间，加入者为host
// 如果有，则加入房间，在participant表中插入。
// 当有人退出时，删除该用户。
// 如果退出的是房主，则广播踢人消息。房间销毁。

type Room struct {
	ID        int32                 `json:"id" gorm:"primarykey"`
	Uuid      string                `json:"uuid" gorm:"type:varchar(150);not null;unique_index:idx_uuid;comment:'uuid'"`
	Name      string                `json:"name" gorm:"type:varchar(150);comment:'房间名称"`
	HostID    int32                 `json:"userId" gorm:"index;comment:'主持人标识'"`
	CreatedAt time.Time             `json:"createAt"`
	UpdatedAt time.Time             `json:"updatedAt"`
	DeletedAt soft_delete.DeletedAt `json:"deletedAt"`
}
