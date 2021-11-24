// Code generated by entc, DO NOT EDIT.

package ent

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/go-faster/gha/internal/ent/chunk"
	"github.com/go-faster/gha/internal/ent/download"
	"github.com/go-faster/gha/internal/ent/predicate"
	"github.com/go-faster/gha/internal/ent/worker"
	"github.com/google/uuid"

	"entgo.io/ent"
)

const (
	// Operation types.
	OpCreate    = ent.OpCreate
	OpDelete    = ent.OpDelete
	OpDeleteOne = ent.OpDeleteOne
	OpUpdate    = ent.OpUpdate
	OpUpdateOne = ent.OpUpdateOne

	// Node types.
	TypeChunk    = "Chunk"
	TypeDownload = "Download"
	TypeWorker   = "Worker"
)

// ChunkMutation represents an operation that mutates the Chunk nodes in the graph.
type ChunkMutation struct {
	config
	op            Op
	typ           string
	id            *string
	created_at    *time.Time
	updated_at    *time.Time
	start         *time.Time
	clearedFields map[string]struct{}
	done          bool
	oldValue      func(context.Context) (*Chunk, error)
	predicates    []predicate.Chunk
}

var _ ent.Mutation = (*ChunkMutation)(nil)

// chunkOption allows management of the mutation configuration using functional options.
type chunkOption func(*ChunkMutation)

// newChunkMutation creates new mutation for the Chunk entity.
func newChunkMutation(c config, op Op, opts ...chunkOption) *ChunkMutation {
	m := &ChunkMutation{
		config:        c,
		op:            op,
		typ:           TypeChunk,
		clearedFields: make(map[string]struct{}),
	}
	for _, opt := range opts {
		opt(m)
	}
	return m
}

// withChunkID sets the ID field of the mutation.
func withChunkID(id string) chunkOption {
	return func(m *ChunkMutation) {
		var (
			err   error
			once  sync.Once
			value *Chunk
		)
		m.oldValue = func(ctx context.Context) (*Chunk, error) {
			once.Do(func() {
				if m.done {
					err = fmt.Errorf("querying old values post mutation is not allowed")
				} else {
					value, err = m.Client().Chunk.Get(ctx, id)
				}
			})
			return value, err
		}
		m.id = &id
	}
}

// withChunk sets the old Chunk of the mutation.
func withChunk(node *Chunk) chunkOption {
	return func(m *ChunkMutation) {
		m.oldValue = func(context.Context) (*Chunk, error) {
			return node, nil
		}
		m.id = &node.ID
	}
}

// Client returns a new `ent.Client` from the mutation. If the mutation was
// executed in a transaction (ent.Tx), a transactional client is returned.
func (m ChunkMutation) Client() *Client {
	client := &Client{config: m.config}
	client.init()
	return client
}

// Tx returns an `ent.Tx` for mutations that were executed in transactions;
// it returns an error otherwise.
func (m ChunkMutation) Tx() (*Tx, error) {
	if _, ok := m.driver.(*txDriver); !ok {
		return nil, fmt.Errorf("ent: mutation is not running in a transaction")
	}
	tx := &Tx{config: m.config}
	tx.init()
	return tx, nil
}

// SetID sets the value of the id field. Note that this
// operation is only accepted on creation of Chunk entities.
func (m *ChunkMutation) SetID(id string) {
	m.id = &id
}

// ID returns the ID value in the mutation. Note that the ID is only available
// if it was provided to the builder or after it was returned from the database.
func (m *ChunkMutation) ID() (id string, exists bool) {
	if m.id == nil {
		return
	}
	return *m.id, true
}

// SetCreatedAt sets the "created_at" field.
func (m *ChunkMutation) SetCreatedAt(t time.Time) {
	m.created_at = &t
}

// CreatedAt returns the value of the "created_at" field in the mutation.
func (m *ChunkMutation) CreatedAt() (r time.Time, exists bool) {
	v := m.created_at
	if v == nil {
		return
	}
	return *v, true
}

// OldCreatedAt returns the old "created_at" field's value of the Chunk entity.
// If the Chunk object wasn't provided to the builder, the object is fetched from the database.
// An error is returned if the mutation operation is not UpdateOne, or the database query fails.
func (m *ChunkMutation) OldCreatedAt(ctx context.Context) (v time.Time, err error) {
	if !m.op.Is(OpUpdateOne) {
		return v, fmt.Errorf("OldCreatedAt is only allowed on UpdateOne operations")
	}
	if m.id == nil || m.oldValue == nil {
		return v, fmt.Errorf("OldCreatedAt requires an ID field in the mutation")
	}
	oldValue, err := m.oldValue(ctx)
	if err != nil {
		return v, fmt.Errorf("querying old value for OldCreatedAt: %w", err)
	}
	return oldValue.CreatedAt, nil
}

