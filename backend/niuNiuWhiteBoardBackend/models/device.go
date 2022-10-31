package models

type Device struct {
	Client      string `json:"client" xorm:"not null default '' comment('客户端') VARCHAR(50)"`
	CreatedTime int64  `json:"created_time" xorm:"not null default 0 comment('注册时间') index INT(10)"`
	ID          int64  `json:"id" xorm:"pk autoincr comment('主键') BIGINT(20)"`
	Ip          int    `json:"ip" xorm:"not null default 0 comment('ip地址') INT(10)"`
	Uid         int64  `json:"uid" xorm:"not null default 0 comment('用户主键') index BIGINT(20)"`
}
