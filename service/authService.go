package service

import (
	"github.com/google/uuid"
	"github.com/mostak-1726/bk_jwt_auth/redis"
)

func GenerateAndStoreAuthToken(uNumber string, expiry int) string {

	token := uuid.NewString()
	err, _ := redis.SetStr(token, uNumber, expiry)

	if err != nil {
		return ""
	}

	return token
}

func VerifyAuthTokenService(token string) bool {
	return redis.HasKey(token)
}

func GetWalletNumberAndRemoveToken(token string) string {
	value, ok := redis.Get(token)
	if ok != nil {
		return ""
	}

	err := redis.Del(token)
	if err != nil {
		return ""
	}

	return value
}
