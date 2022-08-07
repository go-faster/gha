// Code generated by ent, DO NOT EDIT.

package ent

import (
	"context"
	"fmt"
	"log"

	"github.com/go-faster/gha/internal/ent/migrate"
	"github.com/google/uuid"

	"github.com/go-faster/gha/internal/ent/chunk"
	"github.com/go-faster/gha/internal/ent/worker"

	"entgo.io/ent/dialect"
	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
)

// Client is the client that holds all ent builders.
type Client struct {
	config
	// Schema is the client for creating, migrating and dropping schema.
	Schema *migrate.Schema
	// Chunk is the client for interacting with the Chunk builders.
	Chunk *ChunkClient
	// Worker is the client for interacting with the Worker builders.
	Worker *WorkerClient
}

// NewClient creates a new client configured with the given options.
func NewClient(opts ...Option) *Client {
	cfg := config{log: log.Println, hooks: &hooks{}}
	cfg.options(opts...)
	client := &Client{config: cfg}
	client.init()
	return client
}

func (c *Client) init() {
	c.Schema = migrate.NewSchema(c.driver)
	c.Chunk = NewChunkClient(c.config)
	c.Worker = NewWorkerClient(c.config)
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

// Tx returns a new transactional client. The provided context
// is used until the transaction is committed or rolled back.
func (c *Client) Tx(ctx context.Context) (*Tx, error) {
	if _, ok := c.driver.(*txDriver); ok {
		return nil, fmt.Errorf("ent: cannot start a transaction within a transaction")
	}
	tx, err := newTx(ctx, c.driver)
	if err != nil {
		return nil, fmt.Errorf("ent: starting a transaction: %w", err)
	}
	cfg := c.config
	cfg.driver = tx
	return &Tx{
		ctx:    ctx,
		config: cfg,
		Chunk:  NewChunkClient(cfg),
		Worker: NewWorkerClient(cfg),
	}, nil
}

// BeginTx returns a transactional client with specified options.
func (c *Client) BeginTx(ctx context.Context, opts *sql.TxOptions) (*Tx, error) {
	if _, ok := c.driver.(*txDriver); ok {
		return nil, fmt.Errorf("ent: cannot start a transaction within a transaction")
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
		ctx:    ctx,
		config: cfg,
		Chunk:  NewChunkClient(cfg),
		Worker: NewWorkerClient(cfg),
	}, nil
}

// Debug returns a new debug-client. It's used to get verbose logging on specific operations.
//
//	client.Debug().
//		Chunk.
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
	c.Chunk.Use(hooks...)
	c.Worker.Use(hooks...)
}

// ChunkClient is a client for the Chunk schema.
type ChunkClient struct {
	config
}

// NewChunkClient returns a client for the Chunk from the given config.
func NewChunkClient(c config) *ChunkClient {
	return &ChunkClient{config: c}
}

// Use adds a list of mutation hooks to the hooks stack.
// A call to `Use(f, g, h)` equals to `chunk.Hooks(f(g(h())))`.
func (c *ChunkClient) Use(hooks ...Hook) {
	c.hooks.Chunk = append(c.hooks.Chunk, hooks...)
}

