package client

import (
	"github.com/130-133/go-common/utils/redis"
	"time"
)

type RedisLocker struct {
	client *redis.MRedis
}

func NewRedisLocker(c *redis.MRedis) *RedisLocker {
	return &RedisLocker{client: c}
}

func (c *RedisLocker) Set(key string, value interface{}, expiration time.Duration) (bool, error) {
	return c.client.SetNX(key, value, expiration).Result()
}
func (c *RedisLocker) Get(key string) string {
	return c.client.Get(key).Val()
}
func (c *RedisLocker) Del(key string) error {
	return c.client.Del(key).Err()
}
