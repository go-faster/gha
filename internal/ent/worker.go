// Code generated by ent, DO NOT EDIT.

package ent

import (
	"fmt"
	"strings"
	"time"

	"entgo.io/ent/dialect/sql"
	"github.com/go-faster/gha/internal/ent/worker"
	"github.com/google/uuid"
)

// Worker is the model entity for the Worker schema.
type Worker struct {
	config `json:"-"`
	// ID of the ent.
	ID uuid.UUID `json:"id,omitempty"`
	// Time when entity was created.
	CreatedAt time.Time `json:"created_at,omitempty"`
	// Time when entity was updated.
	UpdatedAt time.Time `json:"updated_at,omitempty"`
	// Name holds the value of the "name" field.
	Name string `json:"name,omitempty"`
	// Token holds the value of the "token" field.
	Token string `json:"token,omitempty"`
	// Edges holds the relations/edges for other nodes in the graph.
	// The values are being populated by the WorkerQuery when eager-loading is set.
	Edges WorkerEdges `json:"edges"`
}

// WorkerEdges holds the relations/edges for other nodes in the graph.
type WorkerEdges struct {
	// Chunks holds the value of the chunks edge.
	Chunks []*Chunk `json:"chunks,omitempty"`
	// loadedTypes holds the information for reporting if a
	// type was loaded (or requested) in eager-loading or not.
	loadedTypes [1]bool
}

// ChunksOrErr returns the Chunks value or an error if the edge
// was not loaded in eager-loading.
func (e WorkerEdges) ChunksOrErr() ([]*Chunk, error) {
	if e.loadedTypes[0] {
		return e.Chunks, nil
	}
	return nil, &NotLoadedError{edge: "chunks"}
}

// scanValues returns the types for scanning values from sql.Rows.
func (*Worker) scanValues(columns []string) ([]interface{}, error) {
	values := make([]interface{}, len(columns))
	for i := range columns {
		switch columns[i] {
		case worker.FieldName, worker.FieldToken:
			values[i] = new(sql.NullString)
		case worker.FieldCreatedAt, worker.FieldUpdatedAt:
			values[i] = new(sql.NullTime)
		case worker.FieldID:
			values[i] = new(uuid.UUID)
		default:
			return nil, fmt.Errorf("unexpected column %q for type Worker", columns[i])
		}
	}
	return values, nil
}

// assignValues assigns the values that were returned from sql.Rows (after scanning)
// to the Worker fields.
func (w *Worker) assignValues(columns []string, values []interface{}) error {
	if m, n := len(values), len(columns); m < n {
		return fmt.Errorf("mismatch number of scan values: %d != %d", m, n)
	}
	for i := range columns {
		switch columns[i] {
		case worker.FieldID:
			if value, ok := values[i].(*uuid.UUID); !ok {
				return fmt.Errorf("unexpected type %T for field id", values[i])
			} else if value != nil {
				w.ID = *value
			}
		case worker.FieldCreatedAt:
			if value, ok := values[i].(*sql.NullTime); !ok {
				return fmt.Errorf("unexpected type %T for field created_at", values[i])
			} else if value.Valid {
				w.CreatedAt = value.Time
			}
		case worker.FieldUpdatedAt:
			if value, ok := values[i].(*sql.NullTime); !ok {
				return fmt.Errorf("unexpected type %T for field updated_at", values[i])
			} else if value.Valid {
				w.UpdatedAt = value.Time
			}
		case worker.FieldName:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field name", values[i])
			} else if value.Valid {
				w.Name = value.String
			}
		case worker.FieldToken:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field token", values[i])
			} else if value.Valid {
				w.Token = value.String
			}
		}
	}
	return nil
}

// QueryChunks queries the "chunks" edge of the Worker entity.
func (w *Worker) QueryChunks() *ChunkQuery {
	return (&WorkerClient{config: w.config}).QueryChunks(w)
}

// Update returns a builder for updating this Worker.
// Note that you need to call Worker.Unwrap() before calling this method if this Worker
// was returned from a transaction, and the transaction was committed or rolled back.
func (w *Worker) Update() *WorkerUpdateOne {
	return (&WorkerClient{config: w.config}).UpdateOne(w)
}

// Unwrap unwraps the Worker entity that was returned from a transaction after it was closed,
// so that all future queries will be executed through the driver which created the transaction.
func (w *Worker) Unwrap() *Worker {
	_tx, ok := w.config.driver.(*txDriver)
	if !ok {
		panic("ent: Worker is not a transactional entity")
	}
	w.config.driver = _tx.drv
	return w
}

// String implements the fmt.Stringer.
func (w *Worker) String() string {
	var builder strings.Builder
	builder.WriteString("Worker(")
	builder.WriteString(fmt.Sprintf("id=%v, ", w.ID))
	builder.WriteString("created_at=")
	builder.WriteString(w.CreatedAt.Format(time.ANSIC))
	builder.WriteString(", ")
	builder.WriteString("updated_at=")
	builder.WriteString(w.UpdatedAt.Format(time.ANSIC))
	builder.WriteString(", ")
	builder.WriteString("name=")
	builder.WriteString(w.Name)
	builder.WriteString(", ")
	builder.WriteString("token=")
	builder.WriteString(w.Token)
	builder.WriteByte(')')
	return builder.String()
}

// Workers is a parsable slice of Worker.
type Workers []*Worker

func (w Workers) config(cfg config) {
	for _i := range w {
		w[_i].config = cfg
	}
}
