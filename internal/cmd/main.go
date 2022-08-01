package main

import (
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/joho/godotenv"
	"mc-api/internal/db"
	"mc-api/internal/routes/v1"
)

var (
	Rdb *redis.Client
)

func init() {
	godotenv.Load()
}

func main() {
	// Create Redis client and Gin router
	r := gin.Default()
	db.CreateRedisClient()

	// Register routes
	routes.RegisterV1Routes(r)

	r.Run(":3333")
}
