package service

import (
	"niuNiuSDKBackend/internal/models"
	"niuNiuSDKBackend/internal/server"
)

func MessageHandle(msg *models.Message, s *server.Server) {
	if msg.ContentType == models.SIGNALING {

	} else if msg.ContentType == models.SWITCH_BOARD {
		//TODO: 此处查找出数据库，将所有保存的绘制信息找出，并且conn.Send <- 绘图信息
	} else if msg.ContentType == models.OBJECT_NEW {

	} else if msg.ContentType == models.OBJECT_MODIFY {

	} else if msg.ContentType == models.OBJECT_DELETE {

	} else if msg.ContentType == models.DRAWING_LOCK {

	}
}
