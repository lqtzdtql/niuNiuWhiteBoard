package models

const (
	PONG = "pong"

	//ContentType为消息类型，如心跳，普通信令，一般绘图消息。
	HEAT_BEAT     = 1
	SIGNALING     = 2
	OBJECT_NEW    = 3
	OBJECT_MODIFY = 4
	OBJECT_DELETE = 5
	SWITCH_BOARD  = 6
	DRAWING_LOCK  = 7
	ALL_DRAWING   = 8
)

//心跳不需要保存，发送给指定参与者。
//普通信令不需要保存，广播给房间内所有人。
//图元锁不需要保存，广播给房间内所有人，端上判断是否接收该图元锁。
//图形消息需要保存，每次有人切换新的白板，发送该白板的所有信息给指定参与者。

type Message struct {
	From         string `json:"from"`
	ToRoom       string `json:"toRoom"` //房间号
	ToWhiteBoard string `json:"toWhiteBoard"`
	ObjectId     string `json:"objectId"`
	Content      string `json:"content"`
	ContentType  int32  `json:"contentType"`
}
