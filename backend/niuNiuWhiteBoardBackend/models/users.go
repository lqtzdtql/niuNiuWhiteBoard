package models

import (
	"niuNiuWhiteBoardBackend/config"
)

const (
	UserStateOnline  = "online"
	UserStateOffline = "offline"
)

type Users struct {
	ID          int64  `json:"id"  xorm:"pk autoincr comment('主键') BIGINT(20)"`
	UUID        string `json:"uuid" xorm:"not null unique 'uuid' comment('用户唯一标识符') index VARCHAR(128)"`
	Name        string `json:"name" xorm:"not null default '' comment('用户名') VARCHAR(50)"`
	Mobile      string `json:"mobile" xorm:"not null default '' comment('手机号') VARCHAR(20)"`
	Passwd      string `json:"passwd" xorm:"not null comment('密码') VARCHAR(50)"`
	CreatedTime int64  `json:"created_time" xorm:"not null default 0 comment('创建时间') INT(10)"`
	UpdatedTime int64  `json:"updated_time" xorm:"not null default 0 comment('修改时间') INT(10)"`
}
type UserRow struct {
	ID     int64  `json:"id" xorm:"pk autoincr comment('主键') BIGINT(20)"`
	Name   string `json:"name" xorm:"not null default '' comment('用户名') VARCHAR(50)"`
	Mobile string `json:"mobile" xorm:"not null default '' comment('手机号') VARCHAR(20)"`
}

var usersTable = "users"

func (u *Users) GetRow() bool {
	has, err := conf.MEngine.Get(u)
	if err == nil && has {
		return true
	}
	return false
}
func (u *Users) GetAll() ([]Users, error) {
	var users []Users
	err := conf.MEngine.Find(&users)
	return users, err
}

func (u *Users) Add(trace *Trace, device *Device) (int64, error) {
	session := conf.MEngine.NewSession()
	defer session.Close()
	// add Begin() before any action
	if err := session.Begin(); err != nil {
		return 0, err
	}
	_, err := session.Insert(u)
	if err != nil {
		return 0, err
	}

	trace.Uid = u.ID
	_, err = session.Insert(trace)
	if err != nil {
		return 0, err
	}
	device.Uid = u.ID
	_, err = session.Insert(device)
	if err != nil {
		return 0, err
	}
	return u.ID, session.Commit()
}
func IsExistsMobile(mobile string) bool {
	model := Users{Mobile: mobile}
	return model.GetRow()
}

func (u *Users) GetRowById() (UserRow, error) {
	var userRow UserRow
	_, err := conf.MEngine.Table(usersTable).Where("id=?", u.ID).Get(&userRow)
	return userRow, err
}
