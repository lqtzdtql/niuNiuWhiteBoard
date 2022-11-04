package models

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/go-xorm/xorm"
	"github.com/oklog/ulid/v2"
	"github.com/qiniu/go-sdk/v7/auth"
	"github.com/qiniu/go-sdk/v7/rtc"
	"log"
	"net/http"
	conf "niuNiuWhiteBoardBackend/config"
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

const (
	PermissionUser = "user"
	PermissionHost = "host"
)

const (
	RoomTable        = "rooms"
	ParticipantTable = "participants"
)

type Participant struct {
	ID          int64     `json:"id"  xorm:"id pk autoincr comment('主键') BIGINT(20)"`
	Name        string    `json:"name" xorm:"'name' not null default '' comment('参会者用户名') VARCHAR(50)"`
	UserUUID    string    `json:"user_uuid" xorm:"'user_uuid' not null comment(参会唯一标识符') index VARCHAR(128)"`
	RoomUUID    string    `json:"room_uuid" xorm:"'room_uuid' not null  comment(参会所在房间标识符') index VARCHAR(128)"`
	Permission  string    `json:"permission" xorm:"permission not null default '' comment('用户权限') VARCHAR(20)"`
	CreatedTime time.Time `json:"created_time" xorm:"created_time not null default 'CURRENT_TIMESTAMP' comment('进房时间') TIMESTAMP"`
	UpdatedTime time.Time `json:"updated_time" xorm:"updated_time not null default 'CURRENT_TIMESTAMP' comment('修改时间') TIMESTAMP"`
}

// CreateRoom quick start actions
type RoomNameType struct {
	Name string `json:"name"`
	Type string `json:"type"`
}

func CreateRoom(c *gin.Context) {
	currentUser := c.MustGet("currentUser").(*User)
	db := c.MustGet("db").(*xorm.Engine)

	var nameAndType RoomNameType
	if err := c.BindJSON(&nameAndType); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "invalid args", "code": 401})
		return
	}

	room := Room{
		Name:        nameAndType.Name,
		HostID:      currentUser.ID,
		UUID:        ulid.Make().String(),
		State:       RoomStateActive,
		CreatedTime: time.Now(),
		UpdatedTime: time.Now(),
		Type:        nameAndType.Type,
	}

	if _, err := db.Table(RoomTable).Insert(room); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "build room error", "code": 401})
		log.Println("build room failed")
		return
	}

	participant := Participant{
		UserUUID:    currentUser.UUID,
		RoomUUID:    room.UUID,
		Permission:  PermissionHost,
		Name:        currentUser.Name,
		CreatedTime: time.Now(),
		UpdatedTime: time.Now(),
	}

	if _, err := db.Table(ParticipantTable).Insert(&participant); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "build room error", "code": 401})
		log.Println("participant enter failed")
		return
	}
	//获取participant信息（包含id）
	if _, err := db.Table(ParticipantTable).Where("user_uuid = ? ", participant.UserUUID).Get(&participant); err != nil {
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "build room error", "code": 401})
			log.Println("get participant info failed")
		}
		return
	}

	roomRaw := RoomRaw{}
	//获取room信息（包含id）
	if _, err := db.Table(RoomTable).Where("uuid = ? ", room.UUID).Get(&roomRaw); err != nil {
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "get room info failed", "code": 401})
			log.Println("get room info failed")
		}
		return
	}
	c.JSON(200, roomRaw)
}

func GetRoomInfo(c *gin.Context) {
	uuid := c.Param("uuid")
	db := c.MustGet("db").(*xorm.Engine)
	currentUser := c.MustGet("currentUser").(*User)
	room := Room{}
	has, err := db.Table(RoomTable).Where("uuid=?", uuid).Get(&room)
	if !has {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "room not exit", "code": 401})
		log.Println("room not exist")
		return
	} else if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "server db" + err.Error(), "code": 501})
		log.Println("server db" + err.Error())
		return
	}
	//获取参与者列表
	err = db.Table(ParticipantTable).Where("room_uuid = ?", room.UUID).Iterate(new(Participant), func(i int, bean interface{}) error {
		p := bean.(*Participant)
		room.Participants = append(room.Participants, *p)
		return nil
	})
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "enter room failed", "code": 401})
		log.Println("enter room failed")
		return
	}

	for _, participant := range room.Participants {
		if participant.UserUUID == currentUser.UUID {
			room.MySelf = participant
		}
	}
	c.JSON(http.StatusOK, room)
}

