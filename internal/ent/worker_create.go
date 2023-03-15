// Code generated by ent, DO NOT EDIT.

package ent

import (
	"context"
	"errors"
	"fmt"
	"time"

	"entgo.io/ent/dialect"
	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"entgo.io/ent/schema/field"
	"github.com/go-faster/gha/internal/ent/chunk"
	"github.com/go-faster/gha/internal/ent/worker"
	"github.com/google/uuid"
)

// WorkerCreate is the builder for creating a Worker entity.
type WorkerCreate struct {
	config
	mutation *WorkerMutation
	hooks    []Hook
	conflict []sql.ConflictOption
}

// SetCreatedAt sets the "created_at" field.
func (wc *WorkerCreate) SetCreatedAt(t time.Time) *WorkerCreate {
	wc.mutation.SetCreatedAt(t)
	return wc
}

// SetNillableCreatedAt sets the "created_at" field if the given value is not nil.
func (wc *WorkerCreate) SetNillableCreatedAt(t *time.Time) *WorkerCreate {
	if t != nil {
		wc.SetCreatedAt(*t)
	}
	return wc
}

// SetUpdatedAt sets the "updated_at" field.
func (wc *WorkerCreate) SetUpdatedAt(t time.Time) *WorkerCreate {
	wc.mutation.SetUpdatedAt(t)
	return wc
}

// SetNillableUpdatedAt sets the "updated_at" field if the given value is not nil.
func (wc *WorkerCreate) SetNillableUpdatedAt(t *time.Time) *WorkerCreate {
	if t != nil {
		wc.SetUpdatedAt(*t)
	}
	return wc
}

// SetName sets the "name" field.
func (wc *WorkerCreate) SetName(s string) *WorkerCreate {
	wc.mutation.SetName(s)
	return wc
}

// SetToken sets the "token" field.
func (wc *WorkerCreate) SetToken(s string) *WorkerCreate {
	wc.mutation.SetToken(s)
	return wc
}

// SetID sets the "id" field.
func (wc *WorkerCreate) SetID(u uuid.UUID) *WorkerCreate {
	wc.mutation.SetID(u)
	return wc
}

// SetNillableID sets the "id" field if the given value is not nil.
func (wc *WorkerCreate) SetNillableID(u *uuid.UUID) *WorkerCreate {
	if u != nil {
		wc.SetID(*u)
	}
	return wc
}

// AddChunkIDs adds the "chunks" edge to the Chunk entity by IDs.
func (wc *WorkerCreate) AddChunkIDs(ids ...string) *WorkerCreate {
	wc.mutation.AddChunkIDs(ids...)
	return wc
}

// AddChunks adds the "chunks" edges to the Chunk entity.
func (wc *WorkerCreate) AddChunks(c ...*Chunk) *WorkerCreate {
	ids := make([]string, len(c))
	for i := range c {
		ids[i] = c[i].ID
	}
	return wc.AddChunkIDs(ids...)
}

// Mutation returns the WorkerMutation object of the builder.
func (wc *WorkerCreate) Mutation() *WorkerMutation {
	return wc.mutation
}

// Save creates the Worker in the database.
func (wc *WorkerCreate) Save(ctx context.Context) (*Worker, error) {
	wc.defaults()
	return withHooks[*Worker, WorkerMutation](ctx, wc.sqlSave, wc.mutation, wc.hooks)
}

