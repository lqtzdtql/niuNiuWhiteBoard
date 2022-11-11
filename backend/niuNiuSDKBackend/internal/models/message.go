package models

// ContentType为信令种类
const (
	HEAT_BEAT     = 1
	UPDATE_BOARD  = 2
	OBJECT_NEW    = 3
	OBJECT_MODIFY = 4
	OBJECT_DELETE = 5
	SWITCH_BOARD  = 6
	DRAWING_LOCK  = 7
	CREATE_BOARD  = 8
	CAN_LOCK      = 9
	LEAVE_ROOM    = 10
	CANVAS_LIST   = 11
)

/*
  心跳不需要保存，发送给指定参与者。
  普通信令不需要保存，广播给房间内所有人。
  图元锁也需要保存，广播给房间内所有人，端上判断是否接收该图元锁。
  图形相关的消息需要保存，每次有人切换新的白板，发送该白板的所有信息给指定参与者。
*/

/*
	redis存储，key为ObjectId，map存储DrawingInfo。--（包含白板id，定时将DrawingInfo转发给拥有该白板的所有用户）
	每次更新时。查询对应的ObjectId，比较时间戳判断是否更新。--diff
	同时，也使用set，存储whiteboard 和 ObjectId。
	每次要获取某个白板全部信息时，先查询mysql通过participant表，获取participant的当前whiteboarduuid
	再通过set以key为whiteboard获得所有的ObjectId,最后通过ObjectId获取所有的DrawingInfo并且发送。
*/

type Message struct {
	From         string `json:"from,omitempty"`
	ToRoom       string `json:"to,omitempty"`
	ToWhiteBoard string `json:"toWhiteBoard,omitempty"`
	ToUser       string `json:"toUser,omitempty"`
	ObjectId     string `json:"objectId,omitempty"`
	ContentType  int32  `json:"contentType"`
	Content      string `json:"content,omitempty"`
	Timestamp    int64  `json:"timestamp,omitempty"`
	IsLock       bool   `json:"isLock,omitempty"`
	ReadOnly     bool   `json:"readOnly,omitempty"`
	LeaveUser    string `json:"leaveUser,omitempty"`
}

// HEAT_BEAT回复
type HeatBeatRes struct {
	ContentType int32 `json:"contentType"`
}

// UPDATE_BOARD，SWITCH_BOARD（切换白板需要将当前用户的currentBoard修改）
type UpdateBoardRes struct {
	ContentType  int32  `json:"contentType"`
	ToWhiteBoard string `json:"toWhiteBoard"`
	Content      string `json:"content"`
}

// OBJECT_NEW，OBJECT_MODIFY，OBJECT_DELETE
type ObjectRes struct {
	ContentType  int32  `json:"contentType"`
	ObjectId     string `json:"objectId"`
	ToWhiteBoard string `json:"toWhiteBoard"`
	Content      string `json:"content"`
}

// DRAWING_LOCK
type DrawingLockRes struct {
	ContentType  int32  `json:"contentType"`
	ObjectId     string `json:"objectId"`
	ToWhiteBoard string `json:"toWhiteBoard"`
	IsLock       bool   `json:"isLock"`
}

// CREATE_BOARD
type CreateBoardRes struct {
	ContentType int32  `json:"contentType"`
	Content     string `json:"content"`
}

// CAN_LOCK
type CanLockRes struct {
	ContentType  int32  `json:"contentType"`
	ObjectId     string `json:"objectId"`
	ToWhiteBoard string `json:"toWhiteBoard"`
}

// LEAVE_ROOM
type LeaveRoomRes struct {
	ContentType int32  `json:"contentType"`
	Content     string `json:"content"`
}
