package main

import (
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	auth "github.com/mostak-1726/bk_jwt_auth"
	"github.com/mostak-1726/bk_jwt_auth/config"
	"github.com/mostak-1726/bk_jwt_auth/conn"
	"golang.org/x/exp/slices"
	"net/http"
)

var e = echo.New()

func main() {
	config.LoadConfig()
	conn.ConnectRedis()
	auth.RegisterRoutes(e)
	// echo middlewares
	e.Use(middleware.CORS())
	e.Use(middleware.Secure())
	e.Use(middleware.Recover())

	// remove trailing slashes from each requests
	e.Pre(middleware.RemoveTrailingSlash())
	e.Use(echojwt.WithConfig(getEchoJwtConfig()))

	port := config.App().Port
	e.Logger.Fatal(e.Start(":" + port))
}

func getEchoJwtConfig() echojwt.Config {

	jts := "AwesomeTokenSecret"

	return echojwt.Config{
		Skipper: func(c echo.Context) bool {
			return slices.Contains(config.App().JwtSkipper, c.Request().URL.Path)
		},
		SigningKey: []byte(jts),
		ErrorHandler: func(c echo.Context, err error) error {
			fmt.Println(c.Request().URL)
			return c.JSON(http.StatusUnauthorized, "{Success: false,Message: `Invalid or expired token`}")
		},
		SuccessHandler: func(c echo.Context) {
			token := c.Get("user").(*jwt.Token)
			claims, _ := token.Claims.(jwt.MapClaims)
			wNumber := claims["wallet_number"].(string)
			c.Set("wallet_number", wNumber)
		},
	}
}
