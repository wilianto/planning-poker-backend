package main

import (
	"context"
	"fmt"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/joho/godotenv"
	"github.com/wilianto/planning-poker-backend/model/schema/ent"

	_ "github.com/lib/pq"
)

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
	app.Use(logger.New(logger.ConfigDefault))
	app.Get("/", func(c *fiber.Ctx) error {
		return c.JSON("Hello, World!")
	})
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON("OK")
	})

	api := app.Group("/api/v1")
	room := api.Group("/room")

	room.Post("/", func(c *fiber.Ctx) error {
		var req struct {
			Name string `json:"name"`
		}

		if err := c.BodyParser(&req); err != nil {
			log.Infof("failed parsing request body", "error", err)
			return c.Status(fiber.StatusBadRequest).JSON(err)
		}

		room, err := client.Room.Create().
			SetName(req.Name).
			SetConfig(map[string]interface{}{}).
			Save(c.Context())

		if err != nil {
			log.Errorw("failed creating room", "error", err)
			return c.Status(fiber.StatusInternalServerError).JSON(err)
		}
		return c.Status(fiber.StatusCreated).JSON(room)
	})

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
