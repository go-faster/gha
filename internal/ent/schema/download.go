package schema

import "entgo.io/ent"

// Download holds the schema definition for the Download entity.
type Download struct {
	ent.Schema
}

// Fields of the Download.
func (Download) Fields() []ent.Field {
	return nil
}

// Edges of the Download.
func (Download) Edges() []ent.Edge {
	return nil
}

func (Download) Mixin() []ent.Mixin {
	return []ent.Mixin{
		IDMixin{},
		TimeMixin{},
	}
}
