package server

import (
	"context"
	"encoding/json"
	"sync"
	"time"

	"niuNiuSDKBackend/common/database"
	"niuNiuSDKBackend/common/log"
	"niuNiuSDKBackend/internal/models"

	"github.com/go-redis/redis/v8"
	"github.com/oklog/ulid/v2"
)

var MyServer = NewServer()

type Server struct {
	Clients    map[string]*Client
	mutex      *sync.Mutex
	Broadcast  chan []byte
	Register   chan *Client
	UnRegister chan *Client
}

func NewServer() *Server {
	return &Server{
		mutex:      &sync.Mutex{},
		Clients:    make(map[string]*Client),
		Broadcast:  make(chan []byte),
		Register:   make(chan *Client),
		UnRegister: make(chan *Client),
	}
}

// 消息类型
const (
	HeartbeatCheckTime = 1  // 心跳检测几秒检测一次
	HeartbeatTime      = 10 // 心跳距离上一次的最大时间
)

// 维持心跳
func (s *Server) heartbeat(ctx context.Context) {
	for {
		clients := make([]*Client, 0)
		// 获取所有的Clients
		s.mutex.Lock()
		for _, c := range s.Clients {
			clients = append(clients, c)
		}
		s.mutex.Unlock()
		for _, c := range clients {
			if time.Now().Unix()-c.HeartbeatTime > HeartbeatTime {
				log.Logger.Info("loginoutll", log.Any("loginoutll", c.UUID))
				user := ExitRoom(c)
				msg := &models.LeaveEnterRoomRes{
					ContentType: models.LEAVE_ROOM,
					UserName:    user.Name,
				}
				message, _ := json.Marshal(msg)
				for _, c := range clients {
					if c.UUID != c.UUID {
						c.Send <- message
					}
				}
				c.Conn.Close()
				delete(s.Clients, c.UUID)
			}
		}
		time.Sleep(time.Second * HeartbeatCheckTime)
	}
}

func (s *Server) register(ctx context.Context) {
	log.Logger.Info("start register", log.Any("start server", "start server..."))
	for {
		select {
		case conn := <-s.Register:
			//时机：进房后开始
			log.Logger.Info("login", log.Any("login", conn.UUID))
			s.Clients[conn.UUID] = conn
			user := models.Participant{}
			if _, err := database.MEngine.Table(models.ParticipantTable).Where("uuid = ? ", conn.UUID).Get(&user); err != nil {
				log.Logger.Error("get participant info failed", log.Any("get participant info failed", err.Error()))
			}
			//进房后广播有人进入房间
			msg := &models.LeaveEnterRoomRes{
				ContentType: models.ENTER_ROOM,
				UserName:    user.Name,
			}

			s.mutex.Lock()
			clients := make([]*Client, 0)
			for _, c := range s.Clients {
				clients = append(clients, c)
			}
			s.mutex.Unlock()

			message, _ := json.Marshal(msg)
			for _, c := range clients {
				if c.UUID != conn.UUID {
					c.Send <- message
				}
			}
		case conn := <-s.UnRegister:
			log.Logger.Info("loginout", log.Any("loginout", conn.UUID))
		case message := <-s.Broadcast:
			MessageHandle(ctx, message, s)
		}
	}
}

// 管理连接
func (s *Server) Start(ctx context.Context) {
	// 检查心跳
	go func() {
		defer func() {
			if r := recover(); r != nil {
				log.Logger.Error("recover", log.Any("recover", r))
			}
		}()
		s.heartbeat(ctx)
	}()

	// 登录注销
	go func() {
		//defer func() {
		//	if r := recover(); r != nil {
		//		log.Logger.Error("recover", log.Any("recover", r))
		//	}
		//}()
		s.register(ctx)
	}()
}