// ResetCreatedAt resets all changes to the "created_at" field.
func (m *ChunkMutation) ResetCreatedAt() {
	m.created_at = nil
}

// SetUpdatedAt sets the "updated_at" field.
func (m *ChunkMutation) SetUpdatedAt(t time.Time) {
	m.updated_at = &t
}

// UpdatedAt returns the value of the "updated_at" field in the mutation.
func (m *ChunkMutation) UpdatedAt() (r time.Time, exists bool) {
	v := m.updated_at
	if v == nil {
		return
	}
	return *v, true
}

// OldUpdatedAt returns the old "updated_at" field's value of the Chunk entity.
// If the Chunk object wasn't provided to the builder, the object is fetched from the database.
// An error is returned if the mutation operation is not UpdateOne, or the database query fails.
func (m *ChunkMutation) OldUpdatedAt(ctx context.Context) (v time.Time, err error) {
	if !m.op.Is(OpUpdateOne) {
		return v, fmt.Errorf("OldUpdatedAt is only allowed on UpdateOne operations")
	}
	if m.id == nil || m.oldValue == nil {
		return v, fmt.Errorf("OldUpdatedAt requires an ID field in the mutation")
	}
	oldValue, err := m.oldValue(ctx)
	if err != nil {
		return v, fmt.Errorf("querying old value for OldUpdatedAt: %w", err)
	}
	return oldValue.UpdatedAt, nil
}

// ResetUpdatedAt resets all changes to the "updated_at" field.
func (m *ChunkMutation) ResetUpdatedAt() {
	m.updated_at = nil
}

// SetStart sets the "start" field.
func (m *ChunkMutation) SetStart(t time.Time) {
	m.start = &t
}

// Start returns the value of the "start" field in the mutation.
func (m *ChunkMutation) Start() (r time.Time, exists bool) {
	v := m.start
	if v == nil {
		return
	}
	return *v, true
}

// OldStart returns the old "start" field's value of the Chunk entity.
// If the Chunk object wasn't provided to the builder, the object is fetched from the database.
// An error is returned if the mutation operation is not UpdateOne, or the database query fails.
func (m *ChunkMutation) OldStart(ctx context.Context) (v time.Time, err error) {
	if !m.op.Is(OpUpdateOne) {
		return v, fmt.Errorf("OldStart is only allowed on UpdateOne operations")
	}
	if m.id == nil || m.oldValue == nil {
		return v, fmt.Errorf("OldStart requires an ID field in the mutation")
	}
	oldValue, err := m.oldValue(ctx)
	if err != nil {
		return v, fmt.Errorf("querying old value for OldStart: %w", err)
	}
	return oldValue.Start, nil
}

// ResetStart resets all changes to the "start" field.
func (m *ChunkMutation) ResetStart() {
	m.start = nil
}

// Where appends a list predicates to the ChunkMutation builder.
func (m *ChunkMutation) Where(ps ...predicate.Chunk) {
	m.predicates = append(m.predicates, ps...)
}

// Op returns the operation name.
func (m *ChunkMutation) Op() Op {
	return m.op
}

// Type returns the node type of this mutation (Chunk).
func (m *ChunkMutation) Type() string {
	return m.typ
}

// Fields returns all fields that were changed during this mutation. Note that in
// order to get all numeric fields that were incremented/decremented, call
// AddedFields().
func (m *ChunkMutation) Fields() []string {
	fields := make([]string, 0, 3)
	if m.created_at != nil {
		fields = append(fields, chunk.FieldCreatedAt)
	}
	if m.updated_at != nil {
		fields = append(fields, chunk.FieldUpdatedAt)
	}
	if m.start != nil {
		fields = append(fields, chunk.FieldStart)
	}
	return fields
}

// Field returns the value of a field with the given name. The second boolean
// return value indicates that this field was not set, or was not defined in the
// schema.
func (m *ChunkMutation) Field(name string) (ent.Value, bool) {
	switch name {
	case chunk.FieldCreatedAt:
		return m.CreatedAt()
	case chunk.FieldUpdatedAt:
		return m.UpdatedAt()
	case chunk.FieldStart:
		return m.Start()
	}
	return nil, false
}

// OldField returns the old value of the field from the database. An error is
// returned if the mutation operation is not UpdateOne, or the query to the
// database failed.
func (m *ChunkMutation) OldField(ctx context.Context, name string) (ent.Value, error) {
	switch name {
	case chunk.FieldCreatedAt:
		return m.OldCreatedAt(ctx)
	case chunk.FieldUpdatedAt:
		return m.OldUpdatedAt(ctx)
	case chunk.FieldStart:
		return m.OldStart(ctx)
	}
	return nil, fmt.Errorf("unknown Chunk field %s", name)
}

