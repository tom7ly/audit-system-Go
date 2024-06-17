package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// Account holds the schema definition for the Account entity.
type Account struct {
	ent.Schema
}

// Fields of the Account.
func (Account) Fields() []ent.Field {
	return []ent.Field{
		field.Float("balance"),
		field.Time("last_transfer_time").Default(time.Now),
	}
}

// Edges of the Account.
func (Account) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("user", User.Type).
			Ref("accounts").
			Unique().
			Required(),
		edge.To("outgoing_transactions", Transaction.Type).
			StorageKey(edge.Column("from_account_id")),
		edge.To("incoming_transactions", Transaction.Type).
			StorageKey(edge.Column("to_account_id")),
	}
}
