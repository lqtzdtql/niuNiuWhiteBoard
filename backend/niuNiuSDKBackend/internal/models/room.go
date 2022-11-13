package models

import (
	"time"
)

/*
	连接建立时，根据roomuuid查询room数据库，查看有没有该房间
	如果没有，则建立房间，加入者为host
	如果有，则加入房间，在participant表中插入。
	当有人退出时，删除该用户。
	如果退出的是房主，则广播踢人消息。房间销毁。
*/

const (
	RoomTable        = "rooms"
	ParticipantTable = "participants"
	WhiteBoardTable  = "whiteboards"
)

type Room struct {
	ID          int64     `json:"id" xorm:"'id' pk autoincr BIGINT(20)"`
	UUID        string    `json:"uuid" xorm:"'uuid' not null default '' VARCHAR(128)"`
	Name        string    `json:"name" xorm:"'name' not null default '' VARCHAR(50)"`
	HostUUID    string    `json:"host_uuid" xorm:"'host_uuid' not null default '' VARCHAR(128)"`
	CreatedTime time.Time `json:"created_time" xorm:"'created_time'"`
	UpdatedTime time.Time `json:"updated_time" xorm:"'updated_time' updated"`
	DeletedTime time.Time `json:"deleted_time" xorm:"'deleted_time' datetime deleted"`
}
