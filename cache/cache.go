package cache

import (
	"context"
	"errors"
	"github.com/go-redis/redis/v8"
	"os"
	"time"
)

type Cache struct {
	Enabled bool
	Rdb     *redis.Client
}

func New() *Cache {
	if os.Getenv("CACHE") == "true" {
		return &Cache{
			Enabled: true,
			Rdb: redis.NewClient(&redis.Options{
				Addr:     os.Getenv("REDIS_HOST"),
				Password: os.Getenv("REDIS_PASS"),
				DB:       0,
			}),
		}
	}

	return &Cache{
		Enabled: false,
	}
}

func (ch *Cache) GetServer(c context.Context, key string) (string, error) {
	if !ch.Enabled {
		return "", errors.New("cache disabled")
	}

	return ch.Rdb.Get(c, key).Result()
}

func (ch *Cache) SetServer(c context.Context, key string, json string) error {
	if !ch.Enabled {
		return nil
	}

	return ch.Rdb.Set(c, key, json, time.Second*5).Err()
}
