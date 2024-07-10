package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/field"
	"time"
)

// User holds the schema definition for the User entity.
type User struct {
	ent.Schema
}

// Fields of the User.
func (User) Fields() []ent.Field {
	return []ent.Field{
		field.Uint64("id"),
		field.Time("created_at").SchemaType(TimestampTz).Default(time.Now),
		field.String("login").Unique(),
		field.String("password_hash").MaxLen(32),
	}
}

// Edges of the User.
func (User) Edges() []ent.Edge {
	return nil
}
