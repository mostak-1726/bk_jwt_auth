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

func main() {
	var e = echo.New()
	config.LoadConfig()
	conn.ConnectRedis()
	registerRoutes(e)
	// echo middlewares
	e.Use(middleware.CORS())
	e.Use(middleware.Secure())
	e.Use(middleware.Recover())

	// remove trailing slashes from each requests
	e.Pre(middleware.RemoveTrailingSlash())
	e.Use(echojwt.WithConfig(getEchoJwtConfig()))

	e.Logger.Fatal(e.Start(":" + "8090"))
}

func getEchoJwtConfig() echojwt.Config {

	jts := "AwesomeTokenSecret"

	return echojwt.Config{
		Skipper: func(c echo.Context) bool {
			route := []string{
				"/bkash/auth",
				"bkash/auth/verify",
			}
			return slices.Contains(route, c.Request().URL.Path)
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
func registerRoutes(e *echo.Echo) {
	config.LoadConfig()
	conn.ConnectRedis()
	handler := auth.NewCapAuthIntegrator("mostak", "12345", 3600, "14580760-b5d9-42d7-aa3a-51d20caeff6a", "testSecret")
	e.POST("/bkash/auth", handler.GenerateAuthToken)
	e.POST("/bkash/auth/verify", handler.VerifyAuthToken)

}
