package auth

import (
	"github.com/go-redis/redis"
	"log"
)

var redisClient *redis.Client

func connectRedis(conf RedisConfig) {

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

func getRedisClient() *redis.Client {
	return redisClient
}
