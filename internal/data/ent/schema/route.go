package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/field"
)

// Route holds the schema definition for the Route entity.
type Route struct {
	ent.Schema
}

// Fields of the Route.
func (Route) Fields() []ent.Field {
	return []ent.Field{
		field.Uint64("id").StorageKey("route_id"),
		field.String("route_name"),
		field.Float("load"),
		field.String("cargo_type"),
		field.Bool("is_actual").Default(false),
	}
}

// Edges of the Route.
func (Route) Edges() []ent.Edge {
	return nil
}
