package models

import (
	"github.com/gin-gonic/gin"
	"github.com/go-xorm/xorm"
	"log"
	"net/http"
	"net/url"
	"niuNiuWhiteBoardBackend/common/utils"
	conf "niuNiuWhiteBoardBackend/config"
)

func GetParam(c *gin.Context, key string) (string, bool) {
	val := c.GetHeader(key)
	if val != "" {
		return val, true
	}
	val, err := c.Cookie(key)
	if err != nil {
		return "", false
	}
	return val, true
}

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
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "login failed", "code": 401})
			log.Println("login failed")
			return
		}
		ret, err := ParseToken(accessToken)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "login failed", "code": 401})
			log.Println("login failed")
			return
		}
		user := User{}
		if has, err := db.Table(UsersTable).Where("id=?", ret.UserId).Get(&user); err != nil {
			if !has {
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "mobile not exist", "code": 401})
				log.Println("mobile not exist")
			} else {
				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "server db" + err.Error(), "code": 501})
				log.Println("server db" + err.Error())
			}
			return
		}
		c.Set("currentUser", &user)
		if err := DoLogin(c, user); err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "login failed", "code": 401})
			log.Println("login failed")
			return
		}
		c.Next()
		return
	}
	//cookie
	_, err = c.Cookie(COOKIE_TOKEN)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "user not exist", "code": 401})
		log.Println("user not exist")
		return
	}
	c.Next()
	return
}
