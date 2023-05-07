package controllers

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
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

		if err != nil && err.Error() == "cache disabled" {
			tcpServer, terr := minecraft.PingServer(address)
			if terr != nil {
				fmt.Println(terr)
				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": terr.Error()})
				return
			}

			bytes, _ := json.Marshal(tcpServer)
			if err := ch.SetServer(c, address.Combined, string(bytes)); err != nil {
				fmt.Println("failed to cache")
			}

			c.JSON(http.StatusOK, tcpServer)
			return
		} else if err != nil {
			fmt.Println("failed to retrieve cache")
			fmt.Println(err)
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, redisResponse)
	}
}
