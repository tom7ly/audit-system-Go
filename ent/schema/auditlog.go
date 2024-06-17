package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/field"
)

// AuditLog holds the schema definition for the AuditLog entity.
type AuditLog struct {
	ent.Schema
}

// Fields of the AuditLog.
func (AuditLog) Fields() []ent.Field {
	return []ent.Field{
		field.String("client_id"),
		field.Time("timestamp").Default(time.Now),
		field.String("entity"),
		field.String("mutation"),
		field.JSON("before", map[string]interface{}{}).Optional(),
		field.JSON("after", map[string]interface{}{}).Optional(),
	}
}
