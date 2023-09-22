package schema

import "entgo.io/ent"

// People holds the schema definition for the People entity.
type People struct {
	ent.Schema
}

// Fields of the People.
func (People) Fields() []ent.Field {
	return nil
}

// Edges of the People.
func (People) Edges() []ent.Edge {
	return nil
}
