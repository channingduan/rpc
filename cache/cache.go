package cache

import (
	"context"
	"github.com/channingduan/rpc/config"
	"github.com/go-redis/redis/v8"
	"time"
)

type Cache struct {
	cache *redis.Client
}

func Register(config *config.CacheConfig) *Cache {
	return &Cache{
		cache: redis.NewClient(&redis.Options{
			Addr:     config.Addr,
			Username: config.Username,
			Password: config.Password,
			DB:       0,
		}),
	}
}

func (c *Cache) NewCache() *redis.Client {
	return c.cache
}
func (c *Cache) Set(key, value string) error {

	result := c.cache.Set(context.TODO(), key, value, time.Duration(0))
	return result.Err()
}

func (c *Cache) Get(key, defaultValue string) string {

	str := c.cache.Get(context.TODO(), key).String()
	if str == "" {
		str = defaultValue
	}
	return str
}
