package server

import (
	"encoding/json"
	"niuNiuSDKBackend/common/log"
	"niuNiuSDKBackend/internal/models"
	"niuNiuSDKBackend/internal/service"
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
	HeartbeatCheckTime = 9  // 心跳检测几秒检测一次
	HeartbeatTime      = 20 // 心跳距离上一次的最大时间
)

// 维持心跳
func (s *Server) heartbeat() {
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
				delete(s.Clients, c.UUID)
				close(c.Send)
			}
		}
		time.Sleep(time.Second * HeartbeatCheckTime)
	}
}

func (s *Server) register() {
	log.Logger.Info("start register", log.Any("start server", "start server..."))
	for {
		select {
		case conn := <-s.Register:
			//时机：进房后开始
			log.Logger.Info("login", log.Any("login", "new user login in"+conn.UUID))
			s.Clients[conn.UUID] = conn
			msg := &models.Message{
				From:   "niuNiuWhiteBoard",
				ToRoom: conn.UUID,
				//TODO: 此处查找出指定房间的数据库，将所有保存的绘制信息找出，并且conn.Send <- 绘图信息
				Content: "welcome!",
			}
			message, _ := json.Marshal(msg)
			conn.Send <- message
		case conn := <-s.UnRegister:
			log.Logger.Info("loginout", log.Any("loginout", conn.UUID))
			if _, ok := s.Clients[conn.UUID]; ok {
				close(conn.Send)
				delete(s.Clients, conn.UUID)
			}
		case message := <-s.Broadcast:
			msg := &models.Message{}
			json.Unmarshal(message, msg)
			service.MessageHandle(msg, s)
		}
	}
}

// 管理连接
func (s *Server) Start() {
	// 检查心跳
	go func() {
		defer func() {
			if r := recover(); r != nil {
				log.Logger.Error("recover", log.Any("recover", r))
			}
		}()
		s.heartbeat()
	}()

	// 注册注销
	go func() {
		defer func() {
			if r := recover(); r != nil {
				log.Logger.Error("recover", log.Any("recover", r))
			}
		}()
		s.register()
	}()
}
