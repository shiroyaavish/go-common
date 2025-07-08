package fiber_cache

import "github.com/go-redis/redis/v8"

var client *redis.Client

func SetClient(redisClient *redis.Client) {
	client = redisClient
}

func GetClient() *redis.Client {
	return client
}
