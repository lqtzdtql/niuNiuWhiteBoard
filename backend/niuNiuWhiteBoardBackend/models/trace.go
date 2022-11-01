package models

type Trace struct {
	CreatedTime int64 `json:"created_time" xorm:"not null default 0 comment('注册时间') index INT(10)"`
	ID          int64 `json:"id" xorm:"pk autoincr comment('主键') BIGINT(20)"`
	Ip          int   `json:"ip" xorm:"not null comment('ip') INT(10)"`
	Uid         int64 `json:"uid" xorm:"not null default 0 comment('用户主键') index(UT) BIGINT(20)"`
}
