package room

import (
	"github.com/gofiber/fiber/v2"
	"github.com/wilianto/planning-poker-backend/model/schema/ent"
)

func InitHttpEndpoints(app fiber.Router, client *ent.Client) {
	service := NewService(client)
	roomHttp := NewHttpTransport(service)

	room := app.Group("/rooms")
	room.Post("/", roomHttp.Create)
}
