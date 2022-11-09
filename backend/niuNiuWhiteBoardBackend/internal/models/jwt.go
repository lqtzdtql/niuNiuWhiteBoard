package models

import (
	"errors"
	"niuNiuWhiteBoardBackend/common/log"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

type CustomClaims struct {
	UserId int64
	jwt.StandardClaims
}

// 解析token
func ParseToken(tokenString string) (*CustomClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			log.Logger.Info("Unexpected signing method", log.Any("Unexpected signing method", token.Header["alg"]))
			return nil, errors.New("unexpected signing method")
		}
		return []byte(SECRETKEY), nil
	})
	if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
		return claims, nil
	} else {
		return nil, err
	}
}

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

// 产生token
func (cc *CustomClaims) MakeToken() (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, cc)
	return token.SignedString([]byte(SECRETKEY))
}

// 判断是否https
func IsHttps(c *gin.Context) bool {
	if c.GetHeader(HEADER_FORWARDED_PROTO) == "https" || c.Request.TLS != nil {
		return true
	}
	return false
}
