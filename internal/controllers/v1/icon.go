package controllers

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"mc-api/internal/cache"
	"mc-api/internal/minecraft"
	"net/http"
)

func GetIcon(ch *cache.Cache) func(c *gin.Context) {
	return func(c *gin.Context) {
		address, err := minecraft.ParseIP(c.Param("ip"))
		if err != nil {
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		cached, err := ch.GetServer(c, address.Combined)
		if err != nil && (err.Error() == "cache disabled" || errors.Is(err, redis.Nil)) {
			server, err := minecraft.PingServer(address)
			if err != nil {
				c.AbortWithStatus(http.StatusInternalServerError)
				return
			}

			icon, err := minecraft.GetIcon(server)
			if err != nil {
				fmt.Println(err)
				c.AbortWithStatus(http.StatusInternalServerError)
				return
			}

			c.Data(http.StatusOK, "image/png", icon)
			return
		} else if err != nil {
			fmt.Println(err)

			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		icon, err := minecraft.GetIcon(cached)
		if err != nil {
			fmt.Println(err)
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		c.Data(http.StatusOK, "image/png", icon)
	}
}
