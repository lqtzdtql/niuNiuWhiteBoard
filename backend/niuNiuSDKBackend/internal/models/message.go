package models

// ContentType为信令种类
const (
	HEAT_BEAT         = 1
	UPDATE_BOARD      = 2
	OBJECT_NEW        = 3
	OBJECT_MODIFY     = 4
	OBJECT_DELETE     = 5
	SWITCH_BOARD      = 6
	DRAWING_LOCK      = 7
	CREATE_BOARD      = 8
	CAN_LOCK          = 9
	LEAVE_ROOM        = 10
	CANVAS_LIST       = 11
	CUSTOMIZE_MESSAGE = 12
	ENTER_ROOM        = 13
)

type Message struct {
	From         string `json:"from,omitempty"`
	ToRoom       string `json:"toRoom,omitempty"`
	ToWhiteBoard string `json:"toWhiteBoard,omitempty"`
	ToUser       string `json:"toUser,omitempty"`
	ObjectId     string `json:"objectId,omitempty"`
	ContentType  int32  `json:"contentType"`
	Content      string `json:"content,omitempty"`
	Timestamp    int64  `json:"timestamp,omitempty"`
	IsLock       bool   `json:"isLock,omitempty"`
	ReadOnly     bool   `json:"readOnly,omitempty"`
	UserName     string `json:"userName,omitempty"`
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

type ObjectInRedis struct {
	Content   string `json:"content"`
	Timestamp int64  `json:"timestamp"`
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

type BoardInfo struct {
	CanvasId string `json:"canvasId"`
}

// LEAVE_ROOM or ENTER_ROOM
type LeaveEnterRoomRes struct {
	ContentType int32  `json:"contentType"`
	UserName    string `json:"userName"`
}

// CUSTOMIZE_MESSAGE
type CustomizeRes struct {
	ContentType int32  `json:"contentType"`
	Content     string `json:"content"`
}
