package main

import (
	"context"
	"fmt"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/swagger"
	"github.com/joho/godotenv"

	"github.com/wilianto/planning-poker-backend/model/schema/ent"
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

	client, err := initDB()
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

	api := app.Group("/api/v1")
	initRoomEndpoints(api, client)

	appPort := os.Getenv("APP_PORT")
	app.Listen(fmt.Sprintf(":%s", appPort))
}

func initDB() (*ent.Client, error) {
	dataSourceName := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		os.Getenv("POSTGRES_HOST"),
		os.Getenv("POSTGRES_PORT"),
		os.Getenv("POSTGRES_USER"),
		os.Getenv("POSTGRES_PASSWORD"),
		os.Getenv("POSTGRES_DB"),
	)
	client, err := ent.Open("postgres", dataSourceName)
	if err != nil {
		return nil, fmt.Errorf("failed opening connection to postgres: %w", err)
	}

	if err = client.Schema.Create(context.Background()); err != nil {
		return nil, fmt.Errorf("failed creating schema resources: %w", err)
	}

	return client, nil
}

func initRoomEndpoints(app fiber.Router, client *ent.Client) {
	service := room.NewService(client)
	roomHttp := room.NewHttpTransport(service)
	room := app.Group("/rooms")
	room.Post("/", roomHttp.Create)
}
