package main

import "C"
import (
	"github.com/gin-gonic/gin"

	"niuNiuWhiteBoardBackend/common/database"
	"niuNiuWhiteBoardBackend/common/log"
	"niuNiuWhiteBoardBackend/config"
	"niuNiuWhiteBoardBackend/models"
)

func main() {
	log.InitLogger(conf.Cfg.LogConfig.Path, conf.Cfg.LogConfig.Level)
	log.Logger.Info("config", log.Any("config", conf.Cfg))
	//初始化数据
	gin.SetMode(gin.DebugMode) //开发环境
	//gin.SetMode(gin.ReleaseMode) //线上环境

	//初始化mysql
	db, err := database.InitDatabase()
	if err != nil {
		log.Logger.Fatal("Database init fatal")
	}

	r := gin.Default()

	// API ping
	r.GET("/ping", func(c *gin.Context) {
		c.String(200, "pong")
	})
	r.Use(func(c *gin.Context) {
		c.Set("db", db)
	})

	r.POST("/signup", models.SignupByMobile)
	r.POST("/login", models.Login)
	v1 := r.Group("v1", models.Auth)
	{
		v1.GET("/userinfo/:uuid", models.Info)

		v1.POST("/rooms", models.CreateRoom)
		v1.GET("/roomlist", models.ListRoom)
		v1.GET("/rooms/:uuid", models.EnterRoom, models.GetRoomInfo)
		v1.GET("/rooms/:uuid/exit", models.ExitRoom)
		v1.GET("/rooms/:uuid/rtc", models.EnterRoom, models.GetRoomRTC)
		v1.GET("/logout", models.Logout)
	}

	r.Run(":8282") // listen and serve on 0.0.0.0:8080
}
