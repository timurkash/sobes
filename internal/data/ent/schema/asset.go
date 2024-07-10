package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/field"
	"time"
)

// Asset holds the schema definition for the Asset entity.
type Asset struct {
	ent.Schema
}

// Fields of the Asset.
func (Asset) Fields() []ent.Field {
	return []ent.Field{
		field.Uint64("id").StorageKey("uid"),
		field.String("name"),
		field.Bytes("data"),
		field.Time("created_at").SchemaType(TimestampTz).Default(time.Now),
	}
}

// Edges of the Asset.
func (Asset) Edges() []ent.Edge {
	return nil
}
