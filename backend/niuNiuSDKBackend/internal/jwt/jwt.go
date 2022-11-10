package jwt

import (
	"errors"

	"niuNiuSDKBackend/common/log"

	"github.com/dgrijalva/jwt-go"
)

const (
	SECRETKEY = "42wqTE23123wffLU94342wgldgFs"
	MAXAGE    = 3600 * 24
)

type ClientClaims struct {
	SK         string
	RoomName   string
	UserName   string
	Permission string
	jwt.StandardClaims
}

// 解析token
func ParseToken(tokenString string) (*ClientClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &ClientClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			log.Logger.Info("Unexpected signing method", log.Any("Unexpected signing method", token.Header["alg"]))
			return nil, errors.New("unexpected signing method")
		}
		return []byte(SECRETKEY), nil
	})
	if claims, ok := token.Claims.(*ClientClaims); ok && token.Valid {
		return claims, nil
	} else {
		return nil, err
	}
}
