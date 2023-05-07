package main

import (
	"flag"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"mc-api/internal/cache"
	"mc-api/internal/routes/v1"
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
	routes.RegisterV1Routes(router, ch)

	router.Run(fmt.Sprintf(":%d", port))
}
