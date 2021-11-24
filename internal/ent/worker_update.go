// Code generated by entc, DO NOT EDIT.

package ent

import (
	"context"
	"fmt"
	"time"

	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"entgo.io/ent/schema/field"
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

// Mutation returns the WorkerMutation object of the builder.
func (wu *WorkerUpdate) Mutation() *WorkerMutation {
	return wu.mutation
}

// Save executes the query and returns the number of nodes affected by the update operation.
func (wu *WorkerUpdate) Save(ctx context.Context) (int, error) {
	var (
		err      error
		affected int
	)
	wu.defaults()
	if len(wu.hooks) == 0 {
		affected, err = wu.sqlSave(ctx)
	} else {
		var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
			mutation, ok := m.(*WorkerMutation)
			if !ok {
				return nil, fmt.Errorf("unexpected mutation type %T", m)
			}
			wu.mutation = mutation
			affected, err = wu.sqlSave(ctx)
			mutation.done = true
			return affected, err
		})
		for i := len(wu.hooks) - 1; i >= 0; i-- {
			if wu.hooks[i] == nil {
				return 0, fmt.Errorf("ent: uninitialized hook (forgotten import ent/runtime?)")
			}
			mut = wu.hooks[i](mut)
		}
		if _, err := mut.Mutate(ctx, wu.mutation); err != nil {
			return 0, err
		}
	}
	return affected, err
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
	_spec := &sqlgraph.UpdateSpec{
		Node: &sqlgraph.NodeSpec{
			Table:   worker.Table,
			Columns: worker.Columns,
			ID: &sqlgraph.FieldSpec{
				Type:   field.TypeUUID,
				Column: worker.FieldID,
			},
		},
	}
	if ps := wu.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if value, ok := wu.mutation.UpdatedAt(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeTime,
			Value:  value,
			Column: worker.FieldUpdatedAt,
		})
	}
	if value, ok := wu.mutation.Name(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: worker.FieldName,
		})
	}
	if value, ok := wu.mutation.Token(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: worker.FieldToken,
		})
	}
	if n, err = sqlgraph.UpdateNodes(ctx, wu.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{worker.Label}
		} else if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{err.Error(), err}
		}
		return 0, err
	}
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

// Mutation returns the WorkerMutation object of the builder.
func (wuo *WorkerUpdateOne) Mutation() *WorkerMutation {
	return wuo.mutation
}

// Select allows selecting one or more fields (columns) of the returned entity.
// The default is selecting all fields defined in the entity schema.
func (wuo *WorkerUpdateOne) Select(field string, fields ...string) *WorkerUpdateOne {
	wuo.fields = append([]string{field}, fields...)
	return wuo
}

// Save executes the query and returns the updated Worker entity.
func (wuo *WorkerUpdateOne) Save(ctx context.Context) (*Worker, error) {
	var (
		err  error
		node *Worker
	)
	wuo.defaults()
	if len(wuo.hooks) == 0 {
		node, err = wuo.sqlSave(ctx)
	} else {
		var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
			mutation, ok := m.(*WorkerMutation)
			if !ok {
				return nil, fmt.Errorf("unexpected mutation type %T", m)
			}
			wuo.mutation = mutation
			node, err = wuo.sqlSave(ctx)
			mutation.done = true
			return node, err
		})
		for i := len(wuo.hooks) - 1; i >= 0; i-- {
			if wuo.hooks[i] == nil {
				return nil, fmt.Errorf("ent: uninitialized hook (forgotten import ent/runtime?)")
			}
			mut = wuo.hooks[i](mut)
		}
		if _, err := mut.Mutate(ctx, wuo.mutation); err != nil {
			return nil, err
		}
	}
	return node, err
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
	_spec := &sqlgraph.UpdateSpec{
		Node: &sqlgraph.NodeSpec{
			Table:   worker.Table,
			Columns: worker.Columns,
			ID: &sqlgraph.FieldSpec{
				Type:   field.TypeUUID,
				Column: worker.FieldID,
			},
		},
	}
	id, ok := wuo.mutation.ID()
	if !ok {
		return nil, &ValidationError{Name: "ID", err: fmt.Errorf("missing Worker.ID for update")}
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
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeTime,
			Value:  value,
			Column: worker.FieldUpdatedAt,
		})
	}
	if value, ok := wuo.mutation.Name(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: worker.FieldName,
		})
	}
	if value, ok := wuo.mutation.Token(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: worker.FieldToken,
		})
	}
	_node = &Worker{config: wuo.config}
	_spec.Assign = _node.assignValues
	_spec.ScanValues = _node.scanValues
	if err = sqlgraph.UpdateNode(ctx, wuo.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{worker.Label}
		} else if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{err.Error(), err}
		}
		return nil, err
	}
	return _node, nil
}
