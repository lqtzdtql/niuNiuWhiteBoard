package models

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-xorm/xorm"
	"github.com/oklog/ulid/v2"
	"log"
	"net/http"
	"niuNiuWhiteBoardBackend/common/utils"
	conf "niuNiuWhiteBoardBackend/config"
	"strconv"
	"time"

	"github.com/dgrijalva/jwt-go"
)

const (
	ACCESS_TOKEN           = "Access_Token"
	REFRESH_TOKEN          = "Refresh_Token"
	COOKIE_TOKEN           = "UserId"
	HEADER_AUTH            = "Authorization"
	HEADER_FORWARDED_PROTO = "X-Forwarded-Proto"
)

const (
	SECRETKEY = "42wqTE23123wffLU94342wgldgFs"
	MAXAGE    = 3600 * 24
)

type UserMobile struct {
	Mobile string `form:"mobile" json:"mobile" binding:"required"`
	Passwd string `form:"passwd" json:"passwd" binding:"required,max=20,min=6"`
	Code   string `form:"code" json:"code" binding:"required,len=6"`
}
type UserMobileCode struct {
	Mobile string `form:"mobile" json:"mobile" binding:"required"`
	Code   string `form:"code" json:"code" binding:"required,len=6"`
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
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "login invalid args", "code": 401})
		log.Println("login invalid args")
		return
	}
	user := User{Mobile: userMobile.Mobile}
	if has, err := db.Table(UsersTable).Get(&user); err != nil {
		if !has {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "mobile not exist", "code": 401})
			log.Println("mobile not exist")
		} else {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "server db" + err.Error(), "code": 501})
			log.Println("server db" + err.Error())
		}
		return
	}

	if utils.Sha1En(userMobile.Passwd) != user.Passwd {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "password err", "code": 401})
		log.Println("password err")
		return
	}

	user.UserState = UserStateOnline
	if err := DoLogin(c, user); err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "cookie set error", "code": 401})
		log.Println("cookie set error")
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "login success",
		"code":    200,
	})
	c.Set("currentUser", user)
	c.Next()
	return
}

// 注销登录
func Logout(c *gin.Context) {
	user := c.MustGet("user").(*User)
	secure := IsHttps(c)

	c.SetCookie(COOKIE_TOKEN, "", -1, "/", "", secure, true)
	c.SetCookie(ACCESS_TOKEN, "", -1, "/", "", secure, true)
	c.SetCookie(REFRESH_TOKEN, "", -1, "/", "", secure, true)
	user.UserState = UserStateOffline
	c.JSON(http.StatusOK, gin.H{
		"message": "logout success",
		"token":   200,
	})
	return
}

// 手机号注册
func SignupByMobile(c *gin.Context) {
	db := c.MustGet("db").(*xorm.Engine)
	var userMobile UserMobile
	if err := c.BindJSON(&userMobile); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "login invalid args", "code": 401})
		log.Println("login invalid args")
		return
	}
	fmt.Println(userMobile)
	user := User{Mobile: userMobile.Mobile}
	has, err := db.Table(UsersTable).Exist(&user)
	if has {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "mobile has existed", "code": 401})
		log.Println("mobile has existed")
		return
	}
	if err != nil && !has {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "server db" + err.Error(), "code": 501})
		log.Println("server db " + err.Error())
		return
	}

	user.UUID = ulid.Make().String()
	user.Passwd = utils.Sha1En(userMobile.Passwd)
	user.CreatedTime = time.Now().Unix()
	user.UpdatedTime = time.Now().Unix()

	if _, err = db.Table(UsersTable).Insert(user); err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "register failed", "code": 401})
		log.Println("register failed")
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "register success",
		"code":    200,
	})
	return
}

// 判断是否https
func IsHttps(c *gin.Context) bool {
	if c.GetHeader(HEADER_FORWARDED_PROTO) == "https" || c.Request.TLS != nil {
		return true
	}
	return false
}

type CustomClaims struct {
	UserId int64
	jwt.StandardClaims
}

// 产生token
func (cc *CustomClaims) MakeToken() (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, cc)
	return token.SignedString([]byte(SECRETKEY))
}

func Info(c *gin.Context) {
	uuid := c.Param("uuid")
	println(uuid)
	db := c.MustGet("db").(*xorm.Engine)

	user := User{}
	userRow := UserRow{}
	user.UUID = uuid
	has, err := db.Table(UsersTable).Where("uuid=?", user.UUID).Get(&userRow)
	if !has {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "user not exist", "code": 401})
		log.Println("user not exist")
		return
	} else if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "server db" + err.Error(), "code": 501})
		log.Println("server db" + err.Error())
		return
	}
	fmt.Println(userRow)
	//隐藏手机号中间数字
	s := userRow.Mobile
	userRow.Mobile = string([]byte(s)[0:3]) + "****" + string([]byte(s)[6:])
	c.JSON(http.StatusOK, userRow)
	return
}

func DoLogin(c *gin.Context, user User) error {
	secure := IsHttps(c)
	if conf.Cfg.OpenJwt { //返回jwt
		customClaims := &CustomClaims{
			UserId: user.ID,
			StandardClaims: jwt.StandardClaims{
				ExpiresAt: time.Now().Add(time.Duration(MAXAGE) * time.Second).Unix(), // 过期时间，必须设置
			},
		}
		accessToken, err := customClaims.MakeToken()
		if err != nil {
			return err
		}
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
		c.Header(ACCESS_TOKEN, accessToken)
		c.Header(REFRESH_TOKEN, refreshToken)
		c.SetCookie(ACCESS_TOKEN, accessToken, MAXAGE, "/", "", secure, true)
		c.SetCookie(REFRESH_TOKEN, refreshToken, MAXAGE, "/", "", secure, true)
	}
	id := strconv.Itoa(int(user.ID))
	c.SetCookie(COOKIE_TOKEN, id, MAXAGE, "/", "", secure, true)

	return nil
}

// 解析token
func ParseToken(tokenString string) (*CustomClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(SECRETKEY), nil
	})
	if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
		return claims, nil
	} else {
		return nil, err
	}
}
