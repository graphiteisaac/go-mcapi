package controllers

import (
	"encoding/base64"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"mc-api/internal/cache"
	"mc-api/internal/minecraft"
	"net/http"
	"strings"
)

func b64toimg(input string) ([]byte, error) {
	image := input[strings.IndexByte(input, ',')+1:]

	return base64.StdEncoding.DecodeString(image)
}

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

			img, err := b64toimg(server.Icon)
			if err != nil {
				fmt.Println(err)
				c.AbortWithStatus(http.StatusInternalServerError)
				return
			}

			c.Data(http.StatusOK, "image/png", img)
			return
		} else if err != nil {
			fmt.Println(err)

			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		imgBytes, err := b64toimg(cached)
		if err != nil {
			fmt.Println(err)
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		c.Data(http.StatusOK, "image/png", imgBytes)
	}
}
