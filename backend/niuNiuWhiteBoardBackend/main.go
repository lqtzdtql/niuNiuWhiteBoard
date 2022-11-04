package main

import (
	"github.com/gin-gonic/gin"
	"log"
	"niuNiuWhiteBoardBackend/common/database"
	"niuNiuWhiteBoardBackend/config"
	"niuNiuWhiteBoardBackend/models"
)

func main() {
	//初始化数据
	conf.Load()
	gin.SetMode(gin.DebugMode) //开发环境
	//gin.SetMode(gin.ReleaseMode) //线上环境

	//初始化mysql
	db, err := database.InitDatabase()
	if err != nil {
		log.Fatal("Database init fatal")
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
	api := r.Group("v1", models.Auth)
	{
		api.GET("/userinfo/:uuid", models.Info)
		api.POST("/rooms", models.CreateRoom)
		api.GET("/rooms/:uuid", models.EnterRoom, models.GetRoomInfo)
		api.GET("/rooms/:uuid/exit", models.ExitRoom)
		api.GET("/rooms/:uuid/rtc", models.EnterRoom, models.GetRoomRTC)
		api.GET("/logout", models.Logout)
	}

	r.Run(":8282") // listen and serve on 0.0.0.0:8080
}
