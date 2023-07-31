package main

import (
	"flag"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"mc-api/cache"
	"mc-api/routes"
	"net/http"
	"os"
	"regexp"
)

var port uint

func main() {
	port = *flag.Uint("port", 3333, "the port to run the server on")
	flag.Parse()

	// load in .env variables with the godotenv loader
	godotenv.Load()

	// create a new Gin router
	router := gin.Default()
	ch := cache.New()

	// Register routes
	version := "1.0.0"
	f, err := os.ReadFile("README.md")
	if err != nil {
		fmt.Println("could not open readme, defaulting to version 1.0.0")
	}

	version = regexp.MustCompile("(Version `(\\d\\.\\d\\.\\d))`").FindStringSubmatch(string(f))[2]

	routes.ServerAPIV1(router, ch, version)
	routes.RegisterDocs(router, version)
	// fallback redirect
	router.NoRoute(func(c *gin.Context) {
		c.Redirect(http.StatusPermanentRedirect, "/docs")
	})

	router.Run(fmt.Sprintf(":%d", port))
}