func ExitRoom(conn *Client) models.Participant {
	user := models.Participant{}
	if _, err := database.MEngine.Table(models.ParticipantTable).Where("uuid = ? ", conn.UUID).Get(&user); err != nil {
		log.Logger.Error("get participant info failed", log.Any("get participant info failed", err.Error()))
	}
	_, err := database.MEngine.Table(models.ParticipantTable).Where("uuid = ? AND room_uuid = ?", conn.UUID, conn.RoomUUID).Delete(&user)
	if err != nil {
		log.Logger.Error("exit room failed", log.Any("exit room failed", err.Error()))
	}
	num, _ := database.MEngine.Table(models.ParticipantTable).Where("room_uuid = ?", conn.UUID).Count(&models.Participant{})
	if num == 0 {
		//如果人数为0， 则销毁房间
		if _, err = database.MEngine.Table(models.RoomTable).Where("uuid = ?", conn.RoomUUID).Delete(&models.Room{}); err != nil {
			log.Logger.Error("eixt room success, but room not close", log.Any("eixt room success, but room not close", err.Error()))
		}
	}
	log.Logger.Info("delete", log.Any("delete", conn.UUID))
	return user
}

func MessageHandle(ctx context.Context, message []byte, s *Server) {
	msg := &models.Message{}
	json.Unmarshal(message, msg)
	s.mutex.Lock()
	clients := make([]*Client, 0)
	for _, c := range s.Clients {
		clients = append(clients, c)
	}
	s.mutex.Unlock()

	switch msg.ContentType {
	case models.UPDATE_BOARD:
		upDateBoard(ctx, msg, s)
	case models.SWITCH_BOARD:
		switchBoard(ctx, msg, s, clients)
	case models.OBJECT_NEW:
		objectNew(ctx, msg, clients)
	case models.OBJECT_MODIFY:
		objectModify(ctx, msg, clients)
	case models.OBJECT_DELETE:
		objectDelete(ctx, msg, clients)
	case models.DRAWING_LOCK:
		drawingLock(ctx, msg, s, clients)
	case models.CREATE_BOARD:
		createBoard(ctx, msg, clients)
	case models.CANVAS_LIST:
		canvasList(ctx, msg, s)
	case models.LEAVE_ROOM:
		leaveRoom(ctx, msg, s)
	case models.CUSTOMIZE_MESSAGE:
		customize(ctx, msg, clients)
	case models.HOST_CURRENT:
		getHostCanvasId(ctx, msg, s)
	default:
		for _, c := range clients {
			if msg.ToRoom == c.RoomUUID && msg.From != c.UUID {
				c.Send <- message
			}
		}
	}
}

func broadcast(clients []*Client, msg *models.Message, response []byte) {
	for _, c := range clients {
		if msg.ToRoom == c.RoomUUID && msg.From != c.UUID {
			c.Send <- response
		}
	}
}

func broadcast2All(clients []*Client, msg *models.Message, response []byte) {
	for _, c := range clients {
		if msg.ToRoom == c.RoomUUID {
			c.Send <- response
		}
	}
}

func upDateBoard(ctx context.Context, msg *models.Message, s *Server) {
	upDateRes := models.UpdateBoardRes{
		ContentType:  msg.ContentType,
		ToWhiteBoard: msg.ToWhiteBoard,
	}
	user := models.Participant{}
	if _, err := database.MEngine.Table(models.ParticipantTable).Where("uuid = ? ", msg.From).Get(&user); err != nil {
		log.Logger.Error("get participant info failed", log.Any("get participant info failed", err.Error()))
	}
	if user.Permission == models.PermissionHost {
		err := database.Rdb.Set(ctx, msg.From, msg.ToWhiteBoard, 0).Err()
		if err != nil {
			log.Logger.Error("set host canvasId failed", log.Any("set host canvasId failed", err.Error()))
		}
	}
	// 查找数据库，找出该白板内的所有图形
	objectIds, _ := database.Rdb.SMembers(ctx, msg.ToWhiteBoard).Result()
	var contentString = make([]string, 0)
	for _, objectId := range objectIds {
		contentAndTime := models.ObjectInRedis{}
		obj, _ := database.Rdb.Get(ctx, objectId).Result()
		json.Unmarshal([]byte(obj), &contentAndTime)
		contentString = append(contentString, contentAndTime.Content)
	}
	content, _ := json.Marshal(contentString)
	upDateRes.Content = string(content)
	res, _ := json.Marshal(upDateRes)
	//回复给发送该信令的人
	log.Logger.Debug("upDateBoard test msg.From", log.Any("upDateBoard test msg.From", msg.From))
	s.Clients[msg.From].Send <- res

}

