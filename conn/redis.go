package conn

import (
	"github.com/go-redis/redis"
	"github.com/mostak-1726/bk_jwt_auth/config"
	"log"
)

var redisClient *redis.Client

func ConnectRedis() {
	conf := config.Redis()

	log.Print("connecting to redis at ", conf.Host, ":", conf.Port, "...")

	redisClient = redis.NewClient(&redis.Options{
		Addr:     conf.Host + ":" + conf.Port,
		Password: conf.Pass,
		DB:       conf.Db,
	})

	if _, err := redisClient.Ping().Result(); err != nil {
		log.Print("failed to connect redis: ", err)
		panic(err)
	}

	log.Print("redis connection successful...")
}

func Redis() *redis.Client {
	if redisClient == nil {
		ConnectRedis()
	}
	return redisClient
}
