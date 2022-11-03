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
	RoomTypeTeaching = "teaching_room"
	RoomTypePlaying  = "playing_room"
)

type Room struct {
	ID           int64         `json:"id" xorm:"id pk autoincr comment('主键') BIGINT(20)"`
	UUID         string        `json:"uuid" xorm:"uuid not null unique 'uuid' comment('房间唯一标识符') index VARCHAR(128)"` // 房间对外的 UUID，同时将作为 RTC Room Name
	Name         string        `json:"name" xorm:"name not null default '' comment('房间名称') VARCHAR(50)"`
	HostID       int64         `json:"host_id"  xorm:" host_id not null default 0 comment('主持人标识符') BIGINT(20)"`
	CreatedTime  int64         `json:"created_time" xorm:"created_time not null default 0 comment('创建时间') index INT(10)"`
	UpdatedTime  int64         `json:"updated_time" xorm:"updated_time not null default 0 comment('修改时间') index INT(10)"`
	State        string        `json:"state " xorm:"-"` // 房间状态: active, stored
	Type         string        `json:"type"  xorm:"type not null default '' comment('房间类型') index VARCHAR(20)"`
	MySelf       Participant   `json:"myself" xorm:"-"`
	Participants []Participant `json:"participants" xorm:"-"`
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
	ID          int64  `json:"id"  xorm:"id pk autoincr comment('主键') BIGINT(20)"`
	UUID        string `json:"uuid" xorm:"uuid not null unique 'uuid' comment(参会唯一标识符') index VARCHAR(128)"`
	UserID      int64  `json:"user_id,omitempty" xorm:"user_id not null default 0 comment('用户ID') index INT(20)"`
	RoomID      int64  `json:"room_id,omitempty" xorm:"room_id not null default 0 comment('房间ID') index INT(20)" `
	Permission  string `json:"permission" xorm:"permission not null default '' comment('用户权限') VARCHAR(20)"`
	Name        string `json:"name" xorm:"name not null default '' comment('参会者姓名') VARCHAR(50)"`
	CreatedTime int64  `json:"created_time" xorm:"created_time not null default 0 comment('创建时间') index INT(10)"`
	UpdatedTime int64  `json:"updated_time" xorm:"updated_time not null default 0 comment('修改时间') index INT(10)"`
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
		CreatedTime: time.Now().Unix(),
		UpdatedTime: time.Now().Unix(),
		Type:        nameAndType.Type,
	}

	if _, err := db.Table(RoomTable).Insert(room); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "build room error", "code": 401})
		log.Println("build room failed")
		return
	}

	participant := Participant{
		UserID:      currentUser.ID,
		RoomID:      room.ID,
		Permission:  PermissionHost,
		UUID:        currentUser.UUID,
		Name:        currentUser.Name,
		CreatedTime: time.Now().Unix(),
		UpdatedTime: time.Now().Unix(),
	}

	if _, err := db.Table(ParticipantTable).Insert(participant); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "participant enter failed", "code": 401})
		log.Println("participant enter failed")
		return
	}
	room.Participants = append(room.Participants, participant)

	c.JSON(200, room)
}

func GetRoomInfo(c *gin.Context) {
	uuid := c.Param("uuid")
	db := c.MustGet("db").(*xorm.Engine)
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
		"userID": currentUser.UUID,
		"token":  token,
	})
}

func GetRoom(c *gin.Context) {
	uuid := c.Param("uuid")
	db := c.MustGet("db").(*xorm.Engine)
	currentUser := c.MustGet("currentUser").(*User)
	room := Room{}

	if _, err := db.Table(RoomTable).Where("uuid=?", uuid).Get(&room); err != nil {
		if errors.Is(err, xorm.ErrNotExist) {
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"message": "enter room failed", "code": 401})
		} else {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "enter room failed", "code": 401})
		}
		return
	}

	participant := &Participant{
		UserID:     currentUser.ID,
		RoomID:     room.ID,
		Permission: PermissionUser,
		UUID:       currentUser.UUID,
		Name:       currentUser.Name,
	}

	if has, err := db.Table(ParticipantTable).Where("user_id = ? AND room_id = ?", currentUser.ID, room.ID).Get(participant); err != nil {
		if !has {
			if _, err := db.Table(ParticipantTable).Insert(participant); err != nil {
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "enter room failed", "code": 401})
				log.Println("enter room failed")
			}
		} else {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "enter room failed", "code": 401})
			log.Println("enter room failed")
		}
		return
	}
	if err := db.Table(ParticipantTable).Where("room_id = ?", room.ID).Find(&room.Participants); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "enter room failed", "code": 401})
		log.Println("enter room failed")
		return
	}

	for _, participant := range room.Participants {
		if participant.UserID == currentUser.ID {
			room.MySelf = participant
		}
	}

	c.Set("room", room)
	c.Next()
}
