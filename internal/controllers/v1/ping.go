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

func PingServer(ch *cache.Cache) func(c *gin.Context) {
	return func(c *gin.Context) {
		address, err := minecraft.ParseIP(c.Param("ip"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		redisResponse, err := ch.GetServer(c, address.Combined)

		if err != nil && (err.Error() == "cache disabled" || errors.Is(err, redis.Nil)) {
			tcpServer, terr := minecraft.PingServer(address)
			if terr != nil {
				fmt.Println(terr)
				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": terr.Error()})
				return
			}

			if err := ch.SetServer(c, address.Combined, tcpServer); err != nil {
				fmt.Println("failed to cache")
			}

			server, err := minecraft.CreateResponseObj(tcpServer, address.IPv4, address.Port, false)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			c.JSON(http.StatusOK, server)
			return
		} else if err != nil {
			fmt.Println("failed to retrieve cache")
			fmt.Println(err)
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		server, err := minecraft.CreateResponseObj(redisResponse, address.IPv4, address.Port, true)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, server)
	}
}