// SetField sets the value of a field with the given name. It returns an error if
// the field is not defined in the schema, or if the type mismatched the field
// type.
func (m *ChunkMutation) SetField(name string, value ent.Value) error {
	switch name {
	case chunk.FieldCreatedAt:
		v, ok := value.(time.Time)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetCreatedAt(v)
		return nil
	case chunk.FieldUpdatedAt:
		v, ok := value.(time.Time)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetUpdatedAt(v)
		return nil
	case chunk.FieldStart:
		v, ok := value.(time.Time)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetStart(v)
		return nil
	}
	return fmt.Errorf("unknown Chunk field %s", name)
}

// AddedFields returns all numeric fields that were incremented/decremented during
// this mutation.
func (m *ChunkMutation) AddedFields() []string {
	return nil
}

// AddedField returns the numeric value that was incremented/decremented on a field
// with the given name. The second boolean return value indicates that this field
// was not set, or was not defined in the schema.
func (m *ChunkMutation) AddedField(name string) (ent.Value, bool) {
	return nil, false
}

// AddField adds the value to the field with the given name. It returns an error if
// the field is not defined in the schema, or if the type mismatched the field
// type.
func (m *ChunkMutation) AddField(name string, value ent.Value) error {
	switch name {
	}
	return fmt.Errorf("unknown Chunk numeric field %s", name)
}

// ClearedFields returns all nullable fields that were cleared during this
// mutation.
func (m *ChunkMutation) ClearedFields() []string {
	return nil
}

// FieldCleared returns a boolean indicating if a field with the given name was
// cleared in this mutation.
func (m *ChunkMutation) FieldCleared(name string) bool {
	_, ok := m.clearedFields[name]
	return ok
}

// ClearField clears the value of the field with the given name. It returns an
// error if the field is not defined in the schema.
func (m *ChunkMutation) ClearField(name string) error {
	return fmt.Errorf("unknown Chunk nullable field %s", name)
}

// ResetField resets all changes in the mutation for the field with the given name.
// It returns an error if the field is not defined in the schema.
func (m *ChunkMutation) ResetField(name string) error {
	switch name {
	case chunk.FieldCreatedAt:
		m.ResetCreatedAt()
		return nil
	case chunk.FieldUpdatedAt:
		m.ResetUpdatedAt()
		return nil
	case chunk.FieldStart:
		m.ResetStart()
		return nil
	}
	return fmt.Errorf("unknown Chunk field %s", name)
}

// AddedEdges returns all edge names that were set/added in this mutation.
func (m *ChunkMutation) AddedEdges() []string {
	edges := make([]string, 0, 0)
	return edges
}

// AddedIDs returns all IDs (to other nodes) that were added for the given edge
// name in this mutation.
func (m *ChunkMutation) AddedIDs(name string) []ent.Value {
	return nil
}

// RemovedEdges returns all edge names that were removed in this mutation.
func (m *ChunkMutation) RemovedEdges() []string {
	edges := make([]string, 0, 0)
	return edges
}

// RemovedIDs returns all IDs (to other nodes) that were removed for the edge with
// the given name in this mutation.
func (m *ChunkMutation) RemovedIDs(name string) []ent.Value {
	return nil
}

// ClearedEdges returns all edge names that were cleared in this mutation.
func (m *ChunkMutation) ClearedEdges() []string {
	edges := make([]string, 0, 0)
	return edges
}

// EdgeCleared returns a boolean which indicates if the edge with the given name
// was cleared in this mutation.
func (m *ChunkMutation) EdgeCleared(name string) bool {
	return false
}

// ClearEdge clears the value of the edge with the given name. It returns an error
// if that edge is not defined in the schema.
func (m *ChunkMutation) ClearEdge(name string) error {
	return fmt.Errorf("unknown Chunk unique edge %s", name)
}

// ResetEdge resets all changes to the edge with the given name in this mutation.
// It returns an error if the edge is not defined in the schema.
func (m *ChunkMutation) ResetEdge(name string) error {
	return fmt.Errorf("unknown Chunk edge %s", name)
}

// DownloadMutation represents an operation that mutates the Download nodes in the graph.
type DownloadMutation struct {
	config
	op            Op
	typ           string
	id            *uuid.UUID
	created_at    *time.Time
	updated_at    *time.Time
	clearedFields map[string]struct{}
	done          bool
	oldValue      func(context.Context) (*Download, error)
	predicates    []predicate.Download
}

var _ ent.Mutation = (*DownloadMutation)(nil)

// downloadOption allows management of the mutation configuration using functional options.
type downloadOption func(*DownloadMutation)

// newDownloadMutation creates new mutation for the Download entity.
func newDownloadMutation(c config, op Op, opts ...downloadOption) *DownloadMutation {
	m := &DownloadMutation{
		config:        c,
		op:            op,
		typ:           TypeDownload,
		clearedFields: make(map[string]struct{}),
	}
	for _, opt := range opts {
		opt(m)
	}
	return m
}

