package routes

import (
	"github.com/gin-gonic/gin"
	"mc-api/internal/controllers/v1"
)

func RegisterV1Routes(g *gin.Engine) {
	group := g.Group("/v1")
	group.GET("/ping/:ip", controllers.PingServer)
	group.GET("/icon/:ip", controllers.GetIcon)
}
