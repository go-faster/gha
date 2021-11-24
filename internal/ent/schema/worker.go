package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/field"
)

// Worker holds the schema definition for the Worker entity.
type Worker struct {
	ent.Schema
}

// Fields of the Worker.
func (Worker) Fields() []ent.Field {
	return []ent.Field{
		field.String("name").Unique(),
		field.String("token"),
	}
}

// Edges of the Worker.
func (Worker) Edges() []ent.Edge {
	return nil
}

func (Worker) Mixin() []ent.Mixin {
	return []ent.Mixin{
		IDMixin{},
		TimeMixin{},
	}
}
