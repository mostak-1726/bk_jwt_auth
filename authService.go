package auth

import (
	"github.com/google/uuid"
)

func generateAndStoreAuthToken(uNumber string, expiry int) string {

	token := uuid.NewString()
	err, _ := setStr(token, uNumber, expiry)

	if err != nil {
		return ""
	}

	return token
}

func verifyAuthTokenService(token string) bool {
	return hasKey(token)
}

func getWalletNumberAndRemoveToken(token string) string {
	value, ok := get(token)
	if ok != nil {
		return ""
	}

	err := del(token)
	if err != nil {
		return ""
	}

	return value
}
