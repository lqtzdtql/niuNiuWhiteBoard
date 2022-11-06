package models

import "time"

const (
	PermissionUser = "user"
	PermissionHost = "host"
)

type Participant struct {
	ID          int64     `json:"id"  xorm:"id pk autoincr comment('主键') BIGINT(20)"`
	Name        string    `json:"name" xorm:"'name' not null default '' comment('参会者用户名') VARCHAR(50)"`
	UserUUID    string    `json:"user_uuid" xorm:"'user_uuid' not null comment(参会唯一标识符') index VARCHAR(128)"`
	RoomUUID    string    `json:"room_uuid" xorm:"'room_uuid' not null  comment(参会所在房间标识符') index VARCHAR(128)"`
	Permission  string    `json:"permission" xorm:"permission not null default '' comment('用户权限') VARCHAR(20)"`
	CreatedTime time.Time `json:"created_time" xorm:"'created_time' not null default 'CURRENT_TIMESTAMP' comment('进房时间') TIMESTAMP"`
	UpdatedTime time.Time `json:"updated_time" xorm:"'updated_time' updated not null default 'CURRENT_TIMESTAMP' comment('修改时间') TIMESTAMP"`
	DeletedTime time.Time `json:"deleted_time" xorm:"'deleted_time' deleted"`
}

type ParticipantRow struct {
	Name       string `json:"name" xorm:"'name' not null default '' comment('参会者用户名') VARCHAR(50)"`
	UserUUID   string `json:"user_uuid" xorm:"'user_uuid' not null comment(参会唯一标识符') index VARCHAR(128)"`
	Permission string `json:"permission" xorm:"permission not null default '' comment('用户权限') VARCHAR(20)"`
}
