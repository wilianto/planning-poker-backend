package schema

import "entgo.io/ent"

// Card holds the schema definition for the Card entity.
type Card struct {
	ent.Schema
}

// Fields of the Card.
func (Card) Fields() []ent.Field {
	return nil
}

// Edges of the Card.
func (Card) Edges() []ent.Edge {
	return nil
}
