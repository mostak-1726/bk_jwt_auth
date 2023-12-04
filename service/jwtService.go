package service

import (
	"github.com/golang-jwt/jwt/v5"
	"github.com/mostak-1726/bk_jwt_auth/config"
	"github.com/mostak-1726/bk_jwt_auth/consts"
	"time"
)

type jwtCustomClaims struct {
	MobileNumber string `json:"wallet_number"`
	jwt.RegisteredClaims
}

func GenerateJwt(wNumber string, exp time.Time) (string, error) {
	jts := config.App().JwtTokenSecret
	if jts == "" {
		jts = consts.JwtTokenSecret
	}

	claims := &jwtCustomClaims{
		wNumber,
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(exp),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	t, err := token.SignedString([]byte(jts))
	if err != nil {
		return "", err
	}

	return t, nil
}
