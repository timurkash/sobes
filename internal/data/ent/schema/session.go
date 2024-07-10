package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/field"
	"time"
)

// Session holds the schema definition for the Session entity.
type Session struct {
	ent.Schema
}

// Fields of the Session.
func (Session) Fields() []ent.Field {
	return []ent.Field{
		field.String("id"),
		field.Time("created_at").SchemaType(TimestampTz).Default(time.Now),
		field.Uint64("uid").Unique(),
	}
}

// Edges of the Session.
func (Session) Edges() []ent.Edge {
	return nil
}
