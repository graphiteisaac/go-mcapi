package cache

import (
	"context"
	"errors"
	"fmt"
	"github.com/go-redis/redis/v8"
	"mc-api/minecraft"
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

func (ch *Cache) GetServer(c context.Context, address *minecraft.Address) (string, error) {
	if !ch.Enabled {
		return "", errors.New("cache disabled")
	}

	return ch.Rdb.Get(c, fmt.Sprintf("%s:%d", address.Host, address.Port)).Result()
}

func (ch *Cache) SetServer(c context.Context, address *minecraft.Address, json string) error {
	if !ch.Enabled {
		return nil
	}

	return ch.Rdb.Set(c, fmt.Sprintf("%s:%d", address.Host, address.Port), json, time.Minute*5).Err()
}
