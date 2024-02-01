package container

import (
	"concurrent-image-fetcher/config"
	"database/sql"
	"github.com/go-redis/redis/v9"
)

type Container struct {
	Redis  *redis.Client
	Db     *sql.DB
	Config *config.Config
}

func NewContainer(config *config.Config) *Container {
	return &Container{
		Redis:  NewRedis(config.Cache),
		Db:     NewMysql(config.Db),
		Config: config,
	}
}
