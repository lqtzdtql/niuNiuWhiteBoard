package models

import "time"

const (
	UserStateOnline  = "online"
	UserStateOffline = "offline"
)

const UsersTable = "users"

type User struct {
	ID          int64     `json:"id"  xorm:"id pk autoincr comment('主键') BIGINT(20)"`
	UUID        string    `json:"uuid" xorm:"uuid not null unique 'uuid' comment('用户唯一标识符') index VARCHAR(128)"`
	Name        string    `json:"name" xorm:"name not null default '' comment('用户名') VARCHAR(50)"`
	Mobile      string    `json:"mobile" xorm:"mobile not null default '' comment('手机号') VARCHAR(20)"`
	Passwd      string    `json:"passwd" xorm:"passwd not null comment('密码') VARCHAR(50)"`
	UserState   string    `json:"user_state" xorm:"user_state not null default '' comment('用户状态') VARCHAR(20)"`
	CreatedTime time.Time `json:"created_time" xorm:"created_time not null default 'CURRENT_TIMESTAMP' comment('注册时间') TIMESTAMP"`
	UpdatedTime time.Time `json:"updated_time" xorm:"'updated_time' updated not null default 'CURRENT_TIMESTAMP' comment('修改时间') TIMESTAMP"`
	DeletedTime time.Time `json:"deleted" xorm:"'deleted_time' deleted"`
}

type UserRow struct {
	ID     int64  `json:"id" xorm:"id pk autoincr comment('主键') BIGINT(20)"`
	UUID   string `json:"uuid" xorm:"uuid not null unique 'uuid' comment('用户唯一标识符') index VARCHAR(128)"`
	Name   string `json:"name" xorm:"name not null default '' comment('用户名') VARCHAR(50)"`
	Mobile string `json:"mobile" xorm:"mobile not null default '' comment('手机号') VARCHAR(20)"`
}
