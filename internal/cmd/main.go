package main

import (
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"mc-api/internal/config"
	"mc-api/internal/db"
	"mc-api/internal/routes/v1"
)

func init() {
	godotenv.Load()
	config.LoadConfigVars()
}

func main() {
	// Create Redis client and Gin router
	r := gin.Default()
	if config.CacheMode {
		db.CreateRedisClient()
	}

	// Register routes
	routes.RegisterV1Routes(r)

	r.Run(":3333")
}