// GetRoomRTC get room info
func GetRoomRTC(c *gin.Context) {
	cfg := c.MustGet("config").(*conf.Config)
	room := c.MustGet("room").(*Room)
	currentUser := c.MustGet("currentUser").(*User)

	mg := rtc.NewManager(&auth.Credentials{
		AccessKey: cfg.QiniuService.AccessKey,
		SecretKey: []byte(cfg.QiniuService.SecretKey),
	})

	access := rtc.RoomAccess{
		AppID:      cfg.QiniuService.RTCAppID,
		RoomName:   room.UUID,
		UserID:     currentUser.UUID,
		ExpireAt:   time.Now().Unix() + 600,
		Permission: room.MySelf.Permission,
	}

	token, err := mg.GetRoomToken(access)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "get room token failed", "code": 401})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"user_uuid": currentUser.UUID,
		"token":     token,
	})
}
func ListRoom(c *gin.Context) {
	db := c.MustGet("db").(*xorm.Engine)
	var roomList []RoomRaw
	//获取房间列表
	err := db.Table(RoomTable).Iterate(new(RoomRaw), func(i int, bean interface{}) error {
		p := bean.(*RoomRaw)
		roomList = append(roomList, *p)
		return nil
	})
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "get roomlist failed", "code": 401})
		log.Println("get roomlist failed")
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message":  "get roomlist success",
		"roomlist": roomList,
		"code":     200,
	})

}

func EnterRoom(c *gin.Context) {
	uuid := c.Param("uuid")
	db := c.MustGet("db").(*xorm.Engine)
	currentUser := c.MustGet("currentUser").(*User)

	room := Room{}
	//获取房间信息
	if _, err := db.Table(RoomTable).Where("uuid=?", uuid).Get(&room); err != nil {
		if errors.Is(err, xorm.ErrNotExist) {
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"message": "enter room failed", "code": 401})
		} else {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "enter room failed", "code": 401})
		}
		return
	}

	participant := Participant{
		UserUUID:    currentUser.UUID,
		RoomUUID:    room.UUID,
		Permission:  PermissionUser,
		Name:        currentUser.Name,
		CreatedTime: time.Now(),
		UpdatedTime: time.Now(),
	}
	//判断此用户是否在房中
	has, err := db.Table(ParticipantTable).Where("user_uuid = ? AND room_uuid = ?", currentUser.UUID, room.UUID).Exist(&participant)
	if has && err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "enter room failed", "code": 401})
		log.Println("enter room failed")
	} else if !has {
		//如果没在房中，加入房间
		if _, err := db.Table(ParticipantTable).Insert(&participant); err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "enter room failed", "code": 401})
			log.Println("enter room failed")
			return
		}
	}

	//获取参与者列表
	err = db.Table(ParticipantTable).Where("room_uuid = ?", room.UUID).Iterate(new(Participant), func(i int, bean interface{}) error {
		p := bean.(*Participant)
		room.Participants = append(room.Participants, *p)
		return nil
	})
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "enter room failed", "code": 401})
		log.Println("enter room failed")
		return
	}

	for _, participant := range room.Participants {
		if participant.UserUUID == currentUser.UUID {
			room.MySelf = participant
		}
	}

	c.Set("room", room)
	c.Next()
}

func ExitRoom(c *gin.Context) {
	uuid := c.Param("uuid")
	db := c.MustGet("db").(*xorm.Engine)
	currentUser := c.MustGet("currentUser").(*User)

	room := Room{}
	has, err := db.Table(RoomTable).Where("uuid=?", uuid).Get(&room)
	if !has {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "room not exit", "code": 401})
		log.Println("room not exist")
		return
	} else if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "server database err: " + err.Error(), "code": 501})
		log.Println("server database err" + err.Error())
		return
	}

	participant := Participant{
		UserUUID: currentUser.UUID,
		RoomUUID: room.UUID,
	}
	//判断此用户是否在房中
	has, err = db.Table(ParticipantTable).Where("user_uuid = ? AND room_uuid = ?", currentUser.UUID, room.UUID).Get(&participant)
	if has {
		//如果在，则删除；
		if _, err := db.Table(ParticipantTable).Where("user_uuid = ? AND room_uuid = ?", currentUser.UUID, room.UUID).Delete(&participant); err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "exit room failed", "code": 401})
			log.Println("exit room failed")
			return
		}
	} else if !has {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"message": "not found participant in room", "code": 401})
		log.Println("not found participant in room")
		return
	} else {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "exit room failed", "code": 401})
		log.Println("exit room failed")
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "enter room success",
		"code":    200,
	})
}
