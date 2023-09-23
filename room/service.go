package room

import (
	"context"
	"fmt"

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
		"options": []string{"0", "1/2", "1", "2", "3", "5", "8", "13", "21", "?"},
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
