package routes

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/graphiteisaac/go-mcapi/cache"
	"github.com/graphiteisaac/go-mcapi/minecraft"
)

func ServerAPIV1(g *gin.Engine, ch *cache.Cache, version string) {
	g.GET("/v1/ping/:ip", func(c *gin.Context) {
		var query struct {
			Players  bool `form:"players"`
			Protocol bool `form:"protocol"`
			MOTD     bool `form:"motd"`
			Icon     bool `form:"icon"`
		}
		if err := c.BindQuery(&query); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		address, err := minecraft.ParseIP(c.Param("ip"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// By default, we will assume the result has been cached and
		// only swap it when we error out on the Redis response.
		c.Header("X-Cached", "true")
		c.Header("X-Version", version)

		// Initialize our response object
		server := &minecraft.PingResponse{}

		redisResponse, err := ch.GetServer(c, address)
		if err == nil {
			server, err = minecraft.MarshalJson(redisResponse)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
		} else if err != nil && (err.Error() == "cache disabled" || errors.Is(err, redis.Nil)) {
			c.Header("X-Cached", "false")

			tcpServer, err := minecraft.PingServer(address)
			if err != nil {
				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}

			if err := ch.SetServer(c, address, tcpServer); err != nil {
				fmt.Printf("failed to cache at %s, moving on...\n", time.Now().Format(time.RFC822))
			}

			server, err = minecraft.MarshalJson(tcpServer)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
		} else if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		if query.Players {
			c.JSON(http.StatusOK, server.Players)
			return
		} else if query.MOTD {
			c.JSON(http.StatusOK, server.Description)
			return
		} else if query.Protocol {
			c.JSON(http.StatusOK, server.Version)
			return
		} else if query.Icon {
			c.JSON(http.StatusOK, server)
			return
		}

		// If we aren't specifically requesting the icon we should default to not including it.
		c.JSON(http.StatusOK, gin.H{
			"description": server.Description,
			"players":     server.Players,
			"version":     server.Version,
		})
	})

	g.GET("/v1/icon/:ip", func(c *gin.Context) {
		address, err := minecraft.ParseIP(c.Param("ip"))
		if err != nil {
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		cached, err := ch.GetServer(c, address)
		if err != nil && (err.Error() == "cache disabled" || errors.Is(err, redis.Nil)) {
			server, err := minecraft.PingServer(address)
			if err != nil {
				c.AbortWithStatus(http.StatusInternalServerError)
				return
			}

			icon, err := minecraft.GetIcon(server)
			if err != nil {
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
	})
}
