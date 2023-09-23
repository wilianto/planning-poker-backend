package room

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/google/uuid"
	"github.com/wilianto/planning-poker-backend/common"
)

type HttpTransport struct {
	roomService *Service
}

func NewHttpTransport(roomService *Service) *HttpTransport {
	return &HttpTransport{
		roomService: roomService,
	}
}

type CreateRequest struct {
	Name string `json:"name"`
}

type CreateResponse struct {
	ID        uuid.UUID              `json:"id"`
	Name      string                 `json:"name"`
	Config    map[string]interface{} `json:"config"`
	CreatedAt time.Time              `json:"created_at"`
	UpdatedAt time.Time              `json:"updated_at"`
}

// @Summary Create a new room
// @Description Create a new room with name
// @Tags room
// @Accept json
// @Produce json
// @Param body body CreateRequest true "Create Request"
// @Success 201 {object} CreateResponse
// @Failure 400 {object} common.HttpErrorReponse
// @Failure 500 {object} common.HttpErrorReponse
// @Router /api/v1/rooms [post]
func (h *HttpTransport) Create(c *fiber.Ctx) error {
	var req CreateRequest
	if err := c.BodyParser(&req); err != nil {
		log.Infof("failed parsing request body", "error", err)
		return c.Status(fiber.StatusBadRequest).JSON(common.HttpErrorReponse{Message: err.Error()})
	}

	room, err := h.roomService.Create(c.Context(), req.Name)
	if err != nil {
		log.Errorw("failed creating room", "error", err)
		return c.Status(fiber.StatusInternalServerError).JSON(common.HttpErrorReponse{Message: err.Error()})
	}

	resp := CreateResponse{
		ID:        room.ID,
		Name:      room.Name,
		Config:    room.Config,
		CreatedAt: room.CreatedAt,
		UpdatedAt: room.UpdatedAt,
	}

	return c.Status(fiber.StatusCreated).JSON(resp)
}

type JoinRequest struct {
	Name string `json:"name"`
}

type JoinResponse struct {
	PlayerID uuid.UUID `json:"player_id"`
	RoomID   uuid.UUID `json:"room_id"`
}

func (h *HttpTransport) Join(c *fiber.Ctx) error {
	id := c.Params("id")
	roomID, err := uuid.Parse(id)

	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(common.HttpErrorReponse{Message: "room not found"})
	}

	var req JoinRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(common.HttpErrorReponse{Message: err.Error()})
	}

	player, err := h.roomService.Join(c.Context(), roomID, "player")
	if err != nil {
		log.Errorw("failed joining room", "error", err)
		return c.Status(fiber.StatusInternalServerError).JSON(common.HttpErrorReponse{Message: err.Error()})
	}

	resp := JoinResponse{
		PlayerID: player.ID,
		RoomID:   player.RoomID,
	}

	return c.Status(fiber.StatusOK).JSON(resp)
}