func switchBoard(ctx context.Context, msg *models.Message, s *Server, clients []*Client) {
	switchBoardRes := models.UpdateBoardRes{
		ContentType:  msg.ContentType,
		ToWhiteBoard: msg.ToWhiteBoard,
	}

	// 查找数据库，找出该白板内的所有图形
	objectIds, _ := database.Rdb.SMembers(ctx, msg.ToWhiteBoard).Result()
	var contentString = make([]string, 0)

	for _, objectId := range objectIds {
		contentAndTime := models.ObjectInRedis{}
		obj, _ := database.Rdb.Get(ctx, objectId).Result()
		json.Unmarshal([]byte(obj), &contentAndTime)
		objContent, _ := json.Marshal(contentAndTime.Content)
		contentString = append(contentString, string(objContent))
	}
	content, _ := json.Marshal(contentString)

	if msg.OnlyRead {
		log.Logger.Debug("ReadOnly test", log.Any("ReadOnly test", msg.OnlyRead))
		// SWITCH_BOARD类型的message发给除了发送者以外的该房间的所有人
		switchBoardRes.Content = string(content)
		res2others, _ := json.Marshal(switchBoardRes)
		broadcast(clients, msg, res2others)
	}
	// 给发送者回复UPDATA_BOARD类型的message
	updateBoardRes := models.UpdateBoardRes{
		ContentType:  models.UPDATE_BOARD,
		ToWhiteBoard: msg.ToWhiteBoard,
		Content:      string(content),
	}
	res2From, _ := json.Marshal(updateBoardRes)
	log.Logger.Debug("switchBoard test", log.Any("switchBoard test", res2From))
	s.Clients[msg.From].Send <- res2From
}

func objectNew(ctx context.Context, msg *models.Message, clients []*Client) {
	objectRes := models.ObjectRes{
		ContentType:  msg.ContentType,
		ToWhiteBoard: msg.ToWhiteBoard,
		ObjectId:     msg.ObjectId,
		Content:      msg.Content,
	}
	contentAndTime := models.ObjectInRedis{
		Content:   msg.Content,
		Timestamp: msg.Timestamp,
	}
	contentVal, _ := json.Marshal(contentAndTime)

	err := database.Rdb.Set(ctx, msg.ObjectId, contentVal, 0).Err()
	if err != nil {
		log.Logger.Error("object lost", log.Any("object lost", err.Error()))
	}
	res, _ := json.Marshal(objectRes)
	//以画板id为key，objectid为值，push进redis,以方便后面更新时获取该画板的内容
	err = database.Rdb.SAdd(ctx, msg.ToWhiteBoard, msg.ObjectId).Err()
	if err != nil {
		log.Logger.Error("push objectId to whiteBoard failed", log.Any("push objectId to whiteBoard failed", err.Error()))
	}
	log.Logger.Debug("objectNew test", log.Any("objectNew test", res))
	broadcast(clients, msg, res)
}

func objectModify(ctx context.Context, msg *models.Message, clients []*Client) {
	contentAndTime := models.ObjectInRedis{}
	obj, err := database.Rdb.Get(ctx, msg.ObjectId).Result()
	if err != nil {
		log.Logger.Error("get object failed", log.Any("get object failed", err.Error()))
	}
	json.Unmarshal([]byte(obj), &contentAndTime)
	//修改的时候要比较时间戳
	objectRes := models.ObjectRes{
		ContentType:  msg.ContentType,
		ToWhiteBoard: msg.ToWhiteBoard,
		ObjectId:     msg.ObjectId,
		Content:      msg.Content,
	}
	res, _ := json.Marshal(objectRes)
	if msg.Timestamp > contentAndTime.Timestamp {
		newContentAndTime := models.ObjectInRedis{
			Content:   msg.Content,
			Timestamp: msg.Timestamp,
		}
		newContentVal, _ := json.Marshal(newContentAndTime)
		err := database.Rdb.Set(ctx, msg.ObjectId, newContentVal, 0).Err()
		if err != nil {
			log.Logger.Error("object replace failed", log.Any("object replace failed", err.Error()))
		}
		broadcast(clients, msg, res)
		log.Logger.Debug("objectModify test", log.Any("objectModify test", res))
	}
}

