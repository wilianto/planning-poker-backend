package room

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/wilianto/planning-poker-backend/model/schema/ent"
)

type Service struct {
	client *ent.Client
}

func NewService(client *ent.Client) *Service {
	return &Service{
		client: client,
	}
}

func (s *Service) Create(ctx context.Context, name string) (*ent.Room, error) {
	defaultConfig := map[string]interface{}{
		"options": []string{"0", "0.5", "1", "2", "3", "5", "8", "13", "21", "?"},
	}
	room, err := s.client.Room.Create().
		SetName(name).
		SetConfig(defaultConfig).
		Save(ctx)

	if err != nil {
		return nil, fmt.Errorf("failed creating room: %w", err)
	}

	return room, nil
}

func (s *Service) Join(ctx context.Context, roomID uuid.UUID, playerName string) (*ent.Player, error) {
	room, err := s.client.Room.Get(ctx, roomID)
	if err != nil {
		return nil, fmt.Errorf("failed getting room: %w", err)
	}

	player, err := s.client.Player.Create().
		SetName(playerName).
		SetRoomID(room.ID).
		Save(ctx)

	if err != nil {
		return nil, fmt.Errorf("failed creating player: %w", err)
	}

	return player, nil
}
