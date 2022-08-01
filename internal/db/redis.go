package db

import (
	"github.com/go-redis/redis/v8"
	"os"
)

var (
	Redis *redis.Client
)

func CreateRedisClient() {
	Redis = redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_HOST"),
		Password: os.Getenv("REDIS_PASS"),
		DB:       0,
	})
}
