package auth

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	"github.com/mostak-1726/bk_jwt_auth/config"
	"github.com/mostak-1726/bk_jwt_auth/service"
	"github.com/mostak-1726/bk_jwt_auth/type"
	"net/http"
	"time"
)

type CapAuthIntegrator struct {
	UserName                string
	Password                string
	Request                 _type.AuthTokenRequest
	AuthVerifyReq           _type.AuthTokenVerifyRequest
	JwtTokenExpirySecrete   string
	JwtTokenExpiryInSeconds int
	TestCustomerAppToken    string
	RedisConfig             config.RedisConfig
}

func (ci *CapAuthIntegrator) GenerateAuthToken(c echo.Context) error {
	var err error

	if err = c.Bind(&ci.Request); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Failed to parse request")
	}

	if err = ci.Request.Validate(); err != nil {
		return c.JSON(http.StatusUnprocessableEntity, _type.Response{
			Success: false,
			Message: "Invalid request",
			Errors:  err,
		})
	}
	if ci.Request.Username != ci.UserName || ci.Request.Password != ci.Password {
		return c.JSON(http.StatusUnprocessableEntity, _type.Response{
			Success: false,
			Message: "Invalid userName or Password",
			Errors:  err,
		})
	}

	token := service.GenerateAndStoreAuthToken(ci.Request.MobileNumber, ci.JwtTokenExpiryInSeconds)
	updateTime := time.Now().Format("02-01-2006 15:04:05")

	return c.JSON(http.StatusOK, _type.Response{
		Success: true,
		Message: "Auth token created successfully",
		Data: _type.AuthTokenResponseData{
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
		return c.JSON(http.StatusUnprocessableEntity, _type.Response{
			Success: false,
			Message: "Invalid request",
			Errors:  err,
		})
	}

	jte := ci.JwtTokenExpiryInSeconds

	expiresAt := time.Now().Add(time.Duration(jte) * time.Second)

	if ci.TestCustomerAppToken != "" && ci.TestCustomerAppToken == ci.AuthVerifyReq.Token {
		return verifyTestAuthToken(c, ci.AuthVerifyReq, expiresAt, ci.JwtTokenExpirySecrete)
	}

	vt := service.VerifyAuthTokenService(ci.AuthVerifyReq.Token)
	if !vt {
		return c.JSON(http.StatusBadRequest, _type.Response{
			Success: false,
			Message: "Invalid token",
		})
	}

	wNumber := service.GetWalletNumberAndRemoveToken(ci.AuthVerifyReq.Token)
	token, er := service.GenerateJwt(wNumber, expiresAt, ci.JwtTokenExpirySecrete)
	if er != nil {
		log.Error("Failed to generate jwt token", er)
		return c.JSON(http.StatusBadRequest, _type.Response{
			Success: false,
			Message: "Failed to generate jwt token",
		})
	}

	return c.JSON(http.StatusOK, _type.Response{
		Success: true,
		Message: "Auth token verified successfully",
		Data: _type.AccessTokenResponseData{
			AccessToken: token,
			ExpiresAt:   expiresAt,
		},
	})
}

func verifyTestAuthToken(c echo.Context, req _type.AuthTokenVerifyRequest, expiresAt time.Time, secrete string) error {
	errs := validation.Errors{
		"mobile_number": validation.Validate(req.MobileNumber, validation.Required, validation.Length(11, 11)),
	}.Filter()

	if errs != nil {
		return c.JSON(http.StatusUnprocessableEntity, _type.Response{
			Success: false,
			Message: "Invalid request",
			Errors:  errs,
		})
	}

	token, er := service.GenerateJwt(req.MobileNumber, expiresAt, secrete)
	if er != nil {
		if er != nil {
			log.Error("Failed to generate jwt token", er)
			return c.JSON(http.StatusBadRequest, _type.Response{
				Success: false,
				Message: "Failed to generate jwt token",
			})
		}

	}

	return c.JSON(http.StatusOK, _type.Response{
		Success: true,
		Message: "Auth token verified successfully",
		Data: _type.AccessTokenResponseData{
			AccessToken: token,
			ExpiresAt:   expiresAt,
		},
	})
}

func NewCapAuthIntegrator(userName, password string, expiryInSec int, testToken, secrete string) CapAuthIntegrator {
	return CapAuthIntegrator{
		UserName:                userName,
		Password:                password,
		JwtTokenExpiryInSeconds: expiryInSec,
		TestCustomerAppToken:    testToken,
		JwtTokenExpirySecrete:   secrete,
	}
}