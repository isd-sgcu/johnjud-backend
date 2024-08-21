package database

import (
	"fmt"

	"github.com/isd-sgcu/johnjud-gateway/config"
	"github.com/pkg/errors"
	"github.com/redis/go-redis/v9"
)

func InitRedisConnection(conf *config.Redis) (*redis.Client, error) {
	addr := fmt.Sprintf("%s:%d", conf.Host, conf.Port)

	cache := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: conf.Password,
	})

	if cache == nil {
		return nil, errors.New("Failed to connect to redis server")
	}

	return cache, nil
}
