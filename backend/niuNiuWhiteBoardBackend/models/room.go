package models

const (
	RoomStateActive = "active"
	RoomStateStored = "stored"
)

type Room struct {
	ID           int64         `json:"id" xorm:"pk autoincr comment('主键') BIGINT(20)"`
	UUID         string        `json:"uuid" xorm:"not null unique 'uuid' comment('房间唯一标识符') index VARCHAR(128)"` // 房间对外的 UUID，同时将作为 RTC Room Name
	Name         string        `json:"name" xorm:"not null default '' comment('房间名称') VARCHAR(50)"`
	HostID       string        `json:"host" xorm:"not null default '' comment('主持人名称') VARCHAR(50)"`
	CreatedTime  int64         `json:"created_time" xorm:"not null default 0 comment('创建时间') index INT(10)"`
	UpdatedTime  int64         `json:"updated_time" xorm:"not null default 0 comment('修改时间') index INT(10)"`
	State        string        `json:"state "` // 房间状态: active, stored
	Participants []Participant `json:"participants"`
}

const (
	RoleUser = "user"
	RoleHost = "host"
)

type Participant struct {
	ID        int64  `json:"id,omitempty"      `
	UserID    int64  `json:"user_id,omitempty" `
	RoomID    int64  `json:"room_id,omitempty" `
	Role      string `json:"role"`
	UUID      string `json:"uuid"`
	Name      string `json:"name"`
	CreatedAt int64  `json:"created_at"`
	UpdatedAt int64  `json:"updated_at"`
}
