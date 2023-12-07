package auth

import (
	"github.com/go-redis/redis"
)

var redisClient *redis.Client

func setRedisClient(c *redis.Client) {
	redisClient = c
}
func getRedisClient() *redis.Client {
	return redisClient
}
