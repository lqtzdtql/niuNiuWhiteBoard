package server

import (
	"context"
	"encoding/json"
	"github.com/go-redis/redis/v8"
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
	HeartbeatCheckTime = 3  // 心跳检测几秒检测一次
	HeartbeatTime      = 15 // 心跳距离上一次的最大时间
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
				//TODO:广播user_leave信令
				delete(s.Clients, c.UUID)
				close(c.Send)
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
			log.Logger.Info("login", log.Any("login", "new user login in"+conn.UUID))
			s.Clients[conn.UUID] = conn
			msg := &models.Message{
				From:   "niuNiuWhiteBoard",
				ToRoom: conn.RoomUUID,
				//TODO: 此处给用户返回该房间的所有白板id，CANVAS_LIST
				Content: "welcome!",
			}
			message, _ := json.Marshal(msg)
			conn.Send <- message
		case conn := <-s.UnRegister:
			//TODO:用户退出，LEAVE_ROOM广播给所有人
			log.Logger.Info("loginout", log.Any("loginout", conn.UUID))
			if _, ok := s.Clients[conn.UUID]; ok {
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

	if msg.ContentType == models.UPDATE_BOARD {
		//TODO: 查找数据库，找出该白板内的所有图形
		//TODO: 回复给发送该信令的人
	} else if msg.ContentType == models.SWITCH_BOARD {
		//TODO: 查找数据库，找出该白板内的所有图形
		//TODO: 如果是只读模式，UPDATE_BOARD白板的回复发给该房间的所有人，同时修改房间内所有人的currentBoard
		//TODO：如果是协作模式，UPDATE_BOARD白板的回复只给切换白板的人，修改该用户的currentBoard
	} else if msg.ContentType == models.OBJECT_NEW {
		err := database.Rdb.Set(ctx, msg.ObjectId, message, 0).Err()
		if err != nil {
			log.Logger.Error("object lost", log.Any("object lost", err.Error()))
		}
		objectRes := models.ObjectRes{
			ContentType:  msg.ContentType,
			ToWhiteBoard: msg.ToWhiteBoard,
			ObjectId:     msg.ObjectId,
			Content:      msg.Content,
		}
		res, _ := json.Marshal(objectRes)
		broadcast(clients, msg, res)

	} else if msg.ContentType == models.OBJECT_MODIFY {
		object := &models.Message{}
		obj, err := database.Rdb.Get(ctx, msg.ObjectId).Result()
		if err != nil {
			log.Logger.Error("get object failed", log.Any("get object failed", err.Error()))
		}
		json.Unmarshal([]byte(obj), object)
		//修改的时候要比较时间戳
		if msg.Timestamp > object.Timestamp {
			err := database.Rdb.Set(ctx, msg.ObjectId, message, 0).Err()
			if err != nil {
				log.Logger.Error("object replace failed", log.Any("object replace failed", err.Error()))
			}
		}
		objectRes := models.ObjectRes{
			ContentType:  msg.ContentType,
			ToWhiteBoard: msg.ToWhiteBoard,
			ObjectId:     msg.ObjectId,
			Content:      msg.Content,
		}
		res, _ := json.Marshal(objectRes)
		broadcast(clients, msg, res)

	} else if msg.ContentType == models.OBJECT_DELETE {
		// 删除图形的时候也还要删除图元锁
		err := database.Rdb.Del(ctx, msg.ObjectId+"lock").Err()
		if err != nil {
			log.Logger.Error("drawinglock delete failed", log.Any("drawinglock delete failed", err.Error()))
		}
		err = database.Rdb.Del(ctx, msg.ObjectId).Err()
		if err != nil {
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
	} else if msg.ContentType == models.DRAWING_LOCK {
		// 检查此时是否已经上锁
		_, err := database.Rdb.Get(ctx, msg.ObjectId+"lock").Result()
		if err == redis.Nil {
			//TODO:不存在锁，给发送者回复的CAN_LOCK中 islock 为false，给其他人广播DRAWING_LOCK，islock为true

		} else if err != nil {
			log.Logger.Error("get objectlock failed", log.Any("get objectlock failed", err.Error()))
		} else {
			//TODO:存在锁，给发送者回复的CAN_LOCK中 islock 为true。

		}
	} else if msg.ContentType == models.CREATE_BOARD {
		//TODO:白板table中添加记录，给房间内所有用户广播白板id

	} else if msg.ContentType == models.CANVAS_LIST {
		//TODO: 此处给用户返回该房间的所有白板id
	} else if msg.ContentType == models.LEAVE_ROOM {
		//TODO: 如果from和LeaveUser不一样是踢人
		//TODO: 一样则是用户主动退出
	}
}

func broadcast(clients []*Client, msg *models.Message, message []byte) {
	for _, c := range clients {
		if msg.ToRoom == c.RoomUUID && msg.From != c.UUID {
			c.Send <- message
		}
	}
}
