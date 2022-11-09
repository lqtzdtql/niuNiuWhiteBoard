package models

import (
	"net/http"
	"niuNiuWhiteBoardBackend/common/log"

	"github.com/gin-gonic/gin"
	"github.com/go-xorm/xorm"

	conf "niuNiuWhiteBoardBackend/config"
)

func Auth(c *gin.Context) {
	db := c.MustGet("db").(*xorm.Engine)

	//开启jwt
	if conf.Cfg.OpenJwt {
		accessToken, has := GetParam(c, ACCESS_TOKEN)
		if !has {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "登录失败", "code": 401})
			log.Logger.Warn("login failed", log.Any("accessToken", ACCESS_TOKEN))
			return
		}
		ret, err := ParseToken(accessToken)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "登录失败", "code": 401})
			log.Logger.Error("login failed", log.Any("login failed", err.Error()))
			return
		}
		user := User{}
		has, err = db.Table(UsersTable).Where("id=?", ret.UserId).Get(&user)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "服务器错误", "code": 501})
			log.Logger.Error("server database error", log.Any("server database error", err.Error()))
			return
		}
		if !has {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "手机号不存在", "code": 401})
			log.Logger.Warn("mobile not exist", log.Any("user_id", ret.UserId))
		}

		c.Set("currentUser", &user)
		if err := DoLogin(c, user); err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "登录失败", "code": 401})
			log.Logger.Error("login failed", log.Any("login failed", err.Error()))
			return
		}
		c.Next()
		return
	}
	c.Next()
	return
}
