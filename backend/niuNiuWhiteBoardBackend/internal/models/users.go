package models

import "time"

const (
	UserStateOnline  = "online"
	UserStateOffline = "offline"
)

const UsersTable = "users"

type User struct {
	ID          int64      `json:"id"  xorm:"id pk autoincr BIGINT(20)"`
	UUID        string     `json:"uuid" xorm:"'uuid' not null default '' VARCHAR(128)"`
	Name        string     `json:"name" xorm:"'name' not null default '' VARCHAR(50)"`
	Mobile      string     `json:"mobile" xorm:"'mobile' not null default '' VARCHAR(20)"`
	Passwd      string     `json:"passwd" xorm:"'passwd' not null default '' VARCHAR(50)"`
	UserState   string     `json:"user_state" xorm:"'user_state' not null default '' VARCHAR(20)"`
	CreatedTime time.Time  `json:"created_time" xorm:"'created_time'"`
	UpdatedTime time.Time  `json:"updated_time" xorm:"'updated_time' updated"`
	DeletedTime *time.Time `json:"deleted" xorm:"'deleted_time' deleted"`
}

type UserRow struct {
	UUID   string `json:"uuid" xorm:"'uuid' not null default '' VARCHAR(128)"`
	Name   string `json:"name" xorm:"'name' not null default '' VARCHAR(50)"`
	Mobile string `json:"mobile" xorm:"'mobile' not null default '' VARCHAR(20)"`
}
