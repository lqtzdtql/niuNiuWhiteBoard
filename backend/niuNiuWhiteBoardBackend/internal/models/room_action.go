package models

import (
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-xorm/xorm"
	"github.com/oklog/ulid/v2"
	"github.com/qiniu/go-sdk/v7/auth"
	"github.com/qiniu/go-sdk/v7/rtc"

	conf "niuNiuWhiteBoardBackend/config"
)

func CreateRoom(c *gin.Context) {
	currentUser := c.MustGet("currentUser").(*User)
	db := c.MustGet("db").(*xorm.Engine)

	var nameAndType RoomNameType
	if err := c.BindJSON(&nameAndType); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "创建房间参数错误", "code": 401})
		return
	}

	room := Room{
		Name:        nameAndType.Name,
		HostID:      currentUser.ID,
		UUID:        ulid.Make().String(),
		CreatedTime: time.Now(),
		UpdatedTime: time.Now(),
		Type:        nameAndType.Type,
	}
	if has, _ := db.Table(RoomTable).Where("host_id = ? ", room.HostID).Exist(new(Room)); has {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "房主已在房内，无法创建房间", "code": 401})
		log.Println("host in the room, can not build room")
		return
	}

	if _, err := db.Table(RoomTable).Insert(room); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "创建房间失败", "code": 501})
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
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "创建房间失败", "code": 501})
		log.Println("participant enter failed")
		return
	}

	//获取participant信息（包含id）
	if _, err := db.Table(ParticipantTable).Where("user_uuid = ? ", participant.UserUUID).Get(&participant); err != nil {
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "创建房间失败", "code": 501})
			log.Println("get participant info failed")
		}
		return
	}

	roomRaw := RoomRaw{}
	//获取room信息（包含id）
	if _, err := db.Table(RoomTable).Where("uuid = ? ", room.UUID).Get(&roomRaw); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "创建房间失败", "code": 501})
		log.Println("get room info failed" + err.Error())
		return
	}

	err := db.Table(ParticipantTable).Where("room_uuid = ? AND deleted_time is null", roomRaw.UUID).Iterate(new(ParticipantRow), func(i int, bean interface{}) error {
		p := bean.(*ParticipantRow)
		roomRaw.Participants = append(roomRaw.Participants, *p)
		return nil
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "创建房间失败", "code": 501})
		log.Println("enter room failed")
	}

	c.JSON(200, roomRaw)
}

func GetRoomInfo(c *gin.Context) {
	uuid := c.Param("uuid")
	db := c.MustGet("db").(*xorm.Engine)
	roomRaw := RoomRaw{}
	has, err := db.Table(RoomTable).Where("uuid=?", uuid).Get(&roomRaw)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "服务器错误", "code": 501})
		log.Println("server db" + err.Error())
		return
	}
	if !has {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "房间不存在", "code": 401})
		log.Println("room not exist")
		return
	}
	//获取参与者列表
	err = db.Table(ParticipantTable).Where("room_uuid = ? AND deleted_time is null", roomRaw.UUID).Iterate(new(ParticipantRow), func(i int, bean interface{}) error {
		p := bean.(*ParticipantRow)
		roomRaw.Participants = append(roomRaw.Participants, *p)
		return nil
	})
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "获取房间信息失败", "code": 401})
		log.Println("enter room failed")
		return
	}

	c.JSON(http.StatusOK, roomRaw)
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
	var roomList []RoomInfo
	roomList = make([]RoomInfo, 0)
	//获取房间列表
	err := db.Table(RoomTable).Where("deleted_time IS NULL").Iterate(new(RoomInfo), func(i int, bean interface{}) error {
		p := bean.(*RoomInfo)
		roomList = append(roomList, *p)
		return nil
	})
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "服务器错误", "code": 501})
		log.Println("get roomlist failed")
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":  "获取房间列表成功",
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
	has, err := db.Table(RoomTable).Where("uuid=?", uuid).Get(&room)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "服务器错误", "code": 501})
		return
	}
	if !has {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "房间不存在", "code": 401})
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
	has, _ = db.Table(ParticipantTable).Where("user_uuid = ? AND room_uuid = ?", currentUser.UUID, room.UUID).Exist(&Participant{})
	if !has {
		//如果没在房中，加入房间
		if _, err := db.Table(ParticipantTable).Insert(&participant); err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "服务器错误", "code": 501})
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
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "服务器错误", "code": 501})
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
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "服务器错误", "code": 501})
		log.Println("server database err" + err.Error())
		return
	}
	if !has {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "房间不存在", "code": 401})
		log.Println("room not exist")
		return
	}
	participant := new(Participant)
	participant.UserUUID = currentUser.UUID
	participant.RoomUUID = room.UUID
	//判断此用户是否在房中
	has, err = db.Table(ParticipantTable).Where("user_uuid = ? AND room_uuid = ?", currentUser.UUID, room.UUID).Get(participant)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "服务器错误", "code": 501})
		log.Println("exit room failed")
		return
	}
	if !has {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"message": "用户不在房中", "code": 401})
		log.Println("not found participant in room")
		return
	} else {
		//如果在，则删除；
		_, err := db.Table(ParticipantTable).Where("user_uuid = ? AND room_uuid = ?", currentUser.UUID, room.UUID).Delete(participant)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "服务器错误", "code": 501})
			log.Println("exit room failed")
			return
		}
	}

	num, _ := db.Table(ParticipantTable).Where("room_uuid = ?", room.UUID).Count(&Participant{})
	if num == 0 {
		//如果人数为0， 则销毁房间
		if _, err = db.Table(RoomTable).Where("uuid = ?", room.UUID).Delete(&room); err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "服务器错误", "code": 501})
			log.Println("eixt room success, but room not close: " + err.Error())
			return
		}
	}
	// TODO: 房主退房逻辑，待定

	c.JSON(http.StatusOK, gin.H{
		"message": "退房成功",
		"code":    200,
	})
}
