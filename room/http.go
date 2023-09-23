package room

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
)

type HttpTransport struct {
	roomService *Service
}

func NewHttpTransport(roomService *Service) *HttpTransport {
	return &HttpTransport{
		roomService: roomService,
	}
}

func (h *HttpTransport) Create(c *fiber.Ctx) error {
	var req struct {
		Name string `json:"name"`
	}

	if err := c.BodyParser(&req); err != nil {
		log.Infof("failed parsing request body", "error", err)
		return c.Status(fiber.StatusBadRequest).JSON(err)
	}

	room, err := h.roomService.Create(c.Context(), req.Name)
	if err != nil {
		log.Errorw("failed creating room", "error", err)
		return c.Status(fiber.StatusInternalServerError).JSON(err)
	}

	return c.Status(fiber.StatusCreated).JSON(room)
}
