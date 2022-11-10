package models

import "time"

type WhiteBoard struct {
	ID          int64     `json:"id" xorm:"'id' pk autoincr BIGINT(20)"`
	UUID        string    `json:"uuid" xorm:"'uuid' not null default '' VARCHAR(128)"`
	Name        string    `json:"name" xorm:"'name' not null default '' VARCHAR(50)"`
	CreatedTime time.Time `json:"created_time" xorm:"'created_time'"`
	UpdatedTime time.Time `json:"updated_time" xorm:"'updated_time' updated"`
	DeletedTime time.Time `json:"deleted_time" xorm:"'deleted_time' datetime deleted"`
}
