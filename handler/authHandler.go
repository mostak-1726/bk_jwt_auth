package handler

import (
	"fmt"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/labstack/echo/v4"
	"github.com/mostak-1726/bk_jwt_auth/config"
	"github.com/mostak-1726/bk_jwt_auth/consts"
	"github.com/mostak-1726/bk_jwt_auth/service"
	"github.com/mostak-1726/bk_jwt_auth/type"
	"net/http"
	"time"
)

func GenerateAuthToken(c echo.Context) error {
	var req _type.AuthTokenRequest
	var err error

	if err = c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Failed to parse request")
	}

	if err = req.Validate(); err != nil {
		return c.JSON(http.StatusUnprocessableEntity, _type.Response{
			Success: false,
			Message: "Invalid request",
			Errors:  err,
		})
	}

	token := service.GenerateAndStoreAuthToken(req.MobileNumber)
	updateTime := time.Now().Format(consts.DateTimeFormat)

	return c.JSON(http.StatusOK, _type.Response{
		Success: true,
		Message: "Auth token created successfully",
		Data: _type.AuthTokenResponseData{
			IdToken:    token,
			UpdateTime: updateTime,
		},
	})
}

func VerifyAuthToken(c echo.Context) error {
	var req _type.AuthTokenVerifyRequest
	var err error

	if err = c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Failed to parse request")
	}

	if err = req.Validate(); err != nil {
		return c.JSON(http.StatusUnprocessableEntity, _type.Response{
			Success: false,
			Message: "Invalid request",
			Errors:  err,
		})
	}

	jte := config.App().JwtTokenExpiryInSeconds
	if jte == 0 {
		jte = consts.JwtTokenExpiryInSeconds
	}
	expiresAt := time.Now().Add(time.Duration(jte) * time.Second)

	tt := config.App().TestCustomerAppToken
	if tt != "" && tt == req.Token {
		return VerifyTestAuthToken(c, req, expiresAt)
	}

	vt := service.VerifyAuthTokenService(req.Token)
	if !vt {
		return c.JSON(http.StatusBadRequest, _type.Response{
			Success: false,
			Message: "Invalid token",
		})
	}

	wNumber := service.GetWalletNumberAndRemoveToken(req.Token)
	token, er := service.GenerateJwt(wNumber, expiresAt)
	if er != nil {
		fmt.Println(er)
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

func VerifyTestAuthToken(c echo.Context, req _type.AuthTokenVerifyRequest, expiresAt time.Time) error {
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

	token, er := service.GenerateJwt(req.MobileNumber, expiresAt)
	if er != nil {
		fmt.Println(er)
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
