package redis

import (
	"context"
	"github.com/redis/go-redis/v9"
	"github.com/shiroyaavish/go-common/config"
)

var ctx = context.Background()

type RedisClient struct {
	client *redis.Client
}

func (r *RedisClient) Raw() *redis.Client {
	return r.client
}

func NewRedisClient(cfg *config.RedisConfig) *RedisClient {
	rdb := redis.NewClient(&redis.Options{
		Addr:     cfg.Host,
		Password: cfg.Password, // no password set
		DB:       cfg.DBNumber, // use default DB
	})

	return &RedisClient{client: rdb}
}
