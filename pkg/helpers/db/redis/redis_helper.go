package mysql

import (
	"github.com/exgamer/go-rest-sdk/pkg/config/structures"
	"github.com/redis/go-redis/v9"
)

func OpenRedisConnection(redisConfig *structures.RedisConfig) *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr:     redisConfig.RedisHost,
		Password: redisConfig.RedisPassword,
		DB:       redisConfig.RedisDb,
		PoolSize: redisConfig.PoolSize,
	})

	return client
}

func CloseRedisConnection(client *redis.Client) {
	client.Close()
}
