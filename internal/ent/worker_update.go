// Code generated by ent, DO NOT EDIT.

package ent

import (
	"context"
	"errors"
	"fmt"
	"time"

	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"entgo.io/ent/schema/field"
	"github.com/go-faster/gha/internal/ent/chunk"
	"github.com/go-faster/gha/internal/ent/predicate"
	"github.com/go-faster/gha/internal/ent/worker"
)

// WorkerUpdate is the builder for updating Worker entities.
type WorkerUpdate struct {
	config
	hooks    []Hook
	mutation *WorkerMutation
}

// Where appends a list predicates to the WorkerUpdate builder.
func (wu *WorkerUpdate) Where(ps ...predicate.Worker) *WorkerUpdate {
	wu.mutation.Where(ps...)
	return wu
}

// SetUpdatedAt sets the "updated_at" field.
func (wu *WorkerUpdate) SetUpdatedAt(t time.Time) *WorkerUpdate {
	wu.mutation.SetUpdatedAt(t)
	return wu
}

// SetName sets the "name" field.
func (wu *WorkerUpdate) SetName(s string) *WorkerUpdate {
	wu.mutation.SetName(s)
	return wu
}

// SetToken sets the "token" field.
func (wu *WorkerUpdate) SetToken(s string) *WorkerUpdate {
	wu.mutation.SetToken(s)
	return wu
}

// AddChunkIDs adds the "chunks" edge to the Chunk entity by IDs.
func (wu *WorkerUpdate) AddChunkIDs(ids ...string) *WorkerUpdate {
	wu.mutation.AddChunkIDs(ids...)
	return wu
}

// AddChunks adds the "chunks" edges to the Chunk entity.
func (wu *WorkerUpdate) AddChunks(c ...*Chunk) *WorkerUpdate {
	ids := make([]string, len(c))
	for i := range c {
		ids[i] = c[i].ID
	}
	return wu.AddChunkIDs(ids...)
}

// Mutation returns the WorkerMutation object of the builder.
func (wu *WorkerUpdate) Mutation() *WorkerMutation {
	return wu.mutation
}

// ClearChunks clears all "chunks" edges to the Chunk entity.
func (wu *WorkerUpdate) ClearChunks() *WorkerUpdate {
	wu.mutation.ClearChunks()
	return wu
}

// RemoveChunkIDs removes the "chunks" edge to Chunk entities by IDs.
func (wu *WorkerUpdate) RemoveChunkIDs(ids ...string) *WorkerUpdate {
	wu.mutation.RemoveChunkIDs(ids...)
	return wu
}

// RemoveChunks removes "chunks" edges to Chunk entities.
func (wu *WorkerUpdate) RemoveChunks(c ...*Chunk) *WorkerUpdate {
	ids := make([]string, len(c))
	for i := range c {
		ids[i] = c[i].ID
	}
	return wu.RemoveChunkIDs(ids...)
}

// Save executes the query and returns the number of nodes affected by the update operation.
func (wu *WorkerUpdate) Save(ctx context.Context) (int, error) {
	wu.defaults()
	return withHooks(ctx, wu.sqlSave, wu.mutation, wu.hooks)
}

// SaveX is like Save, but panics if an error occurs.
func (wu *WorkerUpdate) SaveX(ctx context.Context) int {
	affected, err := wu.Save(ctx)
	if err != nil {
		panic(err)
	}
	return affected
}

