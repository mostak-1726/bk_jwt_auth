package redis

import (
	"errors"
	"github.com/mostak-1726/bk_jwt_auth/conn"
	"time"
)

func SetStr(key string, value string, ttl int) (err error, status bool) {
	if key == "" || value == "" {
		return errors.New("empty redis key or value"), false
	}

	conn.Redis().Set(key, value, time.Duration(ttl)*time.Second)
	if ttl > 0 {
		err := conn.Redis().Expire(key, time.Duration(ttl)*time.Second).Err()
		if err != nil {
			return err, false
		}
	}

	return nil, true
}

func Get(key string) (string, error) {
	if key == "" {
		return "", errors.New("empty redis key or value")
	}

	return conn.Redis().Get(key).Result()
}

func Del(keys ...string) error {
	return conn.Redis().Del(keys...).Err()
}

func HasKey(key string) bool {
	return conn.Redis().Exists(key).Val() == 1
}