// SaveX calls Save and panics if Save returns an error.
func (wc *WorkerCreate) SaveX(ctx context.Context) *Worker {
	v, err := wc.Save(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Exec executes the query.
func (wc *WorkerCreate) Exec(ctx context.Context) error {
	_, err := wc.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (wc *WorkerCreate) ExecX(ctx context.Context) {
	if err := wc.Exec(ctx); err != nil {
		panic(err)
	}
}

// defaults sets the default values of the builder before save.
func (wc *WorkerCreate) defaults() {
	if _, ok := wc.mutation.CreatedAt(); !ok {
		v := worker.DefaultCreatedAt()
		wc.mutation.SetCreatedAt(v)
	}
	if _, ok := wc.mutation.UpdatedAt(); !ok {
		v := worker.DefaultUpdatedAt()
		wc.mutation.SetUpdatedAt(v)
	}
	if _, ok := wc.mutation.ID(); !ok {
		v := worker.DefaultID()
		wc.mutation.SetID(v)
	}
}

// check runs all checks and user-defined validators on the builder.
func (wc *WorkerCreate) check() error {
	if _, ok := wc.mutation.CreatedAt(); !ok {
		return &ValidationError{Name: "created_at", err: errors.New(`ent: missing required field "Worker.created_at"`)}
	}
	if _, ok := wc.mutation.UpdatedAt(); !ok {
		return &ValidationError{Name: "updated_at", err: errors.New(`ent: missing required field "Worker.updated_at"`)}
	}
	if _, ok := wc.mutation.Name(); !ok {
		return &ValidationError{Name: "name", err: errors.New(`ent: missing required field "Worker.name"`)}
	}
	if _, ok := wc.mutation.Token(); !ok {
		return &ValidationError{Name: "token", err: errors.New(`ent: missing required field "Worker.token"`)}
	}
	return nil
}

func (wc *WorkerCreate) sqlSave(ctx context.Context) (*Worker, error) {
	if err := wc.check(); err != nil {
		return nil, err
	}
	_node, _spec := wc.createSpec()
	if err := sqlgraph.CreateNode(ctx, wc.driver, _spec); err != nil {
		if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{msg: err.Error(), wrap: err}
		}
		return nil, err
	}
	if _spec.ID.Value != nil {
		if id, ok := _spec.ID.Value.(*uuid.UUID); ok {
			_node.ID = *id
		} else if err := _node.ID.Scan(_spec.ID.Value); err != nil {
			return nil, err
		}
	}
	wc.mutation.id = &_node.ID
	wc.mutation.done = true
	return _node, nil
}

func (wc *WorkerCreate) createSpec() (*Worker, *sqlgraph.CreateSpec) {
	var (
		_node = &Worker{config: wc.config}
		_spec = sqlgraph.NewCreateSpec(worker.Table, sqlgraph.NewFieldSpec(worker.FieldID, field.TypeUUID))
	)
	_spec.OnConflict = wc.conflict
	if id, ok := wc.mutation.ID(); ok {
		_node.ID = id
		_spec.ID.Value = &id
	}
	if value, ok := wc.mutation.CreatedAt(); ok {
		_spec.SetField(worker.FieldCreatedAt, field.TypeTime, value)
		_node.CreatedAt = value
	}
	if value, ok := wc.mutation.UpdatedAt(); ok {
		_spec.SetField(worker.FieldUpdatedAt, field.TypeTime, value)
		_node.UpdatedAt = value
	}
	if value, ok := wc.mutation.Name(); ok {
		_spec.SetField(worker.FieldName, field.TypeString, value)
		_node.Name = value
	}
	if value, ok := wc.mutation.Token(); ok {
		_spec.SetField(worker.FieldToken, field.TypeString, value)
		_node.Token = value
	}
	if nodes := wc.mutation.ChunksIDs(); len(nodes) > 0 {
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
		_spec.Edges = append(_spec.Edges, edge)
	}
	return _node, _spec
}

// OnConflict allows configuring the `ON CONFLICT` / `ON DUPLICATE KEY` clause
// of the `INSERT` statement. For example:
//
//	client.Worker.Create().
//		SetCreatedAt(v).
//		OnConflict(
//			// Update the row with the new values
//			// the was proposed for insertion.
//			sql.ResolveWithNewValues(),
//		).
//		// Override some of the fields with custom
//		// update values.
//		Update(func(u *ent.WorkerUpsert) {
//			SetCreatedAt(v+v).
//		}).
//		Exec(ctx)
func (wc *WorkerCreate) OnConflict(opts ...sql.ConflictOption) *WorkerUpsertOne {
	wc.conflict = opts
	return &WorkerUpsertOne{
		create: wc,
	}
}

// OnConflictColumns calls `OnConflict` and configures the columns
// as conflict target. Using this option is equivalent to using:
//
//	client.Worker.Create().
//		OnConflict(sql.ConflictColumns(columns...)).
//		Exec(ctx)
func (wc *WorkerCreate) OnConflictColumns(columns ...string) *WorkerUpsertOne {
	wc.conflict = append(wc.conflict, sql.ConflictColumns(columns...))
	return &WorkerUpsertOne{
		create: wc,
	}
}

type (
	// WorkerUpsertOne is the builder for "upsert"-ing
	//  one Worker node.
	WorkerUpsertOne struct {
		create *WorkerCreate
	}

	// WorkerUpsert is the "OnConflict" setter.
	WorkerUpsert struct {
		*sql.UpdateSet
	}
)

// SetUpdatedAt sets the "updated_at" field.
func (u *WorkerUpsert) SetUpdatedAt(v time.Time) *WorkerUpsert {
	u.Set(worker.FieldUpdatedAt, v)
	return u
}

// UpdateUpdatedAt sets the "updated_at" field to the value that was provided on create.
func (u *WorkerUpsert) UpdateUpdatedAt() *WorkerUpsert {
	u.SetExcluded(worker.FieldUpdatedAt)
	return u
}

// SetName sets the "name" field.
func (u *WorkerUpsert) SetName(v string) *WorkerUpsert {
	u.Set(worker.FieldName, v)
	return u
}

// UpdateName sets the "name" field to the value that was provided on create.
func (u *WorkerUpsert) UpdateName() *WorkerUpsert {
	u.SetExcluded(worker.FieldName)
	return u
}

// SetToken sets the "token" field.
func (u *WorkerUpsert) SetToken(v string) *WorkerUpsert {
	u.Set(worker.FieldToken, v)
	return u
}

// UpdateToken sets the "token" field to the value that was provided on create.
func (u *WorkerUpsert) UpdateToken() *WorkerUpsert {
	u.SetExcluded(worker.FieldToken)
	return u
}

// UpdateNewValues updates the mutable fields using the new values that were set on create except the ID field.
// Using this option is equivalent to using:
//
//	client.Worker.Create().
//		OnConflict(
//			sql.ResolveWithNewValues(),
//			sql.ResolveWith(func(u *sql.UpdateSet) {
//				u.SetIgnore(worker.FieldID)
//			}),
//		).
//		Exec(ctx)
func (u *WorkerUpsertOne) UpdateNewValues() *WorkerUpsertOne {
	u.create.conflict = append(u.create.conflict, sql.ResolveWithNewValues())
	u.create.conflict = append(u.create.conflict, sql.ResolveWith(func(s *sql.UpdateSet) {
		if _, exists := u.create.mutation.ID(); exists {
			s.SetIgnore(worker.FieldID)
		}
		if _, exists := u.create.mutation.CreatedAt(); exists {
			s.SetIgnore(worker.FieldCreatedAt)
		}
	}))
	return u
}

// Ignore sets each column to itself in case of conflict.
// Using this option is equivalent to using:
//
//	client.Worker.Create().
//	    OnConflict(sql.ResolveWithIgnore()).
//	    Exec(ctx)
func (u *WorkerUpsertOne) Ignore() *WorkerUpsertOne {
	u.create.conflict = append(u.create.conflict, sql.ResolveWithIgnore())
	return u
}

// DoNothing configures the conflict_action to `DO NOTHING`.
// Supported only by SQLite and PostgreSQL.
func (u *WorkerUpsertOne) DoNothing() *WorkerUpsertOne {
	u.create.conflict = append(u.create.conflict, sql.DoNothing())
	return u
}

// Update allows overriding fields `UPDATE` values. See the WorkerCreate.OnConflict
// documentation for more info.
func (u *WorkerUpsertOne) Update(set func(*WorkerUpsert)) *WorkerUpsertOne {
	u.create.conflict = append(u.create.conflict, sql.ResolveWith(func(update *sql.UpdateSet) {
		set(&WorkerUpsert{UpdateSet: update})
	}))
	return u
}

// SetUpdatedAt sets the "updated_at" field.
func (u *WorkerUpsertOne) SetUpdatedAt(v time.Time) *WorkerUpsertOne {
	return u.Update(func(s *WorkerUpsert) {
		s.SetUpdatedAt(v)
	})
}

// UpdateUpdatedAt sets the "updated_at" field to the value that was provided on create.
func (u *WorkerUpsertOne) UpdateUpdatedAt() *WorkerUpsertOne {
	return u.Update(func(s *WorkerUpsert) {
		s.UpdateUpdatedAt()
	})
}

// SetName sets the "name" field.
func (u *WorkerUpsertOne) SetName(v string) *WorkerUpsertOne {
	return u.Update(func(s *WorkerUpsert) {
		s.SetName(v)
	})
}

// UpdateName sets the "name" field to the value that was provided on create.
func (u *WorkerUpsertOne) UpdateName() *WorkerUpsertOne {
	return u.Update(func(s *WorkerUpsert) {
		s.UpdateName()
	})
}

// SetToken sets the "token" field.
func (u *WorkerUpsertOne) SetToken(v string) *WorkerUpsertOne {
	return u.Update(func(s *WorkerUpsert) {
		s.SetToken(v)
	})
}

// UpdateToken sets the "token" field to the value that was provided on create.
func (u *WorkerUpsertOne) UpdateToken() *WorkerUpsertOne {
	return u.Update(func(s *WorkerUpsert) {
		s.UpdateToken()
	})
}

// Exec executes the query.
func (u *WorkerUpsertOne) Exec(ctx context.Context) error {
	if len(u.create.conflict) == 0 {
		return errors.New("ent: missing options for WorkerCreate.OnConflict")
	}
	return u.create.Exec(ctx)
}

// ExecX is like Exec, but panics if an error occurs.
func (u *WorkerUpsertOne) ExecX(ctx context.Context) {
	if err := u.create.Exec(ctx); err != nil {
		panic(err)
	}
}

// Exec executes the UPSERT query and returns the inserted/updated ID.
func (u *WorkerUpsertOne) ID(ctx context.Context) (id uuid.UUID, err error) {
	if u.create.driver.Dialect() == dialect.MySQL {
		// In case of "ON CONFLICT", there is no way to get back non-numeric ID
		// fields from the database since MySQL does not support the RETURNING clause.
		return id, errors.New("ent: WorkerUpsertOne.ID is not supported by MySQL driver. Use WorkerUpsertOne.Exec instead")
	}
	node, err := u.create.Save(ctx)
	if err != nil {
		return id, err
	}
	return node.ID, nil
}

// IDX is like ID, but panics if an error occurs.
func (u *WorkerUpsertOne) IDX(ctx context.Context) uuid.UUID {
	id, err := u.ID(ctx)
	if err != nil {
		panic(err)
	}
	return id
}

// WorkerCreateBulk is the builder for creating many Worker entities in bulk.
type WorkerCreateBulk struct {
	config
	builders []*WorkerCreate
	conflict []sql.ConflictOption
}

// Save creates the Worker entities in the database.
func (wcb *WorkerCreateBulk) Save(ctx context.Context) ([]*Worker, error) {
	specs := make([]*sqlgraph.CreateSpec, len(wcb.builders))
	nodes := make([]*Worker, len(wcb.builders))
	mutators := make([]Mutator, len(wcb.builders))
	for i := range wcb.builders {
		func(i int, root context.Context) {
			builder := wcb.builders[i]
			builder.defaults()
			var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
				mutation, ok := m.(*WorkerMutation)
				if !ok {
					return nil, fmt.Errorf("unexpected mutation type %T", m)
				}
				if err := builder.check(); err != nil {
					return nil, err
				}
				builder.mutation = mutation
				nodes[i], specs[i] = builder.createSpec()
				var err error
				if i < len(mutators)-1 {
					_, err = mutators[i+1].Mutate(root, wcb.builders[i+1].mutation)
				} else {
					spec := &sqlgraph.BatchCreateSpec{Nodes: specs}
					spec.OnConflict = wcb.conflict
					// Invoke the actual operation on the latest mutation in the chain.
					if err = sqlgraph.BatchCreate(ctx, wcb.driver, spec); err != nil {
						if sqlgraph.IsConstraintError(err) {
							err = &ConstraintError{msg: err.Error(), wrap: err}
						}
					}
				}
				if err != nil {
					return nil, err
				}
				mutation.id = &nodes[i].ID
				mutation.done = true
				return nodes[i], nil
			})
			for i := len(builder.hooks) - 1; i >= 0; i-- {
				mut = builder.hooks[i](mut)
			}
			mutators[i] = mut
		}(i, ctx)
	}
	if len(mutators) > 0 {
		if _, err := mutators[0].Mutate(ctx, wcb.builders[0].mutation); err != nil {
			return nil, err
		}
	}
	return nodes, nil
}

