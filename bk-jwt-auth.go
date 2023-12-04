package auth

import (
	"github.com/labstack/echo/v4"
	"github.com/mostak-1726/bk_jwt_auth/config"
	"github.com/mostak-1726/bk_jwt_auth/conn"
	"github.com/mostak-1726/bk_jwt_auth/handler"
)

func RegisterRoutes(e *echo.Echo) {
	config.LoadConfig()
	conn.ConnectRedis()

	e.POST(config.App().TokenGenerationRoute, handler.GenerateAuthToken)
	e.POST(config.App().TokenVerifyRoute, handler.VerifyAuthToken)

}
func GetAuthSkippers() []string {
	return config.App().JwtSkipper
}
