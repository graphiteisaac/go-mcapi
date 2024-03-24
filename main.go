package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"regexp"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/graphiteisaac/go-mcapi/cache"
	"github.com/graphiteisaac/go-mcapi/routes"
)

func main() {
	port := *flag.Uint("port", 3333, "the port to run the server on")
	enableCache := *flag.Bool("cache", false, "enable caching")
	cacheExpiry := *flag.Duration("expiry", time.Second*3, "cache expiry time")
	redisURI := *flag.String("redis-uri", "", "the full redis URI (eg redis://:pass@localhost:6379/0)")
	flag.Parse()

	// create a new Gin router
	router := gin.Default()
	ch, err := cache.New(enableCache, redisURI, cacheExpiry)
	if err != nil {
		log.Fatal(err)
	}

	// Register routes
	version := "1.1.0"
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
