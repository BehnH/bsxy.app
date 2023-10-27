package main

import (
	"bsxy_app/database"
	"bsxy_app/router"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/redis/go-redis/v9"
	"log"
	"os"
)

func main() {
	app := fiber.New()
	app.Use(cors.New())

	// Create the Redis client
	redisUrl := os.Getenv("REDIS_URL")
	opt, _ := redis.ParseURL(redisUrl)
	database.RedisClient = redis.NewClient(opt)

	router.Router(app)
	log.Fatal(app.Listen(":3000"))
}
