package controllers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"mc-api/internal/services"
	"mc-api/internal/util"
	"net/http"
)

func PingServer(c *gin.Context) {
	address, err := util.ParseIP(c.Param("ip"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	redisResponse, err := services.GetServerFromRedis(c, address)

	if err != nil && err.Error() == "does not exist" {
		tcpServer, err := services.PingServerTCP(c, address)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, tcpServer)
		return
	} else if err != nil {
		fmt.Println(err)

		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, redisResponse)
}
