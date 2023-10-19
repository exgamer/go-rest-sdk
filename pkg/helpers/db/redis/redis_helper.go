package redis

import (
	"github.com/gomodule/redigo/redis"
)

func GetString(pool *redis.Pool, key string) (string, error) {
	conn := pool.Get()
	defer conn.Close()

	return redis.String(conn.Do("GET", key))
}

func SetString(pool *redis.Pool, key string, data string, ttl int) (interface{}, error) {
	conn := pool.Get()
	defer conn.Close()

	return conn.Do("SETEX", key, ttl, data)
}
