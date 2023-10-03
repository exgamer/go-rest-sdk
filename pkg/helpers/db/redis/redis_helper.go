package redis

import (
	"github.com/exgamer/go-rest-sdk/pkg/config/structures"
	"github.com/redis/go-redis/v9"
)

func OpenRedisConnection(redisConfig *structures.RedisConfig) *redis.Client {
	poolSize := redisConfig.PoolSize

	client := redis.NewClient(&redis.Options{
		Addr:     redisConfig.RedisHost,
		Password: redisConfig.RedisPassword,
		DB:       redisConfig.RedisDb,
		PoolSize: poolSize,
	})

	return client
}

func CloseRedisConnection(client *redis.Client) {
	client.Close()
}
