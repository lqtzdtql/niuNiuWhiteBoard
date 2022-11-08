package models

// redis存储，key为ObjectId，map存储DrawingInfo。
// 每次更新时。查询对应的ObjectId，比较时间戳判断是否更新。--diff
// 同时，也使用set，存储whiteboard 和 ObjectId。
// 每次要获取某个白板全部信息时，先查询mysql通过participant表，获取participant的当前whiteboarduuid
// 再通过set以key为whiteboard获得所有的ObjectId,最后通过ObjectId获取所有的DrawingInfo并且发送。

type DrawingInfo struct {
	From         string `json:"from,omitempty"`
	To           string `json:"to,omitempty"` //房间号
	ToWhiteBoard string `json:"toWhiteBoard"`
	ObjectId     string `json:"objectId"`
	Content      string `json:"content,omitempty"`
	ContentType  int32  `json:"contentType,omitempty"`
	UpdateAt     int64  `json:"updateAt"`
}