// withDownloadID sets the ID field of the mutation.
func withDownloadID(id uuid.UUID) downloadOption {
	return func(m *DownloadMutation) {
		var (
			err   error
			once  sync.Once
			value *Download
		)
		m.oldValue = func(ctx context.Context) (*Download, error) {
			once.Do(func() {
				if m.done {
					err = fmt.Errorf("querying old values post mutation is not allowed")
				} else {
					value, err = m.Client().Download.Get(ctx, id)
				}
			})
			return value, err
		}
		m.id = &id
	}
}

// withDownload sets the old Download of the mutation.
func withDownload(node *Download) downloadOption {
	return func(m *DownloadMutation) {
		m.oldValue = func(context.Context) (*Download, error) {
			return node, nil
		}
		m.id = &node.ID
	}
}

// Client returns a new `ent.Client` from the mutation. If the mutation was
// executed in a transaction (ent.Tx), a transactional client is returned.
func (m DownloadMutation) Client() *Client {
	client := &Client{config: m.config}
	client.init()
	return client
}

// Tx returns an `ent.Tx` for mutations that were executed in transactions;
// it returns an error otherwise.
func (m DownloadMutation) Tx() (*Tx, error) {
	if _, ok := m.driver.(*txDriver); !ok {
		return nil, fmt.Errorf("ent: mutation is not running in a transaction")
	}
	tx := &Tx{config: m.config}
	tx.init()
	return tx, nil
}

// SetID sets the value of the id field. Note that this
// operation is only accepted on creation of Download entities.
func (m *DownloadMutation) SetID(id uuid.UUID) {
	m.id = &id
}

// ID returns the ID value in the mutation. Note that the ID is only available
// if it was provided to the builder or after it was returned from the database.
func (m *DownloadMutation) ID() (id uuid.UUID, exists bool) {
	if m.id == nil {
		return
	}
	return *m.id, true
}

// SetCreatedAt sets the "created_at" field.
func (m *DownloadMutation) SetCreatedAt(t time.Time) {
	m.created_at = &t
}

// CreatedAt returns the value of the "created_at" field in the mutation.
func (m *DownloadMutation) CreatedAt() (r time.Time, exists bool) {
	v := m.created_at
	if v == nil {
		return
	}
	return *v, true
}

// OldCreatedAt returns the old "created_at" field's value of the Download entity.
// If the Download object wasn't provided to the builder, the object is fetched from the database.
// An error is returned if the mutation operation is not UpdateOne, or the database query fails.
func (m *DownloadMutation) OldCreatedAt(ctx context.Context) (v time.Time, err error) {
	if !m.op.Is(OpUpdateOne) {
		return v, fmt.Errorf("OldCreatedAt is only allowed on UpdateOne operations")
	}
	if m.id == nil || m.oldValue == nil {
		return v, fmt.Errorf("OldCreatedAt requires an ID field in the mutation")
	}
	oldValue, err := m.oldValue(ctx)
	if err != nil {
		return v, fmt.Errorf("querying old value for OldCreatedAt: %w", err)
	}
	return oldValue.CreatedAt, nil
}

// ResetCreatedAt resets all changes to the "created_at" field.
func (m *DownloadMutation) ResetCreatedAt() {
	m.created_at = nil
}

// SetUpdatedAt sets the "updated_at" field.
func (m *DownloadMutation) SetUpdatedAt(t time.Time) {
	m.updated_at = &t
}

// UpdatedAt returns the value of the "updated_at" field in the mutation.
func (m *DownloadMutation) UpdatedAt() (r time.Time, exists bool) {
	v := m.updated_at
	if v == nil {
		return
	}
	return *v, true
}

// OldUpdatedAt returns the old "updated_at" field's value of the Download entity.
// If the Download object wasn't provided to the builder, the object is fetched from the database.
// An error is returned if the mutation operation is not UpdateOne, or the database query fails.
func (m *DownloadMutation) OldUpdatedAt(ctx context.Context) (v time.Time, err error) {
	if !m.op.Is(OpUpdateOne) {
		return v, fmt.Errorf("OldUpdatedAt is only allowed on UpdateOne operations")
	}
	if m.id == nil || m.oldValue == nil {
		return v, fmt.Errorf("OldUpdatedAt requires an ID field in the mutation")
	}
	oldValue, err := m.oldValue(ctx)
	if err != nil {
		return v, fmt.Errorf("querying old value for OldUpdatedAt: %w", err)
	}
	return oldValue.UpdatedAt, nil
}

// ResetUpdatedAt resets all changes to the "updated_at" field.
func (m *DownloadMutation) ResetUpdatedAt() {
	m.updated_at = nil
}

// Where appends a list predicates to the DownloadMutation builder.
func (m *DownloadMutation) Where(ps ...predicate.Download) {
	m.predicates = append(m.predicates, ps...)
}

// Op returns the operation name.
func (m *DownloadMutation) Op() Op {
	return m.op
}

