package models

import (
	"time"
)

// 连接建立时，创建一条participant记录。连接断开时，删除该participant
// CurrentWhiteBoard用于区分该参与者当前的白板。
// 参与者创建或者切换白板后，更新数据库中参与者的白板uuid。
// 查询drawing，如果绘图信息是当前白板的，那就发送给对应的client，否则不发送。

const (
	PermissionUser = "user"
	PermissionHost = "host"
)

type Participant struct {
	ID                int64     `json:"id"  xorm:"id pk autoincr BIGINT(20)"`
	UUID              string    `json:"uuid" xorm:"'uuid' not null default '' VARCHAR(128)"`
	Name              string    `json:"name" xorm:"'name' not null default '' VARCHAR(50)"`
	RoomName          string    `json:"room_name" xorm:"'room_name' not null  default '' VARCHAR(128)"`
	Permission        string    `json:"permission" xorm:"'permission' not null default '' VARCHAR(20)"`
	CurrentWhiteBoard string    `json:"current_white_board"  xorm:"'current_white_board' not null  default '' VARCHAR(128)"`
	CreatedTime       time.Time `json:"created_time" xorm:"'created_time' not null default CURRENT_TIMESTAMP"`
	UpdatedTime       time.Time `json:"updated_time" xorm:"'updated_time' updated not null default CURRENT_TIMESTAMP"`
	DeletedTime       time.Time `json:"deleted_time" xorm:"'deleted_time' datetime deleted"`
}
