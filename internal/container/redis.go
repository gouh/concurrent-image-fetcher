package container

import (
	"concurrent-image-fetcher/config"
	"github.com/go-redis/redis/v9"
)

func NewRedis(config *config.CacheConfig) *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr: config.DSN,
		DB:   0,
	})
}