// Exec executes the query.
func (wu *WorkerUpdate) Exec(ctx context.Context) error {
	_, err := wu.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (wu *WorkerUpdate) ExecX(ctx context.Context) {
	if err := wu.Exec(ctx); err != nil {
		panic(err)
	}
}

// defaults sets the default values of the builder before save.
func (wu *WorkerUpdate) defaults() {
	if _, ok := wu.mutation.UpdatedAt(); !ok {
		v := worker.UpdateDefaultUpdatedAt()
		wu.mutation.SetUpdatedAt(v)
	}
}

func (wu *WorkerUpdate) sqlSave(ctx context.Context) (n int, err error) {
	_spec := sqlgraph.NewUpdateSpec(worker.Table, worker.Columns, sqlgraph.NewFieldSpec(worker.FieldID, field.TypeUUID))
	if ps := wu.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if value, ok := wu.mutation.UpdatedAt(); ok {
		_spec.SetField(worker.FieldUpdatedAt, field.TypeTime, value)
	}
	if value, ok := wu.mutation.Name(); ok {
		_spec.SetField(worker.FieldName, field.TypeString, value)
	}
	if value, ok := wu.mutation.Token(); ok {
		_spec.SetField(worker.FieldToken, field.TypeString, value)
	}
	if wu.mutation.ChunksCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   worker.ChunksTable,
			Columns: []string{worker.ChunksColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(chunk.FieldID, field.TypeString),
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := wu.mutation.RemovedChunksIDs(); len(nodes) > 0 && !wu.mutation.ChunksCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   worker.ChunksTable,
			Columns: []string{worker.ChunksColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(chunk.FieldID, field.TypeString),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := wu.mutation.ChunksIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   worker.ChunksTable,
			Columns: []string{worker.ChunksColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(chunk.FieldID, field.TypeString),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if n, err = sqlgraph.UpdateNodes(ctx, wu.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{worker.Label}
		} else if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{msg: err.Error(), wrap: err}
		}
		return 0, err
	}
	wu.mutation.done = true
	return n, nil
}

// WorkerUpdateOne is the builder for updating a single Worker entity.
type WorkerUpdateOne struct {
	config
	fields   []string
	hooks    []Hook
	mutation *WorkerMutation
}

// SetUpdatedAt sets the "updated_at" field.
func (wuo *WorkerUpdateOne) SetUpdatedAt(t time.Time) *WorkerUpdateOne {
	wuo.mutation.SetUpdatedAt(t)
	return wuo
}

// SetName sets the "name" field.
func (wuo *WorkerUpdateOne) SetName(s string) *WorkerUpdateOne {
	wuo.mutation.SetName(s)
	return wuo
}

// SetToken sets the "token" field.
func (wuo *WorkerUpdateOne) SetToken(s string) *WorkerUpdateOne {
	wuo.mutation.SetToken(s)
	return wuo
}

// AddChunkIDs adds the "chunks" edge to the Chunk entity by IDs.
func (wuo *WorkerUpdateOne) AddChunkIDs(ids ...string) *WorkerUpdateOne {
	wuo.mutation.AddChunkIDs(ids...)
	return wuo
}

// AddChunks adds the "chunks" edges to the Chunk entity.
func (wuo *WorkerUpdateOne) AddChunks(c ...*Chunk) *WorkerUpdateOne {
	ids := make([]string, len(c))
	for i := range c {
		ids[i] = c[i].ID
	}
	return wuo.AddChunkIDs(ids...)
}

// Mutation returns the WorkerMutation object of the builder.
func (wuo *WorkerUpdateOne) Mutation() *WorkerMutation {
	return wuo.mutation
}

// ClearChunks clears all "chunks" edges to the Chunk entity.
func (wuo *WorkerUpdateOne) ClearChunks() *WorkerUpdateOne {
	wuo.mutation.ClearChunks()
	return wuo
}

// RemoveChunkIDs removes the "chunks" edge to Chunk entities by IDs.
func (wuo *WorkerUpdateOne) RemoveChunkIDs(ids ...string) *WorkerUpdateOne {
	wuo.mutation.RemoveChunkIDs(ids...)
	return wuo
}

// RemoveChunks removes "chunks" edges to Chunk entities.
func (wuo *WorkerUpdateOne) RemoveChunks(c ...*Chunk) *WorkerUpdateOne {
	ids := make([]string, len(c))
	for i := range c {
		ids[i] = c[i].ID
	}
	return wuo.RemoveChunkIDs(ids...)
}

// Where appends a list predicates to the WorkerUpdate builder.
func (wuo *WorkerUpdateOne) Where(ps ...predicate.Worker) *WorkerUpdateOne {
	wuo.mutation.Where(ps...)
	return wuo
}

// Select allows selecting one or more fields (columns) of the returned entity.
// The default is selecting all fields defined in the entity schema.
func (wuo *WorkerUpdateOne) Select(field string, fields ...string) *WorkerUpdateOne {
	wuo.fields = append([]string{field}, fields...)
	return wuo
}

// Save executes the query and returns the updated Worker entity.
func (wuo *WorkerUpdateOne) Save(ctx context.Context) (*Worker, error) {
	wuo.defaults()
	return withHooks(ctx, wuo.sqlSave, wuo.mutation, wuo.hooks)
}

// SaveX is like Save, but panics if an error occurs.
func (wuo *WorkerUpdateOne) SaveX(ctx context.Context) *Worker {
	node, err := wuo.Save(ctx)
	if err != nil {
		panic(err)
	}
	return node
}

// Exec executes the query on the entity.
func (wuo *WorkerUpdateOne) Exec(ctx context.Context) error {
	_, err := wuo.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (wuo *WorkerUpdateOne) ExecX(ctx context.Context) {
	if err := wuo.Exec(ctx); err != nil {
		panic(err)
	}
}

// defaults sets the default values of the builder before save.
func (wuo *WorkerUpdateOne) defaults() {
	if _, ok := wuo.mutation.UpdatedAt(); !ok {
		v := worker.UpdateDefaultUpdatedAt()
		wuo.mutation.SetUpdatedAt(v)
	}
}

func (wuo *WorkerUpdateOne) sqlSave(ctx context.Context) (_node *Worker, err error) {
	_spec := sqlgraph.NewUpdateSpec(worker.Table, worker.Columns, sqlgraph.NewFieldSpec(worker.FieldID, field.TypeUUID))
	id, ok := wuo.mutation.ID()
	if !ok {
		return nil, &ValidationError{Name: "id", err: errors.New(`ent: missing "Worker.id" for update`)}
	}
	_spec.Node.ID.Value = id
	if fields := wuo.fields; len(fields) > 0 {
		_spec.Node.Columns = make([]string, 0, len(fields))
		_spec.Node.Columns = append(_spec.Node.Columns, worker.FieldID)
		for _, f := range fields {
			if !worker.ValidColumn(f) {
				return nil, &ValidationError{Name: f, err: fmt.Errorf("ent: invalid field %q for query", f)}
			}
			if f != worker.FieldID {
				_spec.Node.Columns = append(_spec.Node.Columns, f)
			}
		}
	}
	if ps := wuo.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if value, ok := wuo.mutation.UpdatedAt(); ok {
		_spec.SetField(worker.FieldUpdatedAt, field.TypeTime, value)
	}
	if value, ok := wuo.mutation.Name(); ok {
		_spec.SetField(worker.FieldName, field.TypeString, value)
	}
	if value, ok := wuo.mutation.Token(); ok {
		_spec.SetField(worker.FieldToken, field.TypeString, value)
	}
	if wuo.mutation.ChunksCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   worker.ChunksTable,
			Columns: []string{worker.ChunksColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(chunk.FieldID, field.TypeString),
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := wuo.mutation.RemovedChunksIDs(); len(nodes) > 0 && !wuo.mutation.ChunksCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   worker.ChunksTable,
			Columns: []string{worker.ChunksColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(chunk.FieldID, field.TypeString),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := wuo.mutation.ChunksIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   worker.ChunksTable,
			Columns: []string{worker.ChunksColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(chunk.FieldID, field.TypeString),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	_node = &Worker{config: wuo.config}
	_spec.Assign = _node.assignValues
	_spec.ScanValues = _node.scanValues
	if err = sqlgraph.UpdateNode(ctx, wuo.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{worker.Label}
		} else if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{msg: err.Error(), wrap: err}
		}
		return nil, err
	}
	wuo.mutation.done = true
	return _node, nil
}
