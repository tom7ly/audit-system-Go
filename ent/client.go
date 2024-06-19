// Code generated by ent, DO NOT EDIT.

package ent

import (
	"context"
	"errors"
	"fmt"
	"log"
	"reflect"

	"audit-system/ent/migrate"

	"audit-system/ent/account"
	"audit-system/ent/auditlog"
	"audit-system/ent/transaction"
	"audit-system/ent/user"

	"entgo.io/ent"
	"entgo.io/ent/dialect"
	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
)

// Client is the client that holds all ent builders.
type Client struct {
	config
	// Schema is the client for creating, migrating and dropping schema.
	Schema *migrate.Schema
	// Account is the client for interacting with the Account builders.
	Account *AccountClient
	// AuditLog is the client for interacting with the AuditLog builders.
	AuditLog *AuditLogClient
	// Transaction is the client for interacting with the Transaction builders.
	Transaction *TransactionClient
	// User is the client for interacting with the User builders.
	User *UserClient
}

// NewClient creates a new client configured with the given options.
func NewClient(opts ...Option) *Client {
	client := &Client{config: newConfig(opts...)}
	client.init()
	return client
}

func (c *Client) init() {
	c.Schema = migrate.NewSchema(c.driver)
	c.Account = NewAccountClient(c.config)
	c.AuditLog = NewAuditLogClient(c.config)
	c.Transaction = NewTransactionClient(c.config)
	c.User = NewUserClient(c.config)
}

type (
	// config is the configuration for the client and its builder.
	config struct {
		// driver used for executing database requests.
		driver dialect.Driver
		// debug enable a debug logging.
		debug bool
		// log used for logging on debug mode.
		log func(...any)
		// hooks to execute on mutations.
		hooks *hooks
		// interceptors to execute on queries.
		inters *inters
	}
	// Option function to configure the client.
	Option func(*config)
)

// newConfig creates a new config for the client.
func newConfig(opts ...Option) config {
	cfg := config{log: log.Println, hooks: &hooks{}, inters: &inters{}}
	cfg.options(opts...)
	return cfg
}

// options applies the options on the config object.
func (c *config) options(opts ...Option) {
	for _, opt := range opts {
		opt(c)
	}
	if c.debug {
		c.driver = dialect.Debug(c.driver, c.log)
	}
}

// Debug enables debug logging on the ent.Driver.
func Debug() Option {
	return func(c *config) {
		c.debug = true
	}
}

// Log sets the logging function for debug mode.
func Log(fn func(...any)) Option {
	return func(c *config) {
		c.log = fn
	}
}

// Driver configures the client driver.
func Driver(driver dialect.Driver) Option {
	return func(c *config) {
		c.driver = driver
	}
}

// Open opens a database/sql.DB specified by the driver name and
// the data source name, and returns a new client attached to it.
// Optional parameters can be added for configuring the client.
func Open(driverName, dataSourceName string, options ...Option) (*Client, error) {
	switch driverName {
	case dialect.MySQL, dialect.Postgres, dialect.SQLite:
		drv, err := sql.Open(driverName, dataSourceName)
		if err != nil {
			return nil, err
		}
		return NewClient(append(options, Driver(drv))...), nil
	default:
		return nil, fmt.Errorf("unsupported driver: %q", driverName)
	}
}

// ErrTxStarted is returned when trying to start a new transaction from a transactional client.
var ErrTxStarted = errors.New("ent: cannot start a transaction within a transaction")

// Tx returns a new transactional client. The provided context
// is used until the transaction is committed or rolled back.
func (c *Client) Tx(ctx context.Context) (*Tx, error) {
	if _, ok := c.driver.(*txDriver); ok {
		return nil, ErrTxStarted
	}
	tx, err := newTx(ctx, c.driver)
	if err != nil {
		return nil, fmt.Errorf("ent: starting a transaction: %w", err)
	}
	cfg := c.config
	cfg.driver = tx
	return &Tx{
		ctx:         ctx,
		config:      cfg,
		Account:     NewAccountClient(cfg),
		AuditLog:    NewAuditLogClient(cfg),
		Transaction: NewTransactionClient(cfg),
		User:        NewUserClient(cfg),
	}, nil
}

