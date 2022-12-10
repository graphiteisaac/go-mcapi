package routes

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func redir(c *gin.Context) {
	c.Redirect(http.StatusPermanentRedirect, "/docs")
}

func RegisterStaticRoutes(g *gin.Engine) {
	g.Any("/", redir)
	g.Any("/v1", redir)
	g.Static("/docs", "./public")
}
