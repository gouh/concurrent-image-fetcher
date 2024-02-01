package config

import (
	"fmt"
	"os"
)

type CacheConfig struct {
	DSN string
}

func validateCacheConfig() error {
	if os.Getenv("REDIS_DSN") == "" {
		return fmt.Errorf("REDIS_DSN config is not set")
	}
	return nil
}

func LoadCacheConfig() (*CacheConfig, error) {
	err := validateCacheConfig()

	if err != nil {
		return nil, err
	}

	return &CacheConfig{
		DSN: os.Getenv("REDIS_DSN"),
	}, err
}
