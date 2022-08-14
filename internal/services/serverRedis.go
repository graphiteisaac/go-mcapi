package services

import (
	"context"
	"errors"
	"mc-api/internal/config"
	"mc-api/internal/db"
	"mc-api/internal/util"
)

func GetServerFromRedis(c context.Context, addr util.MinecraftAddress) (ping util.PingResponse, err error) {
	if !config.CacheMode {
		return ping, errors.New("does not exist")
	}

	res, err := db.Redis.Get(c, addr.Combined).Result()
	if err != nil {
		return ping, errors.New("does not exist")
	}

	ping, err = util.CreateResponseObj(res, addr, true)
	return
}
