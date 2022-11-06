package models

import (
	"log"
	"net/http"
	"net/url"

	"github.com/gin-gonic/gin"
	"github.com/go-xorm/xorm"

	"niuNiuWhiteBoardBackend/common/utils"
	conf "niuNiuWhiteBoardBackend/config"
)

func Auth(c *gin.Context) {
	db := c.MustGet("db").(*xorm.Engine)
	u, err := url.Parse(c.Request.RequestURI)
	if err != nil {
		panic(err)
	}
	if utils.InArrayString(u.Path, &conf.Cfg.Routes) {
		c.Next()
		return
	}
	//开启jwt
	if conf.Cfg.OpenJwt {
		accessToken, has := GetParam(c, ACCESS_TOKEN)
		if !has {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "登录失败", "code": 401})
			log.Println("login failed")
			return
		}
		ret, err := ParseToken(accessToken)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "登录失败", "code": 401})
			log.Println("login failed")
			return
		}
		user := User{}
		has, err = db.Table(UsersTable).Where("id=?", ret.UserId).Get(&user)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "服务器错误", "code": 501})
			log.Println("server database error: " + err.Error())
			return
		}
		if !has {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "手机号不存在", "code": 401})
			log.Println("mobile not exist")
		}

		c.Set("currentUser", &user)
		if err := DoLogin(c, user); err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "登录失败", "code": 401})
			log.Println("login failed")
			return
		}
		c.Next()
		return
	}
	c.Next()
	return
}