// Create returns a builder for creating a Chunk entity.
func (c *ChunkClient) Create() *ChunkCreate {
	mutation := newChunkMutation(c.config, OpCreate)
	return &ChunkCreate{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// CreateBulk returns a builder for creating a bulk of Chunk entities.
func (c *ChunkClient) CreateBulk(builders ...*ChunkCreate) *ChunkCreateBulk {
	return &ChunkCreateBulk{config: c.config, builders: builders}
}

// Update returns an update builder for Chunk.
func (c *ChunkClient) Update() *ChunkUpdate {
	mutation := newChunkMutation(c.config, OpUpdate)
	return &ChunkUpdate{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// UpdateOne returns an update builder for the given entity.
func (c *ChunkClient) UpdateOne(ch *Chunk) *ChunkUpdateOne {
	mutation := newChunkMutation(c.config, OpUpdateOne, withChunk(ch))
	return &ChunkUpdateOne{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// UpdateOneID returns an update builder for the given id.
func (c *ChunkClient) UpdateOneID(id string) *ChunkUpdateOne {
	mutation := newChunkMutation(c.config, OpUpdateOne, withChunkID(id))
	return &ChunkUpdateOne{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// Delete returns a delete builder for Chunk.
func (c *ChunkClient) Delete() *ChunkDelete {
	mutation := newChunkMutation(c.config, OpDelete)
	return &ChunkDelete{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// DeleteOne returns a builder for deleting the given entity.
func (c *ChunkClient) DeleteOne(ch *Chunk) *ChunkDeleteOne {
	return c.DeleteOneID(ch.ID)
}

// DeleteOne returns a builder for deleting the given entity by its id.
func (c *ChunkClient) DeleteOneID(id string) *ChunkDeleteOne {
	builder := c.Delete().Where(chunk.ID(id))
	builder.mutation.id = &id
	builder.mutation.op = OpDeleteOne
	return &ChunkDeleteOne{builder}
}

// Query returns a query builder for Chunk.
func (c *ChunkClient) Query() *ChunkQuery {
	return &ChunkQuery{
		config: c.config,
	}
}

// Get returns a Chunk entity by its id.
func (c *ChunkClient) Get(ctx context.Context, id string) (*Chunk, error) {
	return c.Query().Where(chunk.ID(id)).Only(ctx)
}

// GetX is like Get, but panics if an error occurs.
func (c *ChunkClient) GetX(ctx context.Context, id string) *Chunk {
	obj, err := c.Get(ctx, id)
	if err != nil {
		panic(err)
	}
	return obj
}

// QueryWorker queries the worker edge of a Chunk.
func (c *ChunkClient) QueryWorker(ch *Chunk) *WorkerQuery {
	query := &WorkerQuery{config: c.config}
	query.path = func(ctx context.Context) (fromV *sql.Selector, _ error) {
		id := ch.ID
		step := sqlgraph.NewStep(
			sqlgraph.From(chunk.Table, chunk.FieldID, id),
			sqlgraph.To(worker.Table, worker.FieldID),
			sqlgraph.Edge(sqlgraph.M2O, true, chunk.WorkerTable, chunk.WorkerColumn),
		)
		fromV = sqlgraph.Neighbors(ch.driver.Dialect(), step)
		return fromV, nil
	}
	return query
}

// Hooks returns the client hooks.
func (c *ChunkClient) Hooks() []Hook {
	return c.hooks.Chunk
}

// WorkerClient is a client for the Worker schema.
type WorkerClient struct {
	config
}

// NewWorkerClient returns a client for the Worker from the given config.
func NewWorkerClient(c config) *WorkerClient {
	return &WorkerClient{config: c}
}

// Use adds a list of mutation hooks to the hooks stack.
// A call to `Use(f, g, h)` equals to `worker.Hooks(f(g(h())))`.
func (c *WorkerClient) Use(hooks ...Hook) {
	c.hooks.Worker = append(c.hooks.Worker, hooks...)
}

// Create returns a builder for creating a Worker entity.
func (c *WorkerClient) Create() *WorkerCreate {
	mutation := newWorkerMutation(c.config, OpCreate)
	return &WorkerCreate{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// CreateBulk returns a builder for creating a bulk of Worker entities.
func (c *WorkerClient) CreateBulk(builders ...*WorkerCreate) *WorkerCreateBulk {
	return &WorkerCreateBulk{config: c.config, builders: builders}
}

// Update returns an update builder for Worker.
func (c *WorkerClient) Update() *WorkerUpdate {
	mutation := newWorkerMutation(c.config, OpUpdate)
	return &WorkerUpdate{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// UpdateOne returns an update builder for the given entity.
func (c *WorkerClient) UpdateOne(w *Worker) *WorkerUpdateOne {
	mutation := newWorkerMutation(c.config, OpUpdateOne, withWorker(w))
	return &WorkerUpdateOne{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// UpdateOneID returns an update builder for the given id.
func (c *WorkerClient) UpdateOneID(id uuid.UUID) *WorkerUpdateOne {
	mutation := newWorkerMutation(c.config, OpUpdateOne, withWorkerID(id))
	return &WorkerUpdateOne{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// Delete returns a delete builder for Worker.
func (c *WorkerClient) Delete() *WorkerDelete {
	mutation := newWorkerMutation(c.config, OpDelete)
	return &WorkerDelete{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// DeleteOne returns a builder for deleting the given entity.
func (c *WorkerClient) DeleteOne(w *Worker) *WorkerDeleteOne {
	return c.DeleteOneID(w.ID)
}

// DeleteOne returns a builder for deleting the given entity by its id.
func (c *WorkerClient) DeleteOneID(id uuid.UUID) *WorkerDeleteOne {
	builder := c.Delete().Where(worker.ID(id))
	builder.mutation.id = &id
	builder.mutation.op = OpDeleteOne
	return &WorkerDeleteOne{builder}
}

// Query returns a query builder for Worker.
func (c *WorkerClient) Query() *WorkerQuery {
	return &WorkerQuery{
		config: c.config,
	}
}

// Get returns a Worker entity by its id.
func (c *WorkerClient) Get(ctx context.Context, id uuid.UUID) (*Worker, error) {
	return c.Query().Where(worker.ID(id)).Only(ctx)
}

// GetX is like Get, but panics if an error occurs.
func (c *WorkerClient) GetX(ctx context.Context, id uuid.UUID) *Worker {
	obj, err := c.Get(ctx, id)
	if err != nil {
		panic(err)
	}
	return obj
}

// QueryChunks queries the chunks edge of a Worker.
func (c *WorkerClient) QueryChunks(w *Worker) *ChunkQuery {
	query := &ChunkQuery{config: c.config}
	query.path = func(ctx context.Context) (fromV *sql.Selector, _ error) {
		id := w.ID
		step := sqlgraph.NewStep(
			sqlgraph.From(worker.Table, worker.FieldID, id),
			sqlgraph.To(chunk.Table, chunk.FieldID),
			sqlgraph.Edge(sqlgraph.O2M, false, worker.ChunksTable, worker.ChunksColumn),
		)
		fromV = sqlgraph.Neighbors(w.driver.Dialect(), step)
		return fromV, nil
	}
	return query
}

// Hooks returns the client hooks.
func (c *WorkerClient) Hooks() []Hook {
	return c.hooks.Worker
}
