package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"github.com/google/uuid"
)

// Card holds the schema definition for the Card entity.
type Card struct {
	ent.Schema
}

// Fields of the Card.
func (Card) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.UUID{}).Default(uuid.New).Unique().Immutable(),
		field.Time("created_at").Default(time.Now),
		field.Time("updated_at").Default(time.Now),
		field.String("value").MaxLen(10),

		field.UUID("player_id", uuid.UUID{}),
		field.UUID("game_id", uuid.UUID{}),
	}
}

// Edges of the Card.
func (Card) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("player", Player.Type).Ref("cards").Field("player_id").Required().Unique(),
		edge.From("game", Game.Type).Ref("cards").Field("game_id").Required().Unique(),
	}
}