// BeginTx returns a transactional client with specified options.
func (c *Client) BeginTx(ctx context.Context, opts *sql.TxOptions) (*Tx, error) {
	if _, ok := c.driver.(*txDriver); ok {
		return nil, errors.New("ent: cannot start a transaction within a transaction")
	}
	tx, err := c.driver.(interface {
		BeginTx(context.Context, *sql.TxOptions) (dialect.Tx, error)
	}).BeginTx(ctx, opts)
	if err != nil {
		return nil, fmt.Errorf("ent: starting a transaction: %w", err)
	}
	cfg := c.config
	cfg.driver = &txDriver{tx: tx, drv: c.driver}
	return &Tx{
		ctx:         ctx,
		config:      cfg,
		Account:     NewAccountClient(cfg),
		AuditLog:    NewAuditLogClient(cfg),
		Transaction: NewTransactionClient(cfg),
		User:        NewUserClient(cfg),
	}, nil
}

// Debug returns a new debug-client. It's used to get verbose logging on specific operations.
//
//	client.Debug().
//		Account.
//		Query().
//		Count(ctx)
func (c *Client) Debug() *Client {
	if c.debug {
		return c
	}
	cfg := c.config
	cfg.driver = dialect.Debug(c.driver, c.log)
	client := &Client{config: cfg}
	client.init()
	return client
}

// Close closes the database connection and prevents new queries from starting.
func (c *Client) Close() error {
	return c.driver.Close()
}

// Use adds the mutation hooks to all the entity clients.
// In order to add hooks to a specific client, call: `client.Node.Use(...)`.
func (c *Client) Use(hooks ...Hook) {
	c.Account.Use(hooks...)
	c.AuditLog.Use(hooks...)
	c.Transaction.Use(hooks...)
	c.User.Use(hooks...)
}

// Intercept adds the query interceptors to all the entity clients.
// In order to add interceptors to a specific client, call: `client.Node.Intercept(...)`.
func (c *Client) Intercept(interceptors ...Interceptor) {
	c.Account.Intercept(interceptors...)
	c.AuditLog.Intercept(interceptors...)
	c.Transaction.Intercept(interceptors...)
	c.User.Intercept(interceptors...)
}

// Mutate implements the ent.Mutator interface.
func (c *Client) Mutate(ctx context.Context, m Mutation) (Value, error) {
	switch m := m.(type) {
	case *AccountMutation:
		return c.Account.mutate(ctx, m)
	case *AuditLogMutation:
		return c.AuditLog.mutate(ctx, m)
	case *TransactionMutation:
		return c.Transaction.mutate(ctx, m)
	case *UserMutation:
		return c.User.mutate(ctx, m)
	default:
		return nil, fmt.Errorf("ent: unknown mutation type %T", m)
	}
}

// AccountClient is a client for the Account schema.
type AccountClient struct {
	config
}

// NewAccountClient returns a client for the Account from the given config.
func NewAccountClient(c config) *AccountClient {
	return &AccountClient{config: c}
}

// Use adds a list of mutation hooks to the hooks stack.
// A call to `Use(f, g, h)` equals to `account.Hooks(f(g(h())))`.
func (c *AccountClient) Use(hooks ...Hook) {
	c.hooks.Account = append(c.hooks.Account, hooks...)
}

// Intercept adds a list of query interceptors to the interceptors stack.
// A call to `Intercept(f, g, h)` equals to `account.Intercept(f(g(h())))`.
func (c *AccountClient) Intercept(interceptors ...Interceptor) {
	c.inters.Account = append(c.inters.Account, interceptors...)
}

