package redis

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/mostak-1726/bk_jwt_auth/conn"
	"strconv"
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

func Set(key string, value interface{}, ttl int) error {
	if key == "" || value == "" {
		return errors.New("empty redis key or value")
	}

	serializedValue, err := json.Marshal(value)
	if err != nil {
		return err
	}

	conn.Redis().Set(key, string(serializedValue), time.Duration(ttl)*time.Second)
	if ttl > 0 {
		conn.Redis().Expire(key, time.Duration(ttl)*time.Second).Err()
	}
	return err
}

func Get(key string) (string, error) {
	if key == "" {
		return "", errors.New("empty redis key or value")
	}

	return conn.Redis().Get(key).Result()
}

func GetBoolean(key string) (bool, error) {
	if key == "" {
		return false, errors.New("empty redis key or value")
	}

	str, err := conn.Redis().Get(key).Result()
	if err != nil {
		return false, err
	}

	return strconv.ParseBool(str)
}

func GetInt(key string) (int, error) {
	if key == "" {
		return 0, errors.New("empty redis key or value")
	}

	str, err := conn.Redis().Get(key).Result()
	if err != nil {
		return 0, err
	}

	return strconv.Atoi(str)
}

func GetStruct(key string, outputStruct interface{}) error {
	if key == "" {
		return errors.New("empty redis key or value")
	}

	serializedValue, err := conn.Redis().Get(key).Result()
	if err != nil {
		return err
	}
	if err := json.Unmarshal([]byte(serializedValue), &outputStruct); err != nil {
		return err
	}

	return nil
}

func Del(keys ...string) error {
	return conn.Redis().Del(keys...).Err()
}

func DelPattern(pattern string) error {
	iter := conn.Redis().Scan(0, pattern, 0).Iterator()

	for iter.Next() {
		err := conn.Redis().Del(iter.Val()).Err()
		if err != nil {
			return err
		}
	}

	if err := iter.Err(); err != nil {
		return err
	}

	return nil
}

func HasKey(key string) bool {
	return conn.Redis().Exists(key).Val() == 1
}

func IsExpire(key string) bool {
	fmt.Println(conn.Redis().TTL(key).Val())
	return conn.Redis().TTL(key).Val() < 0
}

func IncCount(key string, ttl int) error {
	conn.Redis().Expire(key, time.Duration(ttl)*time.Second)
	return conn.Redis().HIncrBy(key, "counter", 1).Err()
}

func SetCounter(token string, limit int, ttl int) error {
	hash := map[string]interface{}{
		"counter": 0,
		"limit":   limit,
	}
	conn.Redis().HMSet(token, hash)
	return conn.Redis().Expire(token, time.Duration(ttl)*time.Second).Err()
}

func GetSet(key string, value *string) (string, error) {
	if key == "" {
		return "", errors.New("empty redis key or value")
	}

	return conn.Redis().GetSet(key, value).Result()
}
