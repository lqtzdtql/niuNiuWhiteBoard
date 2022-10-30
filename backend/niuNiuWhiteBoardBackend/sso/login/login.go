package login

import (
	"fmt"
	"github.com/gin-gonic/gin"
	conf "niuNiuWhiteBoardBackend/sso/config"
	"niuNiuWhiteBoardBackend/sso/models"
	"niuNiuWhiteBoardBackend/sso/response"
	"niuNiuWhiteBoardBackend/utils/common"
	"strconv"
	"time"

	"github.com/dgrijalva/jwt-go"
)

const (
	ACCESS_TOKEN           = "Access_Token"
	REFRESH_TOKEN          = "Refresh_Token"
	COOKIE_TOKEN           = "UserId"
	HEADER_AUTH            = "Authorization"
	HEADER_ETAG_SERVER     = "ETag"
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

// 手机密码
func Login(c *gin.Context) {
	var userMobile UserMobilePasswd
	if err := c.BindJSON(&userMobile); err != nil {
		response.ShowError(c, "bindJSON err")
		return
	}
	model := models.Users{Mobile: userMobile.Mobile}
	if has := model.GetRow(); !has {
		response.ShowError(c, "mobile_not_exists")
		return
	}
	if common.Sha1En(userMobile.Passwd) != model.Passwd {
		response.ShowError(c, "login_error")
		return
	}
	err := DoLogin(c, model)
	if err != nil {
		response.ShowError(c, "fail")
		return
	}
	response.ShowSuccess(c, "success")
	return
}

// 注销登录
func Logout(c *gin.Context) {
	secure := IsHttps(c)
	//access_token  refresh_token 加黑名单

	c.SetCookie(COOKIE_TOKEN, "", -1, "/", "", secure, true)
	c.SetCookie(ACCESS_TOKEN, "", -1, "/", "", secure, true)
	c.SetCookie(REFRESH_TOKEN, "", -1, "/", "", secure, true)
	response.ShowSuccess(c, "success")
	return
}

// 手机号注册
func SignupByMobile(c *gin.Context) {
	var userMobile UserMobile
	if err := c.BindJSON(&userMobile); err != nil {
		response.ShowError(c, "BindJSON err")
		return
	}
	model := models.Users{Mobile: userMobile.Mobile}
	if has := model.GetRow(); has {
		response.ShowError(c, "mobile_exists")
		return
	}

	model.Passwd = common.Sha1En(userMobile.Passwd)
	model.CreatedTime = time.Now().Unix()
	model.UpdatedTime = time.Now().Unix()

	traceModel := models.Trace{CreatedTime: model.CreatedTime}
	traceModel.Ip = common.IpStringToInt(GetClientIp(c))

	deviceModel := models.Device{CreatedTime: model.CreatedTime, Ip: traceModel.Ip, Client: c.GetHeader("User-Agent")}
	_, err := model.Add(&traceModel, &deviceModel)
	if err != nil {
		fmt.Println(err)
		response.ShowError(c, "fail")
		return
	}
	response.ShowSuccess(c, "success")
	return
}

func DoLogin(c *gin.Context, user models.Users) error {
	secure := IsHttps(c)
	if conf.Cfg.OpenJwt { //返回jwt
		customClaims := &CustomClaims{
			UserId: user.Id,
			StandardClaims: jwt.StandardClaims{
				ExpiresAt: time.Now().Add(time.Duration(MAXAGE) * time.Second).Unix(), // 过期时间，必须设置
			},
		}
		accessToken, err := customClaims.MakeToken()
		if err != nil {
			return err
		}
		refreshClaims := &CustomClaims{
			UserId: user.Id,
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
	id := strconv.Itoa(int(user.Id))
	c.SetCookie(COOKIE_TOKEN, id, MAXAGE, "/", "", secure, true)

	return nil
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

func GetClientIp(c *gin.Context) string {
	ip := c.ClientIP()
	if ip == "::1" {
		ip = "127.0.0.1"
	}
	return ip
}

// 产生token
func (cc *CustomClaims) MakeToken() (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, cc)
	return token.SignedString([]byte(SECRETKEY))
}

func Info(c *gin.Context) {
	uid := c.MustGet("uid").(int64)
	fmt.Println(uid)
	model := models.Users{}
	model.Id = uid
	row, err := model.GetRowById()
	if err != nil {
		fmt.Println(err)
		response.ShowError(c, err.Error())
		return
	}
	fmt.Println(row)
	fmt.Println(row.Name)
	//隐藏手机号中间数字
	s := row.Mobile
	row.Mobile = string([]byte(s)[0:3]) + "****" + string([]byte(s)[6:])
	response.ShowData(c, row)
	return
}

//// 解析token
//func ParseToken(tokenString string) (*CustomClaims, error) {
//	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
//		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
//			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
//		}
//		return []byte(SECRETKEY), nil
//	})
//	if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
//		return claims, nil
//	} else {
//		return nil, err
//	}
//}
//
//
//func GetParam(c *gin.Context,key string)(string,bool){
//	val:=c.GetHeader(key)
//	if val!=""{
//		return val,true
//	}
//	val,err :=c.Cookie(key)
//	if err!=nil{
//		return "",false
//	}
//	return val,true
//}
