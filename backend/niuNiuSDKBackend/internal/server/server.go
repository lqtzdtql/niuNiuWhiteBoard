package server

import (
	"context"
	"encoding/json"
	"github.com/go-redis/redis/v8"
	"github.com/oklog/ulid/v2"
	"niuNiuSDKBackend/common/database"
	"niuNiuSDKBackend/common/log"
	"niuNiuSDKBackend/internal/models"
	"sync"
	"time"
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
	HeartbeatCheckTime = 500 // 心跳检测几秒检测一次
	HeartbeatTime      = 10  // 心跳距离上一次的最大时间
)

// 维持心跳
func (s *Server) heartbeat(ctx context.Context) {
	for {
		// 获取所有的Clients
		s.mutex.Lock()
		clients := make([]*Client, len(s.Clients))
		for _, c := range s.Clients {
			clients = append(clients, c)
		}
		s.mutex.Unlock()

		for _, c := range clients {
			if time.Now().Unix()-c.HeartbeatTime > HeartbeatTime {
				log.Logger.Info("loginout", log.Any("loginout", c.UUID))
				s.UnRegister <- c
			}
		}
		time.Sleep(time.Millisecond * HeartbeatCheckTime)
	}
}

func (s *Server) register(ctx context.Context) {
	log.Logger.Info("start register", log.Any("start server", "start server..."))
	for {
		select {
		case conn := <-s.Register:
			//时机：进房后开始
			log.Logger.Info("login", log.Any("login", "new user login in"+conn.UUID))
			s.Clients[conn.UUID] = conn
			//进房后发个心跳类型给client
			msg := &models.HeatBeatRes{
				ContentType: models.HEAT_BEAT,
			}
			message, _ := json.Marshal(msg)
			conn.Send <- message
		case conn := <-s.UnRegister:
			user := &models.Participant{}
			//广播用户退出的消息
			if _, err := database.MEngine.Table(models.ParticipantTable).Where("Uuid = ? ", conn.UUID).Get(&user); err != nil {
				log.Logger.Error("get participant info failed", log.Any("get participant info failed", err.Error()))
			}
			s.mutex.Lock()
			clients := make([]*Client, len(s.Clients))
			for _, c := range s.Clients {
				clients = append(clients, c)
			}
			s.mutex.Unlock()
			msg := &models.LeaveRoomRes{
				ContentType: models.LEAVE_ROOM,
				LeaveUser:   user.Name,
			}
			message, _ := json.Marshal(msg)

			for _, c := range clients {
				if c.UUID != conn.UUID {
					c.Send <- message
				}
			}

			log.Logger.Info("loginout", log.Any("loginout", conn.UUID))
			if _, ok := s.Clients[conn.UUID]; ok {
				_, err := database.MEngine.Table(models.ParticipantTable).Where("user_uuid = ? AND room_uuid = ?", conn.UUID, conn.UUID).Delete(&user)
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
				close(conn.Send)
				delete(s.Clients, conn.UUID)
			}
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

	// 注册注销
	go func() {
		defer func() {
			if r := recover(); r != nil {
				log.Logger.Error("recover", log.Any("recover", r))
			}
		}()
		s.register(ctx)
	}()
}

func MessageHandle(ctx context.Context, message []byte, s *Server) {
	msg := &models.Message{}
	json.Unmarshal(message, msg)
	s.mutex.Lock()
	clients := make([]*Client, len(s.Clients))
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
		drawLock(ctx, msg, s, clients)
	case models.CREATE_BOARD:
		createBoard(ctx, msg, clients)
	case models.CANVAS_LIST:
		canvasList(ctx, msg, s)
	case models.LEAVE_ROOM:
		leaveRoom(ctx, msg, s)
	default:
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
	// 查找数据库，找出该白板内的所有图形
	objectIds, _ := database.Rdb.SMembers(ctx, msg.ToWhiteBoard).Result()
	var contentString []string
	for _, objectId := range objectIds {
		contentAndTime := models.ObjectInRedis{}
		obj, _ := database.Rdb.Get(ctx, objectId).Result()
		json.Unmarshal([]byte(obj), &contentAndTime)
		objContent, _ := json.Marshal(contentAndTime.Content)
		contentString = append(contentString, string(objContent))
	}
	content, _ := json.Marshal(contentString)
	upDateRes.Content = string(content)
	res, _ := json.Marshal(upDateRes)
	// 回复给发送该信令的人
	s.Clients[msg.From].Send <- res
}

func switchBoard(ctx context.Context, msg *models.Message, s *Server, clients []*Client) {
	switchBoardRes := models.UpdateBoardRes{
		ContentType:  msg.ContentType,
		ToWhiteBoard: msg.ToWhiteBoard,
	}
	// 查找数据库，找出该白板内的所有图形
	objectIds, _ := database.Rdb.SMembers(ctx, msg.ToWhiteBoard).Result()
	var contentString []string
	for _, objectId := range objectIds {
		contentAndTime := models.ObjectInRedis{}
		obj, _ := database.Rdb.Get(ctx, objectId).Result()
		json.Unmarshal([]byte(obj), &contentAndTime)
		objContent, _ := json.Marshal(contentAndTime.Content)
		contentString = append(contentString, string(objContent))
	}
	content, _ := json.Marshal(contentString)

	if msg.ReadOnly {
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
	}
	broadcast(clients, msg, res)
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
	log.Logger.Info("object delete success", log.Any("object delete failed", err.Error()))
}

func drawLock(ctx context.Context, msg *models.Message, s *Server, clients []*Client) {
	if msg.IsLock {
		_, err := database.Rdb.Get(ctx, msg.ObjectId+"lock").Result()
		// 检查此时是否已经上锁
		if err == redis.Nil {
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
		}
	} else {
		// 如果存在，就删除这个锁，如果不存在，就忽略
		_, err := database.Rdb.Get(ctx, msg.ObjectId+"lock").Result()
		if err == redis.Nil {
			log.Logger.Info("not exit the drawinglock", log.Any("not exit the drawinglock", msg.ObjectId))
		} else if err != nil {
			log.Logger.Error("get objectlock failed", log.Any("get objectlock failed", err.Error()))
		} else {
			err = database.Rdb.Del(ctx, msg.ObjectId+"lock").Err()
			if err != nil {
				log.Logger.Error("drawlock delete failed", log.Any("drawinglock delete failed", err.Error()))
			}
		}

	}
}

func createBoard(ctx context.Context, msg *models.Message, clients []*Client) {
	canvas := models.WhiteBoard{
		UUID:     ulid.Make().String(),
		RoomUUID: msg.ToRoom,
	}
	//白板table中添加记录，给房间内所有用户广播白板id
	if _, err := database.MEngine.Table(models.WhiteBoardTable).Insert(&canvas); err != nil {
		log.Logger.Error("build room failed", log.Any("build room failed", err.Error()))
	}
	canvasId, _ := json.Marshal(canvas.UUID)
	createCanvas := models.CreateBoardRes{
		ContentType: msg.ContentType,
		Content:     string(canvasId),
	}
	res, _ := json.Marshal(createCanvas)
	broadcast2All(clients, msg, res)
}

func canvasList(ctx context.Context, msg *models.Message, s *Server) {
	// 此处给用户返回该房间的所有白板id
	var whiteBoard []string
	//获取该房间的白板列表
	database.MEngine.Table(models.WhiteBoardTable).Where("room_uuid = ? and deleted_time is null", msg.ToRoom).Iterate(new(string), func(i int, bean interface{}) error {
		p := bean.(*string)
		whiteBoard = append(whiteBoard, *p)
		return nil
	})
	canvasList, _ := json.Marshal(whiteBoard)
	s.Clients[msg.From].Send <- canvasList
}

func leaveRoom(ctx context.Context, msg *models.Message, s *Server) {
	if c, ok := s.Clients[msg.LeaveUser]; ok {
		s.UnRegister <- c
	}
}