// Type returns the node type of this mutation (Download).
func (m *DownloadMutation) Type() string {
	return m.typ
}

// Fields returns all fields that were changed during this mutation. Note that in
// order to get all numeric fields that were incremented/decremented, call
// AddedFields().
func (m *DownloadMutation) Fields() []string {
	fields := make([]string, 0, 2)
	if m.created_at != nil {
		fields = append(fields, download.FieldCreatedAt)
	}
	if m.updated_at != nil {
		fields = append(fields, download.FieldUpdatedAt)
	}
	return fields
}

// Field returns the value of a field with the given name. The second boolean
// return value indicates that this field was not set, or was not defined in the
// schema.
func (m *DownloadMutation) Field(name string) (ent.Value, bool) {
	switch name {
	case download.FieldCreatedAt:
		return m.CreatedAt()
	case download.FieldUpdatedAt:
		return m.UpdatedAt()
	}
	return nil, false
}

// OldField returns the old value of the field from the database. An error is
// returned if the mutation operation is not UpdateOne, or the query to the
// database failed.
func (m *DownloadMutation) OldField(ctx context.Context, name string) (ent.Value, error) {
	switch name {
	case download.FieldCreatedAt:
		return m.OldCreatedAt(ctx)
	case download.FieldUpdatedAt:
		return m.OldUpdatedAt(ctx)
	}
	return nil, fmt.Errorf("unknown Download field %s", name)
}

// SetField sets the value of a field with the given name. It returns an error if
// the field is not defined in the schema, or if the type mismatched the field
// type.
func (m *DownloadMutation) SetField(name string, value ent.Value) error {
	switch name {
	case download.FieldCreatedAt:
		v, ok := value.(time.Time)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetCreatedAt(v)
		return nil
	case download.FieldUpdatedAt:
		v, ok := value.(time.Time)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetUpdatedAt(v)
		return nil
	}
	return fmt.Errorf("unknown Download field %s", name)
}

// AddedFields returns all numeric fields that were incremented/decremented during
// this mutation.
func (m *DownloadMutation) AddedFields() []string {
	return nil
}

// AddedField returns the numeric value that was incremented/decremented on a field
// with the given name. The second boolean return value indicates that this field
// was not set, or was not defined in the schema.
func (m *DownloadMutation) AddedField(name string) (ent.Value, bool) {
	return nil, false
}

// AddField adds the value to the field with the given name. It returns an error if
// the field is not defined in the schema, or if the type mismatched the field
// type.
func (m *DownloadMutation) AddField(name string, value ent.Value) error {
	switch name {
	}
	return fmt.Errorf("unknown Download numeric field %s", name)
}

// ClearedFields returns all nullable fields that were cleared during this
// mutation.
func (m *DownloadMutation) ClearedFields() []string {
	return nil
}

// FieldCleared returns a boolean indicating if a field with the given name was
// cleared in this mutation.
func (m *DownloadMutation) FieldCleared(name string) bool {
	_, ok := m.clearedFields[name]
	return ok
}

// ClearField clears the value of the field with the given name. It returns an
// error if the field is not defined in the schema.
func (m *DownloadMutation) ClearField(name string) error {
	return fmt.Errorf("unknown Download nullable field %s", name)
}

// ResetField resets all changes in the mutation for the field with the given name.
// It returns an error if the field is not defined in the schema.
func (m *DownloadMutation) ResetField(name string) error {
	switch name {
	case download.FieldCreatedAt:
		m.ResetCreatedAt()
		return nil
	case download.FieldUpdatedAt:
		m.ResetUpdatedAt()
		return nil
	}
	return fmt.Errorf("unknown Download field %s", name)
}

// AddedEdges returns all edge names that were set/added in this mutation.
func (m *DownloadMutation) AddedEdges() []string {
	edges := make([]string, 0, 0)
	return edges
}

// AddedIDs returns all IDs (to other nodes) that were added for the given edge
// name in this mutation.
func (m *DownloadMutation) AddedIDs(name string) []ent.Value {
	return nil
}

// RemovedEdges returns all edge names that were removed in this mutation.
func (m *DownloadMutation) RemovedEdges() []string {
	edges := make([]string, 0, 0)
	return edges
}

// RemovedIDs returns all IDs (to other nodes) that were removed for the edge with
// the given name in this mutation.
func (m *DownloadMutation) RemovedIDs(name string) []ent.Value {
	return nil
}

// ClearedEdges returns all edge names that were cleared in this mutation.
func (m *DownloadMutation) ClearedEdges() []string {
	edges := make([]string, 0, 0)
	return edges
}

// EdgeCleared returns a boolean which indicates if the edge with the given name
// was cleared in this mutation.
func (m *DownloadMutation) EdgeCleared(name string) bool {
	return false
}

// ClearEdge clears the value of the edge with the given name. It returns an error
// if that edge is not defined in the schema.
func (m *DownloadMutation) ClearEdge(name string) error {
	return fmt.Errorf("unknown Download unique edge %s", name)
}

