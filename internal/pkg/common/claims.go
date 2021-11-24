package common

import "github.com/golang-jwt/jwt"

type Claims struct {
	Uid int
	jwt.StandardClaims
}
