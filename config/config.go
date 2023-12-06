package config

import (
	"encoding/json"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/spf13/viper"
	_ "github.com/spf13/viper/remote"
	"golang.org/x/exp/slices"
	"log"
	"net/http"
)

type RedisConfig struct {
	Host string
	Port string
	Pass string
	Db   int
	Ttl  int // seconds
}

type Config struct {
	Redis *RedisConfig
}

var config Config

func Redis() *RedisConfig {
	return config.Redis
}

func LoadConfig() {
	setDefaultConfig()

	_ = viper.BindEnv("consul_url")
	_ = viper.BindEnv("consul_path")

	consulURL := viper.GetString("consul_url")
	consulPath := viper.GetString("consul_path")

	if consulURL != "" && consulPath != "" {
		_ = viper.AddRemoteProvider("consul", consulURL, consulPath)

		viper.SetConfigType("json")
		err := viper.ReadRemoteConfig()

		if err != nil {
			log.Println(fmt.Sprintf("%s named \"%s\"", err.Error(), consulPath))
		}

		config = Config{}

		if err := viper.Unmarshal(&config); err != nil {
			panic(err)
		}

		if r, err := json.MarshalIndent(&config, "", "  "); err == nil {
			fmt.Println(string(r))
		}
	} else {
		log.Println("CONSUL_URL or CONSUL_PATH missing! Serving with default config...")
	}
}

func setDefaultConfig() {
	config.Redis = &RedisConfig{
		Host: "127.0.0.1",
		Port: "6379",
		Pass: "secret_redis",
		Db:   1,
		Ttl:  3600,
	}
}

func GetEchoJwtConfig() echojwt.Config {
	jts := "AwesomeTokenSecret"

	return echojwt.Config{
		Skipper: func(c echo.Context) bool {
			authSkipPaths := []string{
				"/bkash/auth",
				"/bkash/auth/verify",
			}
			return slices.Contains(authSkipPaths, c.Request().URL.Path)
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