func objectDelete(ctx context.Context, msg *models.Message, clients []*Client) {
	// 删除图形的时候也还要删除图元锁
	err := database.Rdb.Del(ctx, msg.ObjectId+"lock").Err()
	if err == redis.Nil {
		log.Logger.Info("drawinglock not exit", log.Any("drawinglock not exit", msg.ObjectId+"lock"))
	} else if err != nil {
		log.Logger.Error("drawinglock delete failed", log.Any("drawinglock delete failed", err.Error()))
	}
	err = database.Rdb.Del(ctx, msg.ObjectId).Err()
	if err == redis.Nil {
		log.Logger.Info("object not exit", log.Any("object not exit", msg.ObjectId))
	} else if err != nil {
		log.Logger.Error("object delete failed", log.Any("object delete failed", err.Error()))
	}

	objectRes := models.ObjectRes{
		ContentType:  msg.ContentType,
		ToWhiteBoard: msg.ToWhiteBoard,
		ObjectId:     msg.ObjectId,
		Content:      msg.Content,
	}
	res, _ := json.Marshal(objectRes)
	broadcast(clients, msg, res)
	log.Logger.Debug("objectDelete test", log.Any("createBoard test", res))
}

func drawingLock(ctx context.Context, msg *models.Message, s *Server, clients []*Client) {
	if msg.IsLock {
		_, err := database.Rdb.Get(ctx, msg.ObjectId+"lock").Result()
		// 检查此时是否已经上锁
		if err == redis.Nil {
			database.Rdb.Set(ctx, msg.ObjectId+"lock", msg.ObjectId, 0)
			//不存在锁，给发送者回复的CAN_LOCK中 Islock 为false
			drawingLock2From := models.DrawingLockRes{
				ContentType:  models.CAN_LOCK,
				ObjectId:     msg.ObjectId,
				IsLock:       false,
				ToWhiteBoard: msg.ToWhiteBoard,
			}
			res2From, _ := json.Marshal(drawingLock2From)
			s.Clients[msg.From].Send <- res2From
			//给其他人广播DRAWING_LOCK，islock为true
			drawingLock2Others := models.DrawingLockRes{
				ContentType:  models.DRAWING_LOCK,
				ObjectId:     msg.ObjectId,
				IsLock:       true,
				ToWhiteBoard: msg.ToWhiteBoard,
			}
			res2Others, _ := json.Marshal(drawingLock2Others)
			log.Logger.Debug("drawingLock test not exit lock", log.Any("drawingLock test not exit lock", res2Others))
			broadcast(clients, msg, res2Others)

		} else if err != nil {
			log.Logger.Error("get objectlock failed", log.Any("get objectlock failed", err.Error()))
		} else {
			//存在锁，给发送者回复的CAN_LOCK中 islock 为true。
			exitLock2From := models.DrawingLockRes{
				ContentType:  models.CAN_LOCK,
				ObjectId:     msg.ObjectId,
				IsLock:       true,
				ToWhiteBoard: msg.ToWhiteBoard,
			}
			resExitLock2From, _ := json.Marshal(exitLock2From)
			s.Clients[msg.From].Send <- resExitLock2From
			log.Logger.Debug("drawingLock test Exit lock", log.Any("drawingLock test not exit lock", resExitLock2From))
		}
	} else {
		// 如果存在，就删除这个锁，并且向发送者以外的所有人广播解锁，如果不存在，就忽略
		_, err := database.Rdb.Get(ctx, msg.ObjectId+"lock").Result()
		if err == redis.Nil {
			log.Logger.Info("not exit the drawinglock", log.Any("not exit the drawinglock", msg.ObjectId))
		} else if err != nil {
			log.Logger.Error("get objectlock failed", log.Any("get objectlock failed", err.Error()))
		} else {
			unLock2Others := models.DrawingLockRes{
				ContentType:  models.DRAWING_LOCK,
				ObjectId:     msg.ObjectId,
				IsLock:       false,
				ToWhiteBoard: msg.ToWhiteBoard,
			}
			resUnLock2Others, _ := json.Marshal(unLock2Others)
			broadcast(clients, msg, resUnLock2Others)
			err = database.Rdb.Del(ctx, msg.ObjectId+"lock").Err()
			if err != nil {
				log.Logger.Error("drawlock delete failed", log.Any("drawinglock delete failed", err.Error()))
			}
			log.Logger.Debug("Unlock test", log.Any("Unlock test", resUnLock2Others))
		}
	}
}

