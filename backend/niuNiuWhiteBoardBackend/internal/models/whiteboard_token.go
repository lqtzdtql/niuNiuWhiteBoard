package models

import "github.com/dgrijalva/jwt-go"

// 产生白板token
type ClientClaims struct {
	SK         string
	RoomName   string
	UserName   string
	Permission string
	jwt.StandardClaims
}

func (cc *ClientClaims) MakeWhiteBoardToken() (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, cc)
	return token.SignedString([]byte(SECRETKEY))
}
