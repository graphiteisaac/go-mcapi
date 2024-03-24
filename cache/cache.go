package cache

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/graphiteisaac/go-mcapi/minecraft"
)

type Cache struct {
	Enabled bool
	Expiry  int
	Rdb     *redis.Client
}

func New(enabled bool, url string, expiry time.Duration) (*Cache, error) {
	if !enabled {
		return &Cache{
			Enabled: false,
		}, nil
	}

	opts, err := redis.ParseURL(url)
	if err != nil {
		return nil, err
	}

	client := redis.NewClient(opts)
	c, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	if client.Ping(c).String() != "PONG" {
		return nil, errors.New("failed to ping redis server")
	}

	return &Cache{
		Enabled: true,
		Expiry:  int(expiry.Seconds()),
		Rdb:     client,
	}, nil
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

	return ch.Rdb.Set(c, fmt.Sprintf("%s:%d", address.Host, address.Port), json, time.Minute*time.Duration(ch.Expiry)).Err()
}
