// Code generated by ent, DO NOT EDIT.

package migrate

import (
	"entgo.io/ent/dialect/sql/schema"
	"entgo.io/ent/schema/field"
)

var (
	// AccountsColumns holds the columns for the "accounts" table.
	AccountsColumns = []*schema.Column{
		{Name: "id", Type: field.TypeInt, Increment: true},
		{Name: "balance", Type: field.TypeFloat64},
		{Name: "last_transfer_time", Type: field.TypeTime},
		{Name: "user_accounts", Type: field.TypeInt},
	}
	// AccountsTable holds the schema information for the "accounts" table.
	AccountsTable = &schema.Table{
		Name:       "accounts",
		Columns:    AccountsColumns,
		PrimaryKey: []*schema.Column{AccountsColumns[0]},
		ForeignKeys: []*schema.ForeignKey{
			{
				Symbol:     "accounts_users_accounts",
				Columns:    []*schema.Column{AccountsColumns[3]},
				RefColumns: []*schema.Column{UsersColumns[0]},
				OnDelete:   schema.NoAction,
			},
		},
	}
	// AuditLogsColumns holds the columns for the "audit_logs" table.
	AuditLogsColumns = []*schema.Column{
		{Name: "id", Type: field.TypeInt, Increment: true},
		{Name: "client_id", Type: field.TypeString},
		{Name: "timestamp", Type: field.TypeTime},
		{Name: "entity", Type: field.TypeString},
		{Name: "mutation", Type: field.TypeString},
		{Name: "before", Type: field.TypeJSON, Nullable: true},
		{Name: "after", Type: field.TypeJSON, Nullable: true},
	}
	// AuditLogsTable holds the schema information for the "audit_logs" table.
	AuditLogsTable = &schema.Table{
		Name:       "audit_logs",
		Columns:    AuditLogsColumns,
		PrimaryKey: []*schema.Column{AuditLogsColumns[0]},
	}
	// TransactionsColumns holds the columns for the "transactions" table.
	TransactionsColumns = []*schema.Column{
		{Name: "id", Type: field.TypeInt, Increment: true},
		{Name: "amount", Type: field.TypeFloat64},
		{Name: "timestamp", Type: field.TypeTime},
		{Name: "from_account_id", Type: field.TypeInt},
		{Name: "to_account_id", Type: field.TypeInt},
	}
	// TransactionsTable holds the schema information for the "transactions" table.
	TransactionsTable = &schema.Table{
		Name:       "transactions",
		Columns:    TransactionsColumns,
		PrimaryKey: []*schema.Column{TransactionsColumns[0]},
		ForeignKeys: []*schema.ForeignKey{
			{
				Symbol:     "transactions_accounts_outgoing_transactions",
				Columns:    []*schema.Column{TransactionsColumns[3]},
				RefColumns: []*schema.Column{AccountsColumns[0]},
				OnDelete:   schema.NoAction,
			},
			{
				Symbol:     "transactions_accounts_incoming_transactions",
				Columns:    []*schema.Column{TransactionsColumns[4]},
				RefColumns: []*schema.Column{AccountsColumns[0]},
				OnDelete:   schema.NoAction,
			},
		},
	}
	// UsersColumns holds the columns for the "users" table.
	UsersColumns = []*schema.Column{
		{Name: "id", Type: field.TypeInt, Increment: true},
		{Name: "email", Type: field.TypeString, Unique: true},
		{Name: "name", Type: field.TypeString},
		{Name: "age", Type: field.TypeInt},
	}
	// UsersTable holds the schema information for the "users" table.
	UsersTable = &schema.Table{
		Name:       "users",
		Columns:    UsersColumns,
		PrimaryKey: []*schema.Column{UsersColumns[0]},
	}
	// Tables holds all the tables in the schema.
	Tables = []*schema.Table{
		AccountsTable,
		AuditLogsTable,
		TransactionsTable,
		UsersTable,
	}
)

func init() {
	AccountsTable.ForeignKeys[0].RefTable = UsersTable
	TransactionsTable.ForeignKeys[0].RefTable = AccountsTable
	TransactionsTable.ForeignKeys[1].RefTable = AccountsTable
}
