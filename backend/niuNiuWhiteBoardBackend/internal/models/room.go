package models

import (
	"time"
)

const (
	RoomTypeTeaching = "教学房" //教学房
	RoomTypePlaying  = "游戏房" //游戏房
)

const (
	RoomTable        = "rooms"
	ParticipantTable = "participants"
)

type Room struct {
	ID           int64         `json:"id" xorm:"'id' pk autoincr BIGINT(20)"`
	UUID         string        `json:"uuid" xorm:"'uuid' not null default '' VARCHAR(128)"`
	Name         string        `json:"name" xorm:"'name' not null default '' VARCHAR(50)"`
	HostUUID     string        `json:"host_uuid" xorm:"'host_uuid' not null default '' VARCHAR(128)"`
	HostName     string        `json:"host_name" xorm:"'host_name' not null default '' VARCHAR(50)"`
	CreatedTime  time.Time     `json:"created_time" xorm:"'created_time'"`
	UpdatedTime  time.Time     `json:"updated_time" xorm:"'updated_time' updated"`
	DeletedTime  time.Time     `json:"deleted_time" xorm:"'deleted_time' datetime deleted"`
	Type         string        `json:"type"  xorm:"'type' not null default ''"`
	MySelf       Participant   `json:"myself" xorm:"-"`
	Participants []Participant `json:"participants" xorm:"-"`
}

type RoomRaw struct {
	UUID         string           `json:"uuid" xorm:"'uuid' not null default '' VARCHAR(128)"`
	Name         string           `json:"name" xorm:"name not null default '' VARCHAR(50)"`
	HostUUID     string           `json:"host_uuid" xorm:"'host_uuid' not null default '' VARCHAR(128)"`
	HostName     string           `json:"host_name" xorm:"'host_name' not null default '' VARCHAR(50)"`
	CreatedTime  time.Time        `json:"created_time" xorm:"'created_time'  TIMESTAMP"`
	UpdatedTime  time.Time        `json:"updated_time" xorm:"'updated_time' updated"`
	Type         string           `json:"type"  xorm:"'type' not null default '' VARCHAR(20)"`
	Participants []ParticipantRow `json:"participants" xorm:"-"`
}

type RoomInfo struct {
	UUID     string `json:"uuid" xorm:"'uuid' not null default '' VARCHAR(128)"`
	Name     string `json:"name" xorm:"'name' not null default '' VARCHAR(50)"`
	HostName string `json:"host_name" xorm:"'host_name' not null default '' VARCHAR(50)"`
	Type     string `json:"type"  xorm:"'type' not null default '' VARCHAR(20)"`
}

type RoomNameType struct {
	Name string `json:"name"`
	Type string `json:"type"`
}
