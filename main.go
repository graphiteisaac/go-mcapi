package main

import (
	"flag"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"mc-api/cache"
	"mc-api/routes"
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
	// TODO source version from README
	routes.RegisterV1Routes(router, ch, "1.1.1")

	router.Run(fmt.Sprintf(":%d", port))
}
