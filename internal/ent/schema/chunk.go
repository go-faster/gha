package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/field"
)

// Chunk holds the schema definition for the Chunk entity.
type Chunk struct {
	ent.Schema
}

// Fields of the Chunk.
func (Chunk) Fields() []ent.Field {
	return []ent.Field{
		field.String("id").
			Unique().
			Comment("Like 2006-01-02T15"),
		field.Time("start").
			Comment("Minimum possible time of entry in chunk"),
	}
}

// Edges of the Chunk.
func (Chunk) Edges() []ent.Edge {
	return nil
}

func (Chunk) Mixin() []ent.Mixin {
	return []ent.Mixin{
		TimeMixin{},
	}
}
