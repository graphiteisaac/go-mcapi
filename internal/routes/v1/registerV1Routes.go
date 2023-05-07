package routes

import (
	"github.com/gin-gonic/gin"
	"mc-api/internal/cache"
	"mc-api/internal/controllers/v1"
	"net/http"
)

func RegisterV1Routes(g *gin.Engine, c *cache.Cache) {
	// Static docs routes
	g.Static("/docs", "./public")
	g.GET("/v1/ping/:ip", controllers.PingServer(c))
	g.GET("/v1/icon/:ip", controllers.GetIcon(c))

	// fallback redirect
	g.NoRoute(func(c *gin.Context) {
		c.Redirect(http.StatusPermanentRedirect, "/docs")
	})
}