// Create returns a builder for creating a Account entity.
func (c *AccountClient) Create() *AccountCreate {
	mutation := newAccountMutation(c.config, OpCreate)
	return &AccountCreate{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// CreateBulk returns a builder for creating a bulk of Account entities.
func (c *AccountClient) CreateBulk(builders ...*AccountCreate) *AccountCreateBulk {
	return &AccountCreateBulk{config: c.config, builders: builders}
}

// MapCreateBulk creates a bulk creation builder from the given slice. For each item in the slice, the function creates
// a builder and applies setFunc on it.
func (c *AccountClient) MapCreateBulk(slice any, setFunc func(*AccountCreate, int)) *AccountCreateBulk {
	rv := reflect.ValueOf(slice)
	if rv.Kind() != reflect.Slice {
		return &AccountCreateBulk{err: fmt.Errorf("calling to AccountClient.MapCreateBulk with wrong type %T, need slice", slice)}
	}
	builders := make([]*AccountCreate, rv.Len())
	for i := 0; i < rv.Len(); i++ {
		builders[i] = c.Create()
		setFunc(builders[i], i)
	}
	return &AccountCreateBulk{config: c.config, builders: builders}
}

// Update returns an update builder for Account.
func (c *AccountClient) Update() *AccountUpdate {
	mutation := newAccountMutation(c.config, OpUpdate)
	return &AccountUpdate{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// UpdateOne returns an update builder for the given entity.
func (c *AccountClient) UpdateOne(a *Account) *AccountUpdateOne {
	mutation := newAccountMutation(c.config, OpUpdateOne, withAccount(a))
	return &AccountUpdateOne{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// UpdateOneID returns an update builder for the given id.
func (c *AccountClient) UpdateOneID(id int) *AccountUpdateOne {
	mutation := newAccountMutation(c.config, OpUpdateOne, withAccountID(id))
	return &AccountUpdateOne{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// Delete returns a delete builder for Account.
func (c *AccountClient) Delete() *AccountDelete {
	mutation := newAccountMutation(c.config, OpDelete)
	return &AccountDelete{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// DeleteOne returns a builder for deleting the given entity.
func (c *AccountClient) DeleteOne(a *Account) *AccountDeleteOne {
	return c.DeleteOneID(a.ID)
}

// DeleteOneID returns a builder for deleting the given entity by its id.
func (c *AccountClient) DeleteOneID(id int) *AccountDeleteOne {
	builder := c.Delete().Where(account.ID(id))
	builder.mutation.id = &id
	builder.mutation.op = OpDeleteOne
	return &AccountDeleteOne{builder}
}

// Query returns a query builder for Account.
func (c *AccountClient) Query() *AccountQuery {
	return &AccountQuery{
		config: c.config,
		ctx:    &QueryContext{Type: TypeAccount},
		inters: c.Interceptors(),
	}
}

// Get returns a Account entity by its id.
func (c *AccountClient) Get(ctx context.Context, id int) (*Account, error) {
	return c.Query().Where(account.ID(id)).Only(ctx)
}

// GetX is like Get, but panics if an error occurs.
func (c *AccountClient) GetX(ctx context.Context, id int) *Account {
	obj, err := c.Get(ctx, id)
	if err != nil {
		panic(err)
	}
	return obj
}

// QueryUser queries the user edge of a Account.
func (c *AccountClient) QueryUser(a *Account) *UserQuery {
	query := (&UserClient{config: c.config}).Query()
	query.path = func(context.Context) (fromV *sql.Selector, _ error) {
		id := a.ID
		step := sqlgraph.NewStep(
			sqlgraph.From(account.Table, account.FieldID, id),
			sqlgraph.To(user.Table, user.FieldID),
			sqlgraph.Edge(sqlgraph.M2O, true, account.UserTable, account.UserColumn),
		)
		fromV = sqlgraph.Neighbors(a.driver.Dialect(), step)
		return fromV, nil
	}
	return query
}

// QueryOutgoingTransactions queries the outgoing_transactions edge of a Account.
func (c *AccountClient) QueryOutgoingTransactions(a *Account) *TransactionQuery {
	query := (&TransactionClient{config: c.config}).Query()
	query.path = func(context.Context) (fromV *sql.Selector, _ error) {
		id := a.ID
		step := sqlgraph.NewStep(
			sqlgraph.From(account.Table, account.FieldID, id),
			sqlgraph.To(transaction.Table, transaction.FieldID),
			sqlgraph.Edge(sqlgraph.O2M, false, account.OutgoingTransactionsTable, account.OutgoingTransactionsColumn),
		)
		fromV = sqlgraph.Neighbors(a.driver.Dialect(), step)
		return fromV, nil
	}
	return query
}

// QueryIncomingTransactions queries the incoming_transactions edge of a Account.
func (c *AccountClient) QueryIncomingTransactions(a *Account) *TransactionQuery {
	query := (&TransactionClient{config: c.config}).Query()
	query.path = func(context.Context) (fromV *sql.Selector, _ error) {
		id := a.ID
		step := sqlgraph.NewStep(
			sqlgraph.From(account.Table, account.FieldID, id),
			sqlgraph.To(transaction.Table, transaction.FieldID),
			sqlgraph.Edge(sqlgraph.O2M, false, account.IncomingTransactionsTable, account.IncomingTransactionsColumn),
		)
		fromV = sqlgraph.Neighbors(a.driver.Dialect(), step)
		return fromV, nil
	}
	return query
}

// Hooks returns the client hooks.
func (c *AccountClient) Hooks() []Hook {
	return c.hooks.Account
}

// Interceptors returns the client interceptors.
func (c *AccountClient) Interceptors() []Interceptor {
	return c.inters.Account
}

func (c *AccountClient) mutate(ctx context.Context, m *AccountMutation) (Value, error) {
	switch m.Op() {
	case OpCreate:
		return (&AccountCreate{config: c.config, hooks: c.Hooks(), mutation: m}).Save(ctx)
	case OpUpdate:
		return (&AccountUpdate{config: c.config, hooks: c.Hooks(), mutation: m}).Save(ctx)
	case OpUpdateOne:
		return (&AccountUpdateOne{config: c.config, hooks: c.Hooks(), mutation: m}).Save(ctx)
	case OpDelete, OpDeleteOne:
		return (&AccountDelete{config: c.config, hooks: c.Hooks(), mutation: m}).Exec(ctx)
	default:
		return nil, fmt.Errorf("ent: unknown Account mutation op: %q", m.Op())
	}
}

// AuditLogClient is a client for the AuditLog schema.
type AuditLogClient struct {
	config
}

// NewAuditLogClient returns a client for the AuditLog from the given config.
func NewAuditLogClient(c config) *AuditLogClient {
	return &AuditLogClient{config: c}
}

// Use adds a list of mutation hooks to the hooks stack.
// A call to `Use(f, g, h)` equals to `auditlog.Hooks(f(g(h())))`.
func (c *AuditLogClient) Use(hooks ...Hook) {
	c.hooks.AuditLog = append(c.hooks.AuditLog, hooks...)
}

// Intercept adds a list of query interceptors to the interceptors stack.
// A call to `Intercept(f, g, h)` equals to `auditlog.Intercept(f(g(h())))`.
func (c *AuditLogClient) Intercept(interceptors ...Interceptor) {
	c.inters.AuditLog = append(c.inters.AuditLog, interceptors...)
}

// Create returns a builder for creating a AuditLog entity.
func (c *AuditLogClient) Create() *AuditLogCreate {
	mutation := newAuditLogMutation(c.config, OpCreate)
	return &AuditLogCreate{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// CreateBulk returns a builder for creating a bulk of AuditLog entities.
func (c *AuditLogClient) CreateBulk(builders ...*AuditLogCreate) *AuditLogCreateBulk {
	return &AuditLogCreateBulk{config: c.config, builders: builders}
}

// MapCreateBulk creates a bulk creation builder from the given slice. For each item in the slice, the function creates
// a builder and applies setFunc on it.
func (c *AuditLogClient) MapCreateBulk(slice any, setFunc func(*AuditLogCreate, int)) *AuditLogCreateBulk {
	rv := reflect.ValueOf(slice)
	if rv.Kind() != reflect.Slice {
		return &AuditLogCreateBulk{err: fmt.Errorf("calling to AuditLogClient.MapCreateBulk with wrong type %T, need slice", slice)}
	}
	builders := make([]*AuditLogCreate, rv.Len())
	for i := 0; i < rv.Len(); i++ {
		builders[i] = c.Create()
		setFunc(builders[i], i)
	}
	return &AuditLogCreateBulk{config: c.config, builders: builders}
}

// Update returns an update builder for AuditLog.
func (c *AuditLogClient) Update() *AuditLogUpdate {
	mutation := newAuditLogMutation(c.config, OpUpdate)
	return &AuditLogUpdate{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// UpdateOne returns an update builder for the given entity.
func (c *AuditLogClient) UpdateOne(al *AuditLog) *AuditLogUpdateOne {
	mutation := newAuditLogMutation(c.config, OpUpdateOne, withAuditLog(al))
	return &AuditLogUpdateOne{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// UpdateOneID returns an update builder for the given id.
func (c *AuditLogClient) UpdateOneID(id int) *AuditLogUpdateOne {
	mutation := newAuditLogMutation(c.config, OpUpdateOne, withAuditLogID(id))
	return &AuditLogUpdateOne{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// Delete returns a delete builder for AuditLog.
func (c *AuditLogClient) Delete() *AuditLogDelete {
	mutation := newAuditLogMutation(c.config, OpDelete)
	return &AuditLogDelete{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// DeleteOne returns a builder for deleting the given entity.
func (c *AuditLogClient) DeleteOne(al *AuditLog) *AuditLogDeleteOne {
	return c.DeleteOneID(al.ID)
}

// DeleteOneID returns a builder for deleting the given entity by its id.
func (c *AuditLogClient) DeleteOneID(id int) *AuditLogDeleteOne {
	builder := c.Delete().Where(auditlog.ID(id))
	builder.mutation.id = &id
	builder.mutation.op = OpDeleteOne
	return &AuditLogDeleteOne{builder}
}

// Query returns a query builder for AuditLog.
func (c *AuditLogClient) Query() *AuditLogQuery {
	return &AuditLogQuery{
		config: c.config,
		ctx:    &QueryContext{Type: TypeAuditLog},
		inters: c.Interceptors(),
	}
}

// Get returns a AuditLog entity by its id.
func (c *AuditLogClient) Get(ctx context.Context, id int) (*AuditLog, error) {
	return c.Query().Where(auditlog.ID(id)).Only(ctx)
}

// GetX is like Get, but panics if an error occurs.
func (c *AuditLogClient) GetX(ctx context.Context, id int) *AuditLog {
	obj, err := c.Get(ctx, id)
	if err != nil {
		panic(err)
	}
	return obj
}

// Hooks returns the client hooks.
func (c *AuditLogClient) Hooks() []Hook {
	return c.hooks.AuditLog
}

// Interceptors returns the client interceptors.
func (c *AuditLogClient) Interceptors() []Interceptor {
	return c.inters.AuditLog
}

func (c *AuditLogClient) mutate(ctx context.Context, m *AuditLogMutation) (Value, error) {
	switch m.Op() {
	case OpCreate:
		return (&AuditLogCreate{config: c.config, hooks: c.Hooks(), mutation: m}).Save(ctx)
	case OpUpdate:
		return (&AuditLogUpdate{config: c.config, hooks: c.Hooks(), mutation: m}).Save(ctx)
	case OpUpdateOne:
		return (&AuditLogUpdateOne{config: c.config, hooks: c.Hooks(), mutation: m}).Save(ctx)
	case OpDelete, OpDeleteOne:
		return (&AuditLogDelete{config: c.config, hooks: c.Hooks(), mutation: m}).Exec(ctx)
	default:
		return nil, fmt.Errorf("ent: unknown AuditLog mutation op: %q", m.Op())
	}
}

// TransactionClient is a client for the Transaction schema.
type TransactionClient struct {
	config
}

// NewTransactionClient returns a client for the Transaction from the given config.
func NewTransactionClient(c config) *TransactionClient {
	return &TransactionClient{config: c}
}

// Use adds a list of mutation hooks to the hooks stack.
// A call to `Use(f, g, h)` equals to `transaction.Hooks(f(g(h())))`.
func (c *TransactionClient) Use(hooks ...Hook) {
	c.hooks.Transaction = append(c.hooks.Transaction, hooks...)
}

// Intercept adds a list of query interceptors to the interceptors stack.
// A call to `Intercept(f, g, h)` equals to `transaction.Intercept(f(g(h())))`.
func (c *TransactionClient) Intercept(interceptors ...Interceptor) {
	c.inters.Transaction = append(c.inters.Transaction, interceptors...)
}

// Create returns a builder for creating a Transaction entity.
func (c *TransactionClient) Create() *TransactionCreate {
	mutation := newTransactionMutation(c.config, OpCreate)
	return &TransactionCreate{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// CreateBulk returns a builder for creating a bulk of Transaction entities.
func (c *TransactionClient) CreateBulk(builders ...*TransactionCreate) *TransactionCreateBulk {
	return &TransactionCreateBulk{config: c.config, builders: builders}
}

// MapCreateBulk creates a bulk creation builder from the given slice. For each item in the slice, the function creates
// a builder and applies setFunc on it.
func (c *TransactionClient) MapCreateBulk(slice any, setFunc func(*TransactionCreate, int)) *TransactionCreateBulk {
	rv := reflect.ValueOf(slice)
	if rv.Kind() != reflect.Slice {
		return &TransactionCreateBulk{err: fmt.Errorf("calling to TransactionClient.MapCreateBulk with wrong type %T, need slice", slice)}
	}
	builders := make([]*TransactionCreate, rv.Len())
	for i := 0; i < rv.Len(); i++ {
		builders[i] = c.Create()
		setFunc(builders[i], i)
	}
	return &TransactionCreateBulk{config: c.config, builders: builders}
}

// Update returns an update builder for Transaction.
func (c *TransactionClient) Update() *TransactionUpdate {
	mutation := newTransactionMutation(c.config, OpUpdate)
	return &TransactionUpdate{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// UpdateOne returns an update builder for the given entity.
func (c *TransactionClient) UpdateOne(t *Transaction) *TransactionUpdateOne {
	mutation := newTransactionMutation(c.config, OpUpdateOne, withTransaction(t))
	return &TransactionUpdateOne{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// UpdateOneID returns an update builder for the given id.
func (c *TransactionClient) UpdateOneID(id int) *TransactionUpdateOne {
	mutation := newTransactionMutation(c.config, OpUpdateOne, withTransactionID(id))
	return &TransactionUpdateOne{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// Delete returns a delete builder for Transaction.
func (c *TransactionClient) Delete() *TransactionDelete {
	mutation := newTransactionMutation(c.config, OpDelete)
	return &TransactionDelete{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// DeleteOne returns a builder for deleting the given entity.
func (c *TransactionClient) DeleteOne(t *Transaction) *TransactionDeleteOne {
	return c.DeleteOneID(t.ID)
}

// DeleteOneID returns a builder for deleting the given entity by its id.
func (c *TransactionClient) DeleteOneID(id int) *TransactionDeleteOne {
	builder := c.Delete().Where(transaction.ID(id))
	builder.mutation.id = &id
	builder.mutation.op = OpDeleteOne
	return &TransactionDeleteOne{builder}
}

// Query returns a query builder for Transaction.
func (c *TransactionClient) Query() *TransactionQuery {
	return &TransactionQuery{
		config: c.config,
		ctx:    &QueryContext{Type: TypeTransaction},
		inters: c.Interceptors(),
	}
}

// Get returns a Transaction entity by its id.
func (c *TransactionClient) Get(ctx context.Context, id int) (*Transaction, error) {
	return c.Query().Where(transaction.ID(id)).Only(ctx)
}

// GetX is like Get, but panics if an error occurs.
func (c *TransactionClient) GetX(ctx context.Context, id int) *Transaction {
	obj, err := c.Get(ctx, id)
	if err != nil {
		panic(err)
	}
	return obj
}

// QueryFromAccount queries the from_account edge of a Transaction.
func (c *TransactionClient) QueryFromAccount(t *Transaction) *AccountQuery {
	query := (&AccountClient{config: c.config}).Query()
	query.path = func(context.Context) (fromV *sql.Selector, _ error) {
		id := t.ID
		step := sqlgraph.NewStep(
			sqlgraph.From(transaction.Table, transaction.FieldID, id),
			sqlgraph.To(account.Table, account.FieldID),
			sqlgraph.Edge(sqlgraph.M2O, true, transaction.FromAccountTable, transaction.FromAccountColumn),
		)
		fromV = sqlgraph.Neighbors(t.driver.Dialect(), step)
		return fromV, nil
	}
	return query
}

// QueryToAccount queries the to_account edge of a Transaction.
func (c *TransactionClient) QueryToAccount(t *Transaction) *AccountQuery {
	query := (&AccountClient{config: c.config}).Query()
	query.path = func(context.Context) (fromV *sql.Selector, _ error) {
		id := t.ID
		step := sqlgraph.NewStep(
			sqlgraph.From(transaction.Table, transaction.FieldID, id),
			sqlgraph.To(account.Table, account.FieldID),
			sqlgraph.Edge(sqlgraph.M2O, true, transaction.ToAccountTable, transaction.ToAccountColumn),
		)
		fromV = sqlgraph.Neighbors(t.driver.Dialect(), step)
		return fromV, nil
	}
	return query
}

// Hooks returns the client hooks.
func (c *TransactionClient) Hooks() []Hook {
	return c.hooks.Transaction
}

// Interceptors returns the client interceptors.
func (c *TransactionClient) Interceptors() []Interceptor {
	return c.inters.Transaction
}

func (c *TransactionClient) mutate(ctx context.Context, m *TransactionMutation) (Value, error) {
	switch m.Op() {
	case OpCreate:
		return (&TransactionCreate{config: c.config, hooks: c.Hooks(), mutation: m}).Save(ctx)
	case OpUpdate:
		return (&TransactionUpdate{config: c.config, hooks: c.Hooks(), mutation: m}).Save(ctx)
	case OpUpdateOne:
		return (&TransactionUpdateOne{config: c.config, hooks: c.Hooks(), mutation: m}).Save(ctx)
	case OpDelete, OpDeleteOne:
		return (&TransactionDelete{config: c.config, hooks: c.Hooks(), mutation: m}).Exec(ctx)
	default:
		return nil, fmt.Errorf("ent: unknown Transaction mutation op: %q", m.Op())
	}
}

// UserClient is a client for the User schema.
type UserClient struct {
	config
}

// NewUserClient returns a client for the User from the given config.
func NewUserClient(c config) *UserClient {
	return &UserClient{config: c}
}

// Use adds a list of mutation hooks to the hooks stack.
// A call to `Use(f, g, h)` equals to `user.Hooks(f(g(h())))`.
func (c *UserClient) Use(hooks ...Hook) {
	c.hooks.User = append(c.hooks.User, hooks...)
}

// Intercept adds a list of query interceptors to the interceptors stack.
// A call to `Intercept(f, g, h)` equals to `user.Intercept(f(g(h())))`.
func (c *UserClient) Intercept(interceptors ...Interceptor) {
	c.inters.User = append(c.inters.User, interceptors...)
}

// Create returns a builder for creating a User entity.
func (c *UserClient) Create() *UserCreate {
	mutation := newUserMutation(c.config, OpCreate)
	return &UserCreate{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// CreateBulk returns a builder for creating a bulk of User entities.
func (c *UserClient) CreateBulk(builders ...*UserCreate) *UserCreateBulk {
	return &UserCreateBulk{config: c.config, builders: builders}
}

// MapCreateBulk creates a bulk creation builder from the given slice. For each item in the slice, the function creates
// a builder and applies setFunc on it.
func (c *UserClient) MapCreateBulk(slice any, setFunc func(*UserCreate, int)) *UserCreateBulk {
	rv := reflect.ValueOf(slice)
	if rv.Kind() != reflect.Slice {
		return &UserCreateBulk{err: fmt.Errorf("calling to UserClient.MapCreateBulk with wrong type %T, need slice", slice)}
	}
	builders := make([]*UserCreate, rv.Len())
	for i := 0; i < rv.Len(); i++ {
		builders[i] = c.Create()
		setFunc(builders[i], i)
	}
	return &UserCreateBulk{config: c.config, builders: builders}
}

// Update returns an update builder for User.
func (c *UserClient) Update() *UserUpdate {
	mutation := newUserMutation(c.config, OpUpdate)
	return &UserUpdate{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// UpdateOne returns an update builder for the given entity.
func (c *UserClient) UpdateOne(u *User) *UserUpdateOne {
	mutation := newUserMutation(c.config, OpUpdateOne, withUser(u))
	return &UserUpdateOne{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// UpdateOneID returns an update builder for the given id.
func (c *UserClient) UpdateOneID(id int) *UserUpdateOne {
	mutation := newUserMutation(c.config, OpUpdateOne, withUserID(id))
	return &UserUpdateOne{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// Delete returns a delete builder for User.
func (c *UserClient) Delete() *UserDelete {
	mutation := newUserMutation(c.config, OpDelete)
	return &UserDelete{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// DeleteOne returns a builder for deleting the given entity.
func (c *UserClient) DeleteOne(u *User) *UserDeleteOne {
	return c.DeleteOneID(u.ID)
}

// DeleteOneID returns a builder for deleting the given entity by its id.
func (c *UserClient) DeleteOneID(id int) *UserDeleteOne {
	builder := c.Delete().Where(user.ID(id))
	builder.mutation.id = &id
	builder.mutation.op = OpDeleteOne
	return &UserDeleteOne{builder}
}

// Query returns a query builder for User.
func (c *UserClient) Query() *UserQuery {
	return &UserQuery{
		config: c.config,
		ctx:    &QueryContext{Type: TypeUser},
		inters: c.Interceptors(),
	}
}

// Get returns a User entity by its id.
func (c *UserClient) Get(ctx context.Context, id int) (*User, error) {
	return c.Query().Where(user.ID(id)).Only(ctx)
}

// GetX is like Get, but panics if an error occurs.
func (c *UserClient) GetX(ctx context.Context, id int) *User {
	obj, err := c.Get(ctx, id)
	if err != nil {
		panic(err)
	}
	return obj
}

// QueryAccounts queries the accounts edge of a User.
func (c *UserClient) QueryAccounts(u *User) *AccountQuery {
	query := (&AccountClient{config: c.config}).Query()
	query.path = func(context.Context) (fromV *sql.Selector, _ error) {
		id := u.ID
		step := sqlgraph.NewStep(
			sqlgraph.From(user.Table, user.FieldID, id),
			sqlgraph.To(account.Table, account.FieldID),
			sqlgraph.Edge(sqlgraph.O2M, false, user.AccountsTable, user.AccountsColumn),
		)
		fromV = sqlgraph.Neighbors(u.driver.Dialect(), step)
		return fromV, nil
	}
	return query
}

// Hooks returns the client hooks.
func (c *UserClient) Hooks() []Hook {
	return c.hooks.User
}

// Interceptors returns the client interceptors.
func (c *UserClient) Interceptors() []Interceptor {
	return c.inters.User
}

func (c *UserClient) mutate(ctx context.Context, m *UserMutation) (Value, error) {
	switch m.Op() {
	case OpCreate:
		return (&UserCreate{config: c.config, hooks: c.Hooks(), mutation: m}).Save(ctx)
	case OpUpdate:
		return (&UserUpdate{config: c.config, hooks: c.Hooks(), mutation: m}).Save(ctx)
	case OpUpdateOne:
		return (&UserUpdateOne{config: c.config, hooks: c.Hooks(), mutation: m}).Save(ctx)
	case OpDelete, OpDeleteOne:
		return (&UserDelete{config: c.config, hooks: c.Hooks(), mutation: m}).Exec(ctx)
	default:
		return nil, fmt.Errorf("ent: unknown User mutation op: %q", m.Op())
	}
}

// hooks and interceptors per client, for fast access.
type (
	hooks struct {
		Account, AuditLog, Transaction, User []ent.Hook
	}
	inters struct {
		Account, AuditLog, Transaction, User []ent.Interceptor
	}
)