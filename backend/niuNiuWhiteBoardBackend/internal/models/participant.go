package models

import "time"

const (
	PermissionUser = "user"
	PermissionHost = "host"
)

type Participant struct {
	ID          int64     `json:"id"  xorm:"id pk autoincr BIGINT(20)"`
	Name        string    `json:"name" xorm:"'name' not null default '' VARCHAR(50)"`
	UserUUID    string    `json:"user_uuid" xorm:"'user_uuid' not null default '' VARCHAR(128)"`
	RoomUUID    string    `json:"room_uuid" xorm:"'room_uuid' not null  default '' VARCHAR(128)"`
	Permission  string    `json:"permission" xorm:"'permission' not null default '' VARCHAR(20)"`
	CreatedTime time.Time `json:"created_time" xorm:"'created_time' not null default CURRENT_TIMESTAMP"`
	UpdatedTime time.Time `json:"updated_time" xorm:"'updated_time' updated not null default CURRENT_TIMESTAMP"`
	DeletedTime time.Time `json:"deleted_time" xorm:"'deleted_time' datetime deleted"`
}

type ParticipantRow struct {
	Name       string `json:"name" xorm:"'name' not null default '' VARCHAR(50)"`
	UserUUID   string `json:"user_uuid" xorm:"'user_uuid' not null default '' VARCHAR(128)"`
	Permission string `json:"permission" xorm:"'permission' not null default '' VARCHAR(20)"`
}