func createBoard(ctx context.Context, msg *models.Message, clients []*Client) {
	canvas := models.WhiteBoard{
		UUID:        ulid.Make().String(),
		RoomUUID:    msg.ToRoom,
		CreatedTime: time.Now(),
		UpdatedTime: time.Now(),
	}
	//白板table中添加记录，给房间内所有用户广播白板id
	if _, err := database.MEngine.Table(models.WhiteBoardTable).Insert(&canvas); err != nil {
		log.Logger.Error("build room failed", log.Any("build room failed", err.Error()))
	}

	canvasInfo, _ := json.Marshal(models.BoardInfo{CanvasId: canvas.UUID})
	createCanvas := models.CreateBoardRes{
		ContentType: msg.ContentType,
		UserName:    msg.From,
		Content:     string(canvasInfo),
	}
	res, _ := json.Marshal(createCanvas)
	log.Logger.Debug("createBoard test", log.Any("createBoard test", res))
	broadcast2All(clients, msg, res)
}

func canvasList(ctx context.Context, msg *models.Message, s *Server) {
	// 此处给用户返回该房间的所有白板id
	cavList := make([]string, 0)
	//获取该房间的白板列表
	log.Logger.Debug("canvasList test", log.Any("canvasList test", msg.ToRoom))
	database.MEngine.Table(models.WhiteBoardTable).Where("room_uuid = ? and deleted_time is null", msg.ToRoom).Iterate(new(models.WhiteBoard), func(i int, bean interface{}) error {
		p := bean.(*models.WhiteBoard)
		cavansId := p.UUID
		cavList = append(cavList, cavansId)
		return nil
	})
	canvList, _ := json.Marshal(cavList)
	cavListRes := models.CanvasListRes{
		ContentType: msg.ContentType,
		Content:     string(canvList),
	}
	res, _ := json.Marshal(cavListRes)

	if _, ok := s.Clients[msg.From]; ok {
		s.Clients[msg.From].Send <- res
	}
}

func leaveRoom(ctx context.Context, msg *models.Message, s *Server) {
	user := models.Participant{}
	if _, err := database.MEngine.Table(models.ParticipantTable).Where(" name = ? ", msg.UserName).Get(&user); err != nil {
		log.Logger.Error("get participant info failed", log.Any("get participant info failed", err.Error()))
	}
	log.Logger.Debug("leaveRoom test", log.Any("leaveRoom test", user.Name))
	if c, ok := s.Clients[user.UUID]; ok {
		ExitRoom(c)
		leave := &models.LeaveEnterRoomRes{
			ContentType: models.LEAVE_ROOM,
			UserName:    user.Name,
		}
		message, _ := json.Marshal(leave)

		s.mutex.Lock()
		for _, c := range s.Clients {
			if c.UUID != c.UUID {
				c.Send <- message
			}
		}
		s.mutex.Unlock()
		c.Conn.Close()
		delete(s.Clients, c.UUID)
	}
}

func customize(ctx context.Context, msg *models.Message, clients []*Client) {
	custRes := models.CustomizeRes{
		ContentType: msg.ContentType,
		Content:     msg.Content,
	}
	log.Logger.Debug("customize test", log.Any("customize test", custRes))
	res, _ := json.Marshal(custRes)
	broadcast2All(clients, msg, res)
}

func getHostCanvasId(ctx context.Context, msg *models.Message, s *Server) {
	hostCanvasId, err := database.Rdb.Get(ctx, msg.From).Result()
	if err != nil {
		log.Logger.Error("get object failed", log.Any("get object failed", err.Error()))
	}
	hostCanvasIdRes := models.HostCanvasIdRes{
		ContentType: msg.ContentType,
		Content:     hostCanvasId,
	}
	res, _ := json.Marshal(hostCanvasIdRes)
	if _, ok := s.Clients[msg.From]; ok {
		s.Clients[msg.From].Send <- res
	}
}
