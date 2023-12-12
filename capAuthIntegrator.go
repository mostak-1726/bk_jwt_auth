package auth

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-redis/redis"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	"net/http"
	"time"
)

type CapAuthIntegrator struct {
	UserName                string
	Password                string
	Request                 AuthTokenRequest
	AuthVerifyReq           AuthTokenVerifyRequest
	JwtTokenSecrete         string
	JwtTokenExpiryInSeconds int
	TestCustomerAppToken    string
	RedisClient             *redis.Client
}

func (ci *CapAuthIntegrator) GenerateAuthToken(c echo.Context) error {
	var err error

	if err = c.Bind(&ci.Request); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Failed to parse request")
	}

	if err = ci.Request.Validate(); err != nil {
		return c.JSON(http.StatusUnprocessableEntity, Response{
			Success: false,
			Message: "Invalid request",
			Errors:  err,
		})
	}
	if ci.Request.Username != ci.UserName || ci.Request.Password != ci.Password {
		return c.JSON(http.StatusUnprocessableEntity, Response{
			Success: false,
			Message: "Invalid userName or Password",
			Errors:  err,
		})
	}

	token := generateAndStoreAuthToken(ci.Request.MobileNumber, ci.JwtTokenExpiryInSeconds)
	updateTime := time.Now().Format("02-01-2006 15:04:05")

	return c.JSON(http.StatusOK, Response{
		Success: true,
		Message: "Auth token created successfully",
		Data: AuthTokenResponseData{
			IdToken:    token,
			UpdateTime: updateTime,
		},
	})
}

func (ci *CapAuthIntegrator) VerifyAuthToken(c echo.Context) error {
	var err error

	if err = c.Bind(&ci.AuthVerifyReq); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Failed to parse request")
	}

	if err = ci.AuthVerifyReq.Validate(); err != nil {
		return c.JSON(http.StatusUnprocessableEntity, Response{
			Success: false,
			Message: "Invalid request",
			Errors:  err,
		})
	}

	jte := ci.JwtTokenExpiryInSeconds

	expiresAt := time.Now().Add(time.Duration(jte) * time.Second)

	if ci.TestCustomerAppToken != "" && ci.TestCustomerAppToken == ci.AuthVerifyReq.Token {
		return verifyTestAuthToken(c, ci.AuthVerifyReq, expiresAt, ci.JwtTokenSecrete)
	}

	vt := verifyAuthTokenService(ci.AuthVerifyReq.Token)
	if !vt {
		return c.JSON(http.StatusBadRequest, Response{
			Success: false,
			Message: "Invalid token",
		})
	}

	wNumber := getWalletNumberAndRemoveToken(ci.AuthVerifyReq.Token)
	token, er := generateJwt(wNumber, expiresAt, ci.JwtTokenSecrete)
	if er != nil {
		log.Error("Failed to generate jwt token", er)
		return c.JSON(http.StatusBadRequest, Response{
			Success: false,
			Message: "Failed to generate jwt token",
		})
	}

	return c.JSON(http.StatusOK, Response{
		Success: true,
		Message: "Auth token verified successfully",
		Data: AccessTokenResponseData{
			AccessToken: token,
			ExpiresAt:   expiresAt,
		},
	})
}

func verifyTestAuthToken(c echo.Context, req AuthTokenVerifyRequest, expiresAt time.Time, secrete string) error {
	errs := validation.Errors{
		"mobile_number": validation.Validate(req.MobileNumber, validation.Required, validation.Length(11, 11)),
	}.Filter()

	if errs != nil {
		return c.JSON(http.StatusUnprocessableEntity, Response{
			Success: false,
			Message: "Invalid request",
			Errors:  errs,
		})
	}

	token, er := generateJwt(req.MobileNumber, expiresAt, secrete)
	if er != nil {
		if er != nil {
			log.Error("Failed to generate jwt token", er)
			return c.JSON(http.StatusBadRequest, Response{
				Success: false,
				Message: "Failed to generate jwt token",
			})
		}

	}

	return c.JSON(http.StatusOK, Response{
		Success: true,
		Message: "Auth token verified successfully",
		Data: AccessTokenResponseData{
			AccessToken: token,
			ExpiresAt:   expiresAt,
		},
	})
}

func NewCapAuthIntegrator(config Config) CapAuthIntegrator {
	setRedisClient(config.RedisClient)
	return CapAuthIntegrator{
		UserName:                config.UserName,
		Password:                config.Password,
		JwtTokenExpiryInSeconds: config.ExpiryInSec,
		TestCustomerAppToken:    config.TestCustomerAppToken,
		JwtTokenSecrete:         config.JwtTokenSecret,
	}
}
