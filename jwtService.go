package auth

import (
	"github.com/golang-jwt/jwt/v5"
	"time"
)

type jwtCustomClaims struct {
	MobileNumber string `json:"wallet_number"`
	jwt.RegisteredClaims
}

func generateJwt(wNumber string, exp time.Time, secrete string) (string, error) {

	claims := &jwtCustomClaims{
		wNumber,
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(exp),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	t, err := token.SignedString([]byte(secrete))
	if err != nil {
		return "", err
	}

	return t, nil
}