// SaveX is like Save, but panics if an error occurs.
func (wcb *WorkerCreateBulk) SaveX(ctx context.Context) []*Worker {
	v, err := wcb.Save(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Exec executes the query.
func (wcb *WorkerCreateBulk) Exec(ctx context.Context) error {
	_, err := wcb.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (wcb *WorkerCreateBulk) ExecX(ctx context.Context) {
	if err := wcb.Exec(ctx); err != nil {
		panic(err)
	}
}

// OnConflict allows configuring the `ON CONFLICT` / `ON DUPLICATE KEY` clause
// of the `INSERT` statement. For example:
//
//	client.Worker.CreateBulk(builders...).
//		OnConflict(
//			// Update the row with the new values
//			// the was proposed for insertion.
//			sql.ResolveWithNewValues(),
//		).
//		// Override some of the fields with custom
//		// update values.
//		Update(func(u *ent.WorkerUpsert) {
//			SetCreatedAt(v+v).
//		}).
//		Exec(ctx)
func (wcb *WorkerCreateBulk) OnConflict(opts ...sql.ConflictOption) *WorkerUpsertBulk {
	wcb.conflict = opts
	return &WorkerUpsertBulk{
		create: wcb,
	}
}

// OnConflictColumns calls `OnConflict` and configures the columns
// as conflict target. Using this option is equivalent to using:
//
//	client.Worker.Create().
//		OnConflict(sql.ConflictColumns(columns...)).
//		Exec(ctx)
func (wcb *WorkerCreateBulk) OnConflictColumns(columns ...string) *WorkerUpsertBulk {
	wcb.conflict = append(wcb.conflict, sql.ConflictColumns(columns...))
	return &WorkerUpsertBulk{
		create: wcb,
	}
}

// WorkerUpsertBulk is the builder for "upsert"-ing
// a bulk of Worker nodes.
type WorkerUpsertBulk struct {
	create *WorkerCreateBulk
}

// UpdateNewValues updates the mutable fields using the new values that
// were set on create. Using this option is equivalent to using:
//
//	client.Worker.Create().
//		OnConflict(
//			sql.ResolveWithNewValues(),
//			sql.ResolveWith(func(u *sql.UpdateSet) {
//				u.SetIgnore(worker.FieldID)
//			}),
//		).
//		Exec(ctx)
func (u *WorkerUpsertBulk) UpdateNewValues() *WorkerUpsertBulk {
	u.create.conflict = append(u.create.conflict, sql.ResolveWithNewValues())
	u.create.conflict = append(u.create.conflict, sql.ResolveWith(func(s *sql.UpdateSet) {
		for _, b := range u.create.builders {
			if _, exists := b.mutation.ID(); exists {
				s.SetIgnore(worker.FieldID)
			}
			if _, exists := b.mutation.CreatedAt(); exists {
				s.SetIgnore(worker.FieldCreatedAt)
			}
		}
	}))
	return u
}

// Ignore sets each column to itself in case of conflict.
// Using this option is equivalent to using:
//
//	client.Worker.Create().
//		OnConflict(sql.ResolveWithIgnore()).
//		Exec(ctx)
func (u *WorkerUpsertBulk) Ignore() *WorkerUpsertBulk {
	u.create.conflict = append(u.create.conflict, sql.ResolveWithIgnore())
	return u
}

// DoNothing configures the conflict_action to `DO NOTHING`.
// Supported only by SQLite and PostgreSQL.
func (u *WorkerUpsertBulk) DoNothing() *WorkerUpsertBulk {
	u.create.conflict = append(u.create.conflict, sql.DoNothing())
	return u
}

// Update allows overriding fields `UPDATE` values. See the WorkerCreateBulk.OnConflict
// documentation for more info.
func (u *WorkerUpsertBulk) Update(set func(*WorkerUpsert)) *WorkerUpsertBulk {
	u.create.conflict = append(u.create.conflict, sql.ResolveWith(func(update *sql.UpdateSet) {
		set(&WorkerUpsert{UpdateSet: update})
	}))
	return u
}

// SetUpdatedAt sets the "updated_at" field.
func (u *WorkerUpsertBulk) SetUpdatedAt(v time.Time) *WorkerUpsertBulk {
	return u.Update(func(s *WorkerUpsert) {
		s.SetUpdatedAt(v)
	})
}

// UpdateUpdatedAt sets the "updated_at" field to the value that was provided on create.
func (u *WorkerUpsertBulk) UpdateUpdatedAt() *WorkerUpsertBulk {
	return u.Update(func(s *WorkerUpsert) {
		s.UpdateUpdatedAt()
	})
}

// SetName sets the "name" field.
func (u *WorkerUpsertBulk) SetName(v string) *WorkerUpsertBulk {
	return u.Update(func(s *WorkerUpsert) {
		s.SetName(v)
	})
}

// UpdateName sets the "name" field to the value that was provided on create.
func (u *WorkerUpsertBulk) UpdateName() *WorkerUpsertBulk {
	return u.Update(func(s *WorkerUpsert) {
		s.UpdateName()
	})
}

// SetToken sets the "token" field.
func (u *WorkerUpsertBulk) SetToken(v string) *WorkerUpsertBulk {
	return u.Update(func(s *WorkerUpsert) {
		s.SetToken(v)
	})
}

// UpdateToken sets the "token" field to the value that was provided on create.
func (u *WorkerUpsertBulk) UpdateToken() *WorkerUpsertBulk {
	return u.Update(func(s *WorkerUpsert) {
		s.UpdateToken()
	})
}

// Exec executes the query.
func (u *WorkerUpsertBulk) Exec(ctx context.Context) error {
	for i, b := range u.create.builders {
		if len(b.conflict) != 0 {
			return fmt.Errorf("ent: OnConflict was set for builder %d. Set it on the WorkerCreateBulk instead", i)
		}
	}
	if len(u.create.conflict) == 0 {
		return errors.New("ent: missing options for WorkerCreateBulk.OnConflict")
	}
	return u.create.Exec(ctx)
}

// ExecX is like Exec, but panics if an error occurs.
func (u *WorkerUpsertBulk) ExecX(ctx context.Context) {
	if err := u.create.Exec(ctx); err != nil {
		panic(err)
	}
}