// ResetEdge resets all changes to the edge with the given name in this mutation.
// It returns an error if the edge is not defined in the schema.
func (m *DownloadMutation) ResetEdge(name string) error {
	return fmt.Errorf("unknown Download edge %s", name)
}

// WorkerMutation represents an operation that mutates the Worker nodes in the graph.
type WorkerMutation struct {
	config
	op            Op
	typ           string
	id            *uuid.UUID
	created_at    *time.Time
	updated_at    *time.Time
	name          *string
	token         *string
	clearedFields map[string]struct{}
	done          bool
	oldValue      func(context.Context) (*Worker, error)
	predicates    []predicate.Worker
}

var _ ent.Mutation = (*WorkerMutation)(nil)

// workerOption allows management of the mutation configuration using functional options.
type workerOption func(*WorkerMutation)

// newWorkerMutation creates new mutation for the Worker entity.
func newWorkerMutation(c config, op Op, opts ...workerOption) *WorkerMutation {
	m := &WorkerMutation{
		config:        c,
		op:            op,
		typ:           TypeWorker,
		clearedFields: make(map[string]struct{}),
	}
	for _, opt := range opts {
		opt(m)
	}
	return m
}

// withWorkerID sets the ID field of the mutation.
func withWorkerID(id uuid.UUID) workerOption {
	return func(m *WorkerMutation) {
		var (
			err   error
			once  sync.Once
			value *Worker
		)
		m.oldValue = func(ctx context.Context) (*Worker, error) {
			once.Do(func() {
				if m.done {
					err = fmt.Errorf("querying old values post mutation is not allowed")
				} else {
					value, err = m.Client().Worker.Get(ctx, id)
				}
			})
			return value, err
		}
		m.id = &id
	}
}

// withWorker sets the old Worker of the mutation.
func withWorker(node *Worker) workerOption {
	return func(m *WorkerMutation) {
		m.oldValue = func(context.Context) (*Worker, error) {
			return node, nil
		}
		m.id = &node.ID
	}
}

// Client returns a new `ent.Client` from the mutation. If the mutation was
// executed in a transaction (ent.Tx), a transactional client is returned.
func (m WorkerMutation) Client() *Client {
	client := &Client{config: m.config}
	client.init()
	return client
}

// Tx returns an `ent.Tx` for mutations that were executed in transactions;
// it returns an error otherwise.
func (m WorkerMutation) Tx() (*Tx, error) {
	if _, ok := m.driver.(*txDriver); !ok {
		return nil, fmt.Errorf("ent: mutation is not running in a transaction")
	}
	tx := &Tx{config: m.config}
	tx.init()
	return tx, nil
}

// SetID sets the value of the id field. Note that this
// operation is only accepted on creation of Worker entities.
func (m *WorkerMutation) SetID(id uuid.UUID) {
	m.id = &id
}

// ID returns the ID value in the mutation. Note that the ID is only available
// if it was provided to the builder or after it was returned from the database.
func (m *WorkerMutation) ID() (id uuid.UUID, exists bool) {
	if m.id == nil {
		return
	}
	return *m.id, true
}

// SetCreatedAt sets the "created_at" field.
func (m *WorkerMutation) SetCreatedAt(t time.Time) {
	m.created_at = &t
}

// CreatedAt returns the value of the "created_at" field in the mutation.
func (m *WorkerMutation) CreatedAt() (r time.Time, exists bool) {
	v := m.created_at
	if v == nil {
		return
	}
	return *v, true
}

// OldCreatedAt returns the old "created_at" field's value of the Worker entity.
// If the Worker object wasn't provided to the builder, the object is fetched from the database.
// An error is returned if the mutation operation is not UpdateOne, or the database query fails.
func (m *WorkerMutation) OldCreatedAt(ctx context.Context) (v time.Time, err error) {
	if !m.op.Is(OpUpdateOne) {
		return v, fmt.Errorf("OldCreatedAt is only allowed on UpdateOne operations")
	}
	if m.id == nil || m.oldValue == nil {
		return v, fmt.Errorf("OldCreatedAt requires an ID field in the mutation")
	}
	oldValue, err := m.oldValue(ctx)
	if err != nil {
		return v, fmt.Errorf("querying old value for OldCreatedAt: %w", err)
	}
	return oldValue.CreatedAt, nil
}

// ResetCreatedAt resets all changes to the "created_at" field.
func (m *WorkerMutation) ResetCreatedAt() {
	m.created_at = nil
}

// SetUpdatedAt sets the "updated_at" field.
func (m *WorkerMutation) SetUpdatedAt(t time.Time) {
	m.updated_at = &t
}

// UpdatedAt returns the value of the "updated_at" field in the mutation.
func (m *WorkerMutation) UpdatedAt() (r time.Time, exists bool) {
	v := m.updated_at
	if v == nil {
		return
	}
	return *v, true
}

