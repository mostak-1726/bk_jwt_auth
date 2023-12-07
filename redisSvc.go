package auth

import (
	"errors"
	"time"
)

func setStr(key string, value string, ttl int) (err error, status bool) {
	if key == "" || value == "" {
		return errors.New("empty redis key or value"), false
	}

	getRedisClient().Set(key, value, time.Duration(ttl)*time.Second)
	if ttl > 0 {
		err := getRedisClient().Expire(key, time.Duration(ttl)*time.Second).Err()
		if err != nil {
			return err, false
		}
	}

	return nil, true
}

func get(key string) (string, error) {
	if key == "" {
		return "", errors.New("empty redis key or value")
	}

	return getRedisClient().Get(key).Result()
}

func del(keys ...string) error {
	return getRedisClient().Del(keys...).Err()
}

func hasKey(key string) bool {
	return getRedisClient().Exists(key).Val() == 1
}
