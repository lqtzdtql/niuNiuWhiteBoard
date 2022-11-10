package main

import "C"
import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
	"niuNiuWhiteBoardBackend/internal/models"

	"niuNiuWhiteBoardBackend/common/database"
	"niuNiuWhiteBoardBackend/common/log"
	"niuNiuWhiteBoardBackend/config"
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
		log.Logger.Fatal("Database init fatal", log.Any("Database init fatal", err.Error()))
	}

	r := gin.Default()

	r.Use(func(c *gin.Context) {
		c.Set("db", db)
	})
	r.Use(Cors())

	r.POST("/signup", models.SignupByMobile)
	r.POST("/login", models.Login)
	v1 := r.Group("v1", models.Auth)
	{
		v1.GET("/userinfo/:uuid", models.Info)

		v1.POST("/rooms", models.CreateRoom)
		v1.GET("/roomlist", models.ListRoom)
		v1.GET("/rooms/:uuid", models.GetRoomInfo)
		v1.GET("/rooms/:uuid/exit", models.ExitRoom)
		v1.GET("/rooms/:uuid/rtc", models.EnterRoom, models.GetRoomRTC)
		v1.GET("/rooms/:uuid/whiteboard", models.EnterRoom, models.GetRoomWhiteBoard)
		v1.GET("/logout", models.Logout)
	}

	r.Run(":8282") // listen and serve on 0.0.0.0:8080
}

func Cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		method := c.Request.Method
		origin := c.Request.Header.Get("Origin") //请求头部
		if origin != "" {
			c.Header("Access-Control-Allow-Origin", "*")
			c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE, UPDATE")
			c.Header("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept, Authorization, Access-Token")
			c.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Cache-Control, Content-Language, Content-Type, Refresh-Token")
			c.Header("Access-Control-Allow-Credentials", "true")
		}
		//允许类型校验
		if method == "OPTIONS" {
			c.JSON(http.StatusOK, "ok!")
		}

		defer func() {
			if err := recover(); err != nil {
				log.Logger.Error("HttpError", zap.Any("HttpError", err))
			}
		}()

		c.Next()
	}
}