// OldUpdatedAt returns the old "updated_at" field's value of the Worker entity.
// If the Worker object wasn't provided to the builder, the object is fetched from the database.
// An error is returned if the mutation operation is not UpdateOne, or the database query fails.
func (m *WorkerMutation) OldUpdatedAt(ctx context.Context) (v time.Time, err error) {
	if !m.op.Is(OpUpdateOne) {
		return v, fmt.Errorf("OldUpdatedAt is only allowed on UpdateOne operations")
	}
	if m.id == nil || m.oldValue == nil {
		return v, fmt.Errorf("OldUpdatedAt requires an ID field in the mutation")
	}
	oldValue, err := m.oldValue(ctx)
	if err != nil {
		return v, fmt.Errorf("querying old value for OldUpdatedAt: %w", err)
	}
	return oldValue.UpdatedAt, nil
}

// ResetUpdatedAt resets all changes to the "updated_at" field.
func (m *WorkerMutation) ResetUpdatedAt() {
	m.updated_at = nil
}

// SetName sets the "name" field.
func (m *WorkerMutation) SetName(s string) {
	m.name = &s
}

// Name returns the value of the "name" field in the mutation.
func (m *WorkerMutation) Name() (r string, exists bool) {
	v := m.name
	if v == nil {
		return
	}
	return *v, true
}

// OldName returns the old "name" field's value of the Worker entity.
// If the Worker object wasn't provided to the builder, the object is fetched from the database.
// An error is returned if the mutation operation is not UpdateOne, or the database query fails.
func (m *WorkerMutation) OldName(ctx context.Context) (v string, err error) {
	if !m.op.Is(OpUpdateOne) {
		return v, fmt.Errorf("OldName is only allowed on UpdateOne operations")
	}
	if m.id == nil || m.oldValue == nil {
		return v, fmt.Errorf("OldName requires an ID field in the mutation")
	}
	oldValue, err := m.oldValue(ctx)
	if err != nil {
		return v, fmt.Errorf("querying old value for OldName: %w", err)
	}
	return oldValue.Name, nil
}

// ResetName resets all changes to the "name" field.
func (m *WorkerMutation) ResetName() {
	m.name = nil
}

// SetToken sets the "token" field.
func (m *WorkerMutation) SetToken(s string) {
	m.token = &s
}

// Token returns the value of the "token" field in the mutation.
func (m *WorkerMutation) Token() (r string, exists bool) {
	v := m.token
	if v == nil {
		return
	}
	return *v, true
}

// OldToken returns the old "token" field's value of the Worker entity.
// If the Worker object wasn't provided to the builder, the object is fetched from the database.
// An error is returned if the mutation operation is not UpdateOne, or the database query fails.
func (m *WorkerMutation) OldToken(ctx context.Context) (v string, err error) {
	if !m.op.Is(OpUpdateOne) {
		return v, fmt.Errorf("OldToken is only allowed on UpdateOne operations")
	}
	if m.id == nil || m.oldValue == nil {
		return v, fmt.Errorf("OldToken requires an ID field in the mutation")
	}
	oldValue, err := m.oldValue(ctx)
	if err != nil {
		return v, fmt.Errorf("querying old value for OldToken: %w", err)
	}
	return oldValue.Token, nil
}

// ResetToken resets all changes to the "token" field.
func (m *WorkerMutation) ResetToken() {
	m.token = nil
}

// Where appends a list predicates to the WorkerMutation builder.
func (m *WorkerMutation) Where(ps ...predicate.Worker) {
	m.predicates = append(m.predicates, ps...)
}

// Op returns the operation name.
func (m *WorkerMutation) Op() Op {
	return m.op
}

// Type returns the node type of this mutation (Worker).
func (m *WorkerMutation) Type() string {
	return m.typ
}

// Fields returns all fields that were changed during this mutation. Note that in
// order to get all numeric fields that were incremented/decremented, call
// AddedFields().
func (m *WorkerMutation) Fields() []string {
	fields := make([]string, 0, 4)
	if m.created_at != nil {
		fields = append(fields, worker.FieldCreatedAt)
	}
	if m.updated_at != nil {
		fields = append(fields, worker.FieldUpdatedAt)
	}
	if m.name != nil {
		fields = append(fields, worker.FieldName)
	}
	if m.token != nil {
		fields = append(fields, worker.FieldToken)
	}
	return fields
}

// Field returns the value of a field with the given name. The second boolean
// return value indicates that this field was not set, or was not defined in the
// schema.
func (m *WorkerMutation) Field(name string) (ent.Value, bool) {
	switch name {
	case worker.FieldCreatedAt:
		return m.CreatedAt()
	case worker.FieldUpdatedAt:
		return m.UpdatedAt()
	case worker.FieldName:
		return m.Name()
	case worker.FieldToken:
		return m.Token()
	}
	return nil, false
}

