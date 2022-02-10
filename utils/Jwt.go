package utils

import (
	"admin/setting"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"time"
)

type JsonWebToken struct {

}

type CustomClaims struct {
	UserId int
	Username string
	jwt.StandardClaims
}

func (jw *JsonWebToken)  CreateToken(uid int, username string, maxAge int) string {
	claims :=  &CustomClaims {
		UserId: uid,
		Username: username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Duration(maxAge)*time.Second).Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err:= token.SignedString([]byte(setting.JwtSecret))
	if err != nil {
		return ""
	}

	return tokenString
}

func (jw *JsonWebToken)  VerifyToken(tokenString string) (*CustomClaims, bool) {
	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(setting.JwtSecret), nil
	})

	if err != nil || token == nil {
		return nil,  false
	}

	if customClaims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
		return customClaims, true
	}

	return nil, false
}
