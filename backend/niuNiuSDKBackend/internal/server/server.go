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
				log.Logger.Info("loginout", log.Any("loginout", c.Name))
				delete(s.Clients, c.Name)
				close(c.Send)
			}
		}
		time.Sleep(time.Second * HeartbeatCheckTime)
	}
}

func (s *Server) register() {
	log.Logger.Info("start server", log.Any("start server", "start server..."))
	for {
		select {
		case conn := <-s.Register:
			//时机：进房后开始
			log.Logger.Info("login", log.Any("login", "new user login in"+conn.Name))
			s.Clients[conn.Name] = conn
			msg := &models.Message{
				From: "niuNiuWhiteBoard",
				To:   conn.Name,
				//TODO: 此处查找出数据库，将所有保存的绘制信息找出，并且conn.Send <- 绘图信息
				Content: "welcome!",
			}
			message, _ := json.Marshal(msg)
			conn.Send <- message
		case conn := <-s.UnRegister:
			//时机：退房后。
			log.Logger.Info("loginout", log.Any("loginout", conn.Name))
			if _, ok := s.Clients[conn.Name]; ok {
				close(conn.Send)
				delete(s.Clients, conn.Name)
			}
		case message := <-s.Broadcast:
			msg := &models.Message{}
			json.Unmarshal(message, msg)
			//
			if msg.To != "" {
				if msg.ContentType == models.OBJECT {
					// 图形相关的信息需要被保存，保存的图形只会在存在socket的一个端上进行保存，防止分布式部署后，信息重复问题
					_, exits := s.Clients[msg.From]
					if exits {
						saveMessage(msg)
					}
					sendRoomMessage(msg, s)
				} else if msg.ContentType == models.REPAINT {
					//TODO: 此处查找出数据库，将所有保存的绘制信息找出，并且conn.Send <- 绘图信息
				} else {
					//转发到房间号为to的所有房间
					client, ok := s.Clients[msg.To]
					if ok {
						client.Send <- message
					}
				}
			}
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

// 发送给房间的消息,需要查询该房间所有参与者再依次发送
func sendRoomMessage(msg *models.Message, s *Server) {
	// 发送给群组的消息，查找该群所有的用户进行发送
	users := service.RoomService.GetUserIdByRoomUuid(msg.To)
	for _, user := range users {
		if user.Uuid == msg.From {
			continue
		}

		client, ok := s.Clients[user.Uuid]
		if !ok {
			continue
		}
		// from是个人，to是群聊uuid。所以在返回消息时，将from修改为群聊uuid
		msgSend := models.Message{
			From:        msg.To,
			To:          msg.From,
			Content:     msg.Content,
			ContentType: msg.ContentType,
		}

		message, err := json.Marshal(&msgSend)
		if err == nil {
			client.Send <- message
		}
	}
}

// 保存消息
func saveMessage(message *models.Message) {
	service.MessageService.SaveMessage(*message)
}
