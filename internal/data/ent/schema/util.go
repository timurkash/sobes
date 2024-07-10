package schema

import "entgo.io/ent/dialect"

var TimestampTz = map[string]string{
	dialect.Postgres: "timestamptz",
}
