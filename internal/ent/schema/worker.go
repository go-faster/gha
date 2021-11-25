package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
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
	return []ent.Edge{
		edge.To("chunks", Chunk.Type),
	}
}

func (Worker) Mixin() []ent.Mixin {
	return []ent.Mixin{
		IDMixin{},
		TimeMixin{},
	}
}
