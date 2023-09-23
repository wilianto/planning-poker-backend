package main

import (
	"fmt"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/swagger"
	"github.com/joho/godotenv"

	"github.com/wilianto/planning-poker-backend/common"
	"github.com/wilianto/planning-poker-backend/room"

	_ "github.com/lib/pq"
	_ "github.com/wilianto/planning-poker-backend/docs"
)

// @title Planning Poker API
// @version v1
// @description This is a planning poker API server.

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html
func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	client, err := common.InitDB()
	if err != nil {
		log.Fatalf("failed initializing ent: %v", err)
	}
	defer client.Close()

	app := fiber.New()
	app.Use(logger.New(logger.Config{
		Format: "${time} ${method} ${path} - ${ip} - ${status} - ${latency}\n",
	}))
	app.Get("/", func(c *fiber.Ctx) error {
		return c.JSON("Hello, World!")
	})
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON("OK")
	})
	app.Get("/swagger/*", swagger.HandlerDefault)

	apiV1 := app.Group("/api/v1")
	room.InitHttpEndpoints(apiV1, client)

	appPort := os.Getenv("APP_PORT")
	app.Listen(fmt.Sprintf(":%s", appPort))
}
