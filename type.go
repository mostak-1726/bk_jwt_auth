package auth

import (
	v "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-redis/redis"
	"time"
)

type AuthTokenResponseData struct {
	IdToken    string `json:"id_token"`
	UpdateTime string `json:"update_time"`
}
type Response struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	Errors  error       `json:"errors,omitempty"`
}
type AccessTokenResponseData struct {
	AccessToken string    `json:"access_token"`
	ExpiresAt   time.Time `json:"expires_at"`
}

type AuthTokenRequest struct {
	Username     string `json:"username"`
	Password     string `json:"password"`
	MobileNumber string `json:"mobile_number"`
}

func (q AuthTokenRequest) Validate() error {
	return v.ValidateStruct(&q,
		v.Field(&q.Username, v.Required),
		v.Field(&q.Password, v.Required),
		v.Field(&q.MobileNumber, v.Required, v.Length(11, 11)),
	)
}

type AuthTokenVerifyRequest struct {
	Token        string `json:"token"`
	MobileNumber string `json:"mobile_number,omitempty"`
}

func (q AuthTokenVerifyRequest) Validate() error {
	return v.ValidateStruct(&q,
		v.Field(&q.Token, v.Required),
	)
}

type Config struct {
	UserName             string
	Password             string
	ExpiryInSec          int
	JwtTokenSecrete      string
	TestCustomerAppToken string
	RedisClient          *redis.Client
}
