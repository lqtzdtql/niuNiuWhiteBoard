package models

import (
	"net/http"
	"niuNiuWhiteBoardBackend/common/log"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/go-xorm/xorm"
	"github.com/oklog/ulid/v2"

	"niuNiuWhiteBoardBackend/common/utils"
	conf "niuNiuWhiteBoardBackend/config"
)

const (
	ACCESS_TOKEN           = "Access-Token"
	REFRESH_TOKEN          = "Refresh-Token"
	COOKIE_TOKEN           = "User_UUID"
	HEADER_FORWARDED_PROTO = "X-Forwarded-Proto"
)

const (
	SECRETKEY = "42wqTE23123wffLU94342wgldgFs"
	MAXAGE    = 3600 * 24
)

type UserMobile struct {
	Mobile string `form:"mobile" json:"mobile" binding:"required"`
	Passwd string `form:"passwd" json:"passwd" binding:"required,max=20,min=6"`
	Name   string `form:"name" json:"name" binding:"required"`
}

type UserMobilePasswd struct {
	Mobile string `form:"mobile" json:"mobile" binding:"required"`
	Passwd string `form:"passwd" json:"passwd" binding:"required,max=20,min=6"`
}

type Mobile struct {
	Mobile string `form:"mobile" json:"mobile" binding:"required"`
}

// 手机 + 密码登录
func Login(c *gin.Context) {
	db := c.MustGet("db").(*xorm.Engine)

	var userMobile UserMobilePasswd
	if err := c.BindJSON(&userMobile); err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "登录参数错误", "code": 401})
		log.Logger.Error("login invalid args", log.Any("login invalid args", err.Error()))
		return
	}
	user := User{Mobile: userMobile.Mobile}
	has, err := db.Table(UsersTable).Get(&user)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "服务器错误", "code": 500})
		log.Logger.Error("server database error", log.Any("server database error", err.Error()))
		return
	}
	if !has {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "手机不存在", "code": 401})
		log.Logger.Warn("mobile not exist", log.Any("mobile not exist", userMobile.Mobile))
		return
	}

	if utils.Sha1En(userMobile.Passwd) != user.Passwd {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "密码错误", "code": 401})
		log.Logger.Warn("password error", log.Any("password error", user.Passwd))
		return
	}

	user.UserState = UserStateOnline
	if err := DoLogin(c, user); err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "登录失败", "code": 401})
		log.Logger.Error("token set error", log.Any("token set error", err.Error()))
		return
	}

	userRow := UserRow{}
	_, err = db.Table(UsersTable).Where("uuid=?", user.UUID).Get(&userRow)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "服务器错误", "code": 501})
		log.Logger.Error("server database error", log.Any("server database error", err.Error()))
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":   "登录成功",
		"user_info": userRow,
		"code":      200,
	})
	c.Set("currentUser", user)
	c.Next()
	return
}

// 注销登录
func Logout(c *gin.Context) {
	currentUser := c.MustGet("currentUser").(*User)
	secure := IsHttps(c)

	c.SetCookie(COOKIE_TOKEN, "", -1, "/", "", secure, true)
	c.SetCookie(ACCESS_TOKEN, "", -1, "/", "", secure, true)
	c.SetCookie(REFRESH_TOKEN, "", -1, "/", "", secure, true)
	currentUser.UserState = UserStateOffline
	c.JSON(http.StatusOK, gin.H{
		"message": "注销成功",
		"code":    200,
	})
	return
}

// 手机号注册
func SignupByMobile(c *gin.Context) {
	db := c.MustGet("db").(*xorm.Engine)
	var userMobile UserMobile
	if err := c.BindJSON(&userMobile); err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "注册参数错误", "code": 401})
		log.Logger.Error("login invalid args", log.Any("login invalid args", err.Error()))
		return
	}
	user := User{
		Mobile: userMobile.Mobile,
	}

	has, err := db.Table(UsersTable).Exist(&user)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "服务器错误", "code": 501})
		log.Logger.Error("server database error", log.Any("server database error", err.Error()))
		return
	}
	if has {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "手机已注册", "code": 401})
		log.Logger.Warn("mobile or user_name has existed", log.Any("user_mobile", userMobile.Mobile))
		return
	}

	user.Name = userMobile.Name
	user.UUID = ulid.Make().String()
	user.Passwd = utils.Sha1En(userMobile.Passwd)
	user.CreatedTime = time.Now()
	user.UpdatedTime = time.Now()

	if _, err = db.Table(UsersTable).Insert(user); err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "注册失败", "code": 401})
		log.Logger.Error("register failed", log.Any("register failed", err.Error()))
		return
	}

	userRow := UserRow{}
	_, err = db.Table(UsersTable).Where("uuid=?", user.UUID).Get(&userRow)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "服务器错误", "code": 501})
		log.Logger.Error("server database error", log.Any("server database error", err.Error()))
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":   "注册成功",
		"user_info": userRow,
		"code":      200,
	})
	return
}

func Info(c *gin.Context) {
	uuid := c.Param("uuid")
	println(uuid)
	db := c.MustGet("db").(*xorm.Engine)

	user := User{}
	userRow := UserRow{}
	user.UUID = uuid
	has, err := db.Table(UsersTable).Where("uuid=?", user.UUID).Get(&userRow)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "服务器错误", "code": 501})
		log.Logger.Error("server database error", log.Any("server database error", err.Error()))
		return
	}
	if !has {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "用户不存在", "code": 401})
		log.Logger.Warn("user not exist", log.Any("user_uuid", user.UUID))
		return
	}
	//隐藏手机号中间数字
	s := userRow.Mobile
	userRow.Mobile = string([]byte(s)[0:3]) + "****" + string([]byte(s)[7:])
	c.JSON(http.StatusOK, userRow)
	return
}

func DoLogin(c *gin.Context, user User) error {
	if conf.Cfg.OpenJwt { //返回jwt
		refreshClaims := &CustomClaims{
			UserId: user.ID,
			StandardClaims: jwt.StandardClaims{
				ExpiresAt: time.Now().Add(time.Duration(MAXAGE+1800) * time.Second).Unix(), // 过期时间，必须设置
			},
		}
		refreshToken, err := refreshClaims.MakeToken()
		if err != nil {
			return err
		}
		c.Writer.Header().Set(REFRESH_TOKEN, refreshToken)
	}
	return nil
}
