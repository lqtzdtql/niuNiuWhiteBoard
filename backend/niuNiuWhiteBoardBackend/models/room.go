package models

import (
	"time"
)

const (
	RoomStateActive = "active"
	RoomStateStored = "stored"
)

const (
	RoomTypeTeaching = "teaching_room" //教学房
	RoomTypePlaying  = "playing_room"  //游戏房
)

const (
	RoomTable        = "rooms"
	ParticipantTable = "participants"
)

type Room struct {
	ID           int64         `json:"id" xorm:"id pk autoincr comment('主键') BIGINT(20)"`
	UUID         string        `json:"uuid" xorm:"uuid not null unique 'uuid' comment('房间唯一标识符') index VARCHAR(128)"` // 房间对外的 UUID，同时将作为 RTC Room Name
	Name         string        `json:"name" xorm:"name not null default '' comment('房间名称') VARCHAR(50)"`
	HostID       int64         `json:"host_id"  xorm:" host_id not null default 0 comment('主持人标识符') BIGINT(20)"`
	CreatedTime  time.Time     `json:"created_time" xorm:"created_time not null default 'CURRENT_TIMESTAMP' comment('建房时间') TIMESTAMP"`
	UpdatedTime  time.Time     `json:"updated_time" xorm:"updated_time not null default 'CURRENT_TIMESTAMP' comment('修改时间') TIMESTAMP"`
	State        string        `json:"state " xorm:"-"` // 房间状态: active, stored
	Type         string        `json:"type"  xorm:"type not null default '' comment('房间类型') index VARCHAR(20)"`
	MySelf       Participant   `json:"myself" xorm:"-"`
	Participants []Participant `json:"participants" xorm:"-"`
}

type RoomRaw struct {
	ID          int64     `json:"id" xorm:"id pk autoincr comment('主键') BIGINT(20)"`
	UUID        string    `json:"uuid" xorm:"uuid not null unique 'uuid' comment('房间唯一标识符') index VARCHAR(128)"` // 房间对外的 UUID，同时将作为 RTC Room Name
	Name        string    `json:"name" xorm:"name not null default '' comment('房间名称') VARCHAR(50)"`
	HostID      int64     `json:"host_id"  xorm:" host_id not null default 0 comment('主持人标识符') BIGINT(20)"`
	CreatedTime time.Time `json:"created_time" xorm:"created_time not null default 'CURRENT_TIMESTAMP' comment('创建时间') TIMESTAMP"`
	UpdatedTime time.Time `json:"updated_time" xorm:"updated_time not null default 'CURRENT_TIMESTAMP' comment('修改时间') TIMESTAMP"`
	State       string    `json:"state " xorm:"-"` // 房间状态: active, stored
	Type        string    `json:"type"  xorm:"type not null default '' comment('房间类型') index VARCHAR(20)"`
}
