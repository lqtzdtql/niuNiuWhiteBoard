package models

const (
	UserStateOnline  = "online"
	UserStateOffline = "offline"
)

type User struct {
	ID          int64  `json:"id"  xorm:"id pk autoincr comment('主键') BIGINT(20)"`
	UUID        string `json:"uuid" xorm:"uuid not null unique 'uuid' comment('用户唯一标识符') index VARCHAR(128)"`
	Name        string `json:"name" xorm:"name not null default '' comment('用户名') VARCHAR(50)"`
	Mobile      string `json:"mobile" xorm:"mobile not null default '' comment('手机号') VARCHAR(20)"`
	Passwd      string `json:"passwd" xorm:"passwd not null comment('密码') VARCHAR(50)"`
	UserState   string `json:"user_state" xorm:"user_state not null default '' comment('用户状态') VARCHAR(20)"`
	CreatedTime int64  `json:"created_time" xorm:"created_time not null default 0 comment('创建时间') INT(10)"`
	UpdatedTime int64  `json:"updated_time" xorm:"updated_time not null default 0 comment('修改时间') INT(10)"`
}

type UserRow struct {
	Id     int64  `json:"id" xorm:"id pk autoincr comment('主键') BIGINT(20)"`
	UUID   string `json:"uuid" xorm:"uuid not null unique 'uuid' comment('用户唯一标识符') index VARCHAR(128)"`
	Name   string `json:"name" xorm:"name not null default '' comment('用户名') VARCHAR(50)"`
	Mobile string `json:"mobile" xorm:"mobile not null default '' comment('手机号') VARCHAR(20)"`
}

const UsersTable = "users"