// OldField returns the old value of the field from the database. An error is
// returned if the mutation operation is not UpdateOne, or the query to the
// database failed.
func (m *WorkerMutation) OldField(ctx context.Context, name string) (ent.Value, error) {
	switch name {
	case worker.FieldCreatedAt:
		return m.OldCreatedAt(ctx)
	case worker.FieldUpdatedAt:
		return m.OldUpdatedAt(ctx)
	case worker.FieldName:
		return m.OldName(ctx)
	case worker.FieldToken:
		return m.OldToken(ctx)
	}
	return nil, fmt.Errorf("unknown Worker field %s", name)
}

// SetField sets the value of a field with the given name. It returns an error if
// the field is not defined in the schema, or if the type mismatched the field
// type.
func (m *WorkerMutation) SetField(name string, value ent.Value) error {
	switch name {
	case worker.FieldCreatedAt:
		v, ok := value.(time.Time)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetCreatedAt(v)
		return nil
	case worker.FieldUpdatedAt:
		v, ok := value.(time.Time)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetUpdatedAt(v)
		return nil
	case worker.FieldName:
		v, ok := value.(string)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetName(v)
		return nil
	case worker.FieldToken:
		v, ok := value.(string)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetToken(v)
		return nil
	}
	return fmt.Errorf("unknown Worker field %s", name)
}

// AddedFields returns all numeric fields that were incremented/decremented during
// this mutation.
func (m *WorkerMutation) AddedFields() []string {
	return nil
}

// AddedField returns the numeric value that was incremented/decremented on a field
// with the given name. The second boolean return value indicates that this field
// was not set, or was not defined in the schema.
func (m *WorkerMutation) AddedField(name string) (ent.Value, bool) {
	return nil, false
}

// AddField adds the value to the field with the given name. It returns an error if
// the field is not defined in the schema, or if the type mismatched the field
// type.
func (m *WorkerMutation) AddField(name string, value ent.Value) error {
	switch name {
	}
	return fmt.Errorf("unknown Worker numeric field %s", name)
}

// ClearedFields returns all nullable fields that were cleared during this
// mutation.
func (m *WorkerMutation) ClearedFields() []string {
	return nil
}

// FieldCleared returns a boolean indicating if a field with the given name was
// cleared in this mutation.
func (m *WorkerMutation) FieldCleared(name string) bool {
	_, ok := m.clearedFields[name]
	return ok
}

// ClearField clears the value of the field with the given name. It returns an
// error if the field is not defined in the schema.
func (m *WorkerMutation) ClearField(name string) error {
	return fmt.Errorf("unknown Worker nullable field %s", name)
}

// ResetField resets all changes in the mutation for the field with the given name.
// It returns an error if the field is not defined in the schema.
func (m *WorkerMutation) ResetField(name string) error {
	switch name {
	case worker.FieldCreatedAt:
		m.ResetCreatedAt()
		return nil
	case worker.FieldUpdatedAt:
		m.ResetUpdatedAt()
		return nil
	case worker.FieldName:
		m.ResetName()
		return nil
	case worker.FieldToken:
		m.ResetToken()
		return nil
	}
	return fmt.Errorf("unknown Worker field %s", name)
}

// AddedEdges returns all edge names that were set/added in this mutation.
func (m *WorkerMutation) AddedEdges() []string {
	edges := make([]string, 0, 0)
	return edges
}

// AddedIDs returns all IDs (to other nodes) that were added for the given edge
// name in this mutation.
func (m *WorkerMutation) AddedIDs(name string) []ent.Value {
	return nil
}

// RemovedEdges returns all edge names that were removed in this mutation.
func (m *WorkerMutation) RemovedEdges() []string {
	edges := make([]string, 0, 0)
	return edges
}

// RemovedIDs returns all IDs (to other nodes) that were removed for the edge with
// the given name in this mutation.
func (m *WorkerMutation) RemovedIDs(name string) []ent.Value {
	return nil
}

// ClearedEdges returns all edge names that were cleared in this mutation.
func (m *WorkerMutation) ClearedEdges() []string {
	edges := make([]string, 0, 0)
	return edges
}

// EdgeCleared returns a boolean which indicates if the edge with the given name
// was cleared in this mutation.
func (m *WorkerMutation) EdgeCleared(name string) bool {
	return false
}

// ClearEdge clears the value of the edge with the given name. It returns an error
// if that edge is not defined in the schema.
func (m *WorkerMutation) ClearEdge(name string) error {
	return fmt.Errorf("unknown Worker unique edge %s", name)
}

// ResetEdge resets all changes to the edge with the given name in this mutation.
// It returns an error if the edge is not defined in the schema.
func (m *WorkerMutation) ResetEdge(name string) error {
	return fmt.Errorf("unknown Worker edge %s", name)
}
