package controllers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"mc-api/internal/services"
	"mc-api/internal/util"
	"net/http"
)

func GetIcon(c *gin.Context) {
	address, err := util.ParseIP(c.Param("ip"))
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	server, err := services.GetServerFromRedis(c, address)
	if err != nil && err.Error() == "does not exist" {
		server, err = services.PingServerTCP(c, address)
		if err != nil {
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		fmt.Println(server.Icon)

		img, err := util.Base64StringToImage(server.Icon)
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

	img, err := util.Base64StringToImage(server.Icon)
	if err != nil {
		fmt.Println(err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.Data(http.StatusOK, "image/png", img)
}
