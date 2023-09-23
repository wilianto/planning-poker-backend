package http

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/swagger"
	"github.com/wilianto/planning-poker-backend/room"
)

func Routing(
	app *fiber.App,
	roomService *room.Service,
) {
	app.Use(logger.New(logger.Config{
		Format: "${time} ${method} ${path} - ${ip} - ${status} - ${latency}\n",
	}))
	app.Use(recover.New())
	app.Get("/", func(c *fiber.Ctx) error {
		return c.JSON("Hello, World!")
	})
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON("OK")
	})
	app.Get("/swagger/*", swagger.HandlerDefault)

	roomHttp := room.NewHttpTransport(roomService)

	apiV1 := app.Group("/api/v1")
	apiV1.Route("/rooms", func(r fiber.Router) {
		r.Post("/", roomHttp.Create)
	})
}
