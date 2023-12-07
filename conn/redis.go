package conn

import (
	"github.com/go-redis/redis"
	_type "github.com/mostak-1726/bk_jwt_auth/type"
	"log"
)

var redisClient *redis.Client

func ConnectRedis(conf _type.RedisConfig) {

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
	return redisClient
}
