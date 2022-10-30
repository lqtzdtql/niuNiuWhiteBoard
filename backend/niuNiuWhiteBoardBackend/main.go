package main

import (
	"github.com/gin-gonic/gin"
	"net/url"
	conf "niuNiuWhiteBoardBackend/sso/config"
	"niuNiuWhiteBoardBackend/sso/login"
	"niuNiuWhiteBoardBackend/sso/utils"
)

func main() {

	//初始化数据
	Load()
	gin.SetMode(gin.DebugMode) //开发环境
	//gin.SetMode(gin.ReleaseMode) //线上环境
	r := gin.Default()
	r.Use(Auth)
	r.POST("/logout", login.Logout)
	r.POST("/login", login.Login)
	r.POST("/signup/mobile", login.SignupByMobile)
	r.GET("/my/info", login.Info) //用户信息
	r.GET("/pong", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	r.Run(":8285") // listen and serve on 0.0.0.0:8080
}
func Load() {
	c := conf.Config{}
	c.Routes = []string{"/ping", "/login", "/login/mobile", "/signup/mobile", "/signup/mobile/exist"}
	c.OpenJwt = false //开启jwt
	conf.Set(c)
}

func Auth(c *gin.Context) {
	u, err := url.Parse(c.Request.RequestURI)
	if err != nil {
		panic(err)
	}
	if utils.InArrayString(u.Path, &conf.Cfg.Routes) {
		c.Next()
		return
	}
	c.Next()
	return
}
