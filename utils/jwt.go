package utils

import (
	"errors"
	"github.com/dgrijalva/jwt-go"
	"github.com/wonderivan/logger"
)

var JWTToken jwtToken

type jwtToken struct{}

type CustomClaims struct {
	UserName string `json:"username"`
	Password string `json:"password"`
	jwt.StandardClaims
}

const (
	SECRET = "devops"
)

func (*jwtToken) ParseToken(tokenString string) (claims *CustomClaims, err error) {
	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(SECRET), nil
	})
	if err != nil {
		logger.Error("parse token failed ", err)
		if ve, ok := err.(*jwt.ValidationError); ok {
			if ve.Errors&jwt.ValidationErrorMalformed != 0 {
				return nil, errors.New("TokenMalformed")
			} else if ve.Errors&jwt.ValidationErrorExpired != 0 {
				// Token is expired
				return nil, errors.New("TokenExpired")
			} else if ve.Errors&jwt.ValidationErrorNotValidYet != 0 {
				return nil, errors.New("TokenNotValidYet")
			} else {
				return nil, errors.New("TokenInvalid")
			}
		}
	}
	if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
		return claims, nil
	}
	return nil, errors.New("解析Token失败")
}
