// Code generated by entc, DO NOT EDIT.

package ent

import (
	"fmt"
	"strings"
	"time"

	"entgo.io/ent/dialect/sql"
	"github.com/go-faster/gha/internal/ent/chunk"
	"github.com/go-faster/gha/internal/ent/worker"
	"github.com/google/uuid"
)

// Chunk is the model entity for the Chunk schema.
type Chunk struct {
	config `json:"-"`
	// ID of the ent.
	// Like 2006-01-02T15
	ID string `json:"id,omitempty"`
	// CreatedAt holds the value of the "created_at" field.
	// Time when entity was created.
	CreatedAt time.Time `json:"created_at,omitempty"`
	// UpdatedAt holds the value of the "updated_at" field.
	// Time when entity was updated.
	UpdatedAt time.Time `json:"updated_at,omitempty"`
	// Start holds the value of the "start" field.
	// Minimum possible time of entry in chunk
	Start time.Time `json:"start,omitempty"`
	// LeaseExpiresAt holds the value of the "lease_expires_at" field.
	// State expiration like heartbeat
	LeaseExpiresAt time.Time `json:"lease_expires_at,omitempty"`
	// State holds the value of the "state" field.
	State chunk.State `json:"state,omitempty"`
	// SizeInput holds the value of the "size_input" field.
	SizeInput int64 `json:"size_input,omitempty"`
	// SizeContent holds the value of the "size_content" field.
	SizeContent int64 `json:"size_content,omitempty"`
	// SizeOutput holds the value of the "size_output" field.
	SizeOutput int64 `json:"size_output,omitempty"`
	// Sha256Input holds the value of the "sha256_input" field.
	Sha256Input *string `json:"sha256_input,omitempty"`
	// Sha256Content holds the value of the "sha256_content" field.
	Sha256Content *string `json:"sha256_content,omitempty"`
	// Sha256Output holds the value of the "sha256_output" field.
	Sha256Output *string `json:"sha256_output,omitempty"`
	// Edges holds the relations/edges for other nodes in the graph.
	// The values are being populated by the ChunkQuery when eager-loading is set.
	Edges         ChunkEdges `json:"edges"`
	worker_chunks *uuid.UUID
}

// ChunkEdges holds the relations/edges for other nodes in the graph.
type ChunkEdges struct {
	// Worker holds the value of the worker edge.
	Worker *Worker `json:"worker,omitempty"`
	// loadedTypes holds the information for reporting if a
	// type was loaded (or requested) in eager-loading or not.
	loadedTypes [1]bool
}

// WorkerOrErr returns the Worker value or an error if the edge
// was not loaded in eager-loading, or loaded but was not found.
func (e ChunkEdges) WorkerOrErr() (*Worker, error) {
	if e.loadedTypes[0] {
		if e.Worker == nil {
			// The edge worker was loaded in eager-loading,
			// but was not found.
			return nil, &NotFoundError{label: worker.Label}
		}
		return e.Worker, nil
	}
	return nil, &NotLoadedError{edge: "worker"}
}

// scanValues returns the types for scanning values from sql.Rows.
func (*Chunk) scanValues(columns []string) ([]interface{}, error) {
	values := make([]interface{}, len(columns))
	for i := range columns {
		switch columns[i] {
		case chunk.FieldSizeInput, chunk.FieldSizeContent, chunk.FieldSizeOutput:
			values[i] = new(sql.NullInt64)
		case chunk.FieldID, chunk.FieldState, chunk.FieldSha256Input, chunk.FieldSha256Content, chunk.FieldSha256Output:
			values[i] = new(sql.NullString)
		case chunk.FieldCreatedAt, chunk.FieldUpdatedAt, chunk.FieldStart, chunk.FieldLeaseExpiresAt:
			values[i] = new(sql.NullTime)
		case chunk.ForeignKeys[0]: // worker_chunks
			values[i] = &sql.NullScanner{S: new(uuid.UUID)}
		default:
			return nil, fmt.Errorf("unexpected column %q for type Chunk", columns[i])
		}
	}
	return values, nil
}

// assignValues assigns the values that were returned from sql.Rows (after scanning)
// to the Chunk fields.
func (c *Chunk) assignValues(columns []string, values []interface{}) error {
	if m, n := len(values), len(columns); m < n {
		return fmt.Errorf("mismatch number of scan values: %d != %d", m, n)
	}
	for i := range columns {
		switch columns[i] {
		case chunk.FieldID:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field id", values[i])
			} else if value.Valid {
				c.ID = value.String
			}
		case chunk.FieldCreatedAt:
			if value, ok := values[i].(*sql.NullTime); !ok {
				return fmt.Errorf("unexpected type %T for field created_at", values[i])
			} else if value.Valid {
				c.CreatedAt = value.Time
			}
		case chunk.FieldUpdatedAt:
			if value, ok := values[i].(*sql.NullTime); !ok {
				return fmt.Errorf("unexpected type %T for field updated_at", values[i])
			} else if value.Valid {
				c.UpdatedAt = value.Time
			}
		case chunk.FieldStart:
			if value, ok := values[i].(*sql.NullTime); !ok {
				return fmt.Errorf("unexpected type %T for field start", values[i])
			} else if value.Valid {
				c.Start = value.Time
			}
		case chunk.FieldLeaseExpiresAt:
			if value, ok := values[i].(*sql.NullTime); !ok {
				return fmt.Errorf("unexpected type %T for field lease_expires_at", values[i])
			} else if value.Valid {
				c.LeaseExpiresAt = value.Time
			}
		case chunk.FieldState:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field state", values[i])
			} else if value.Valid {
				c.State = chunk.State(value.String)
			}
		case chunk.FieldSizeInput:
			if value, ok := values[i].(*sql.NullInt64); !ok {
				return fmt.Errorf("unexpected type %T for field size_input", values[i])
			} else if value.Valid {
				c.SizeInput = value.Int64
			}
		case chunk.FieldSizeContent:
			if value, ok := values[i].(*sql.NullInt64); !ok {
				return fmt.Errorf("unexpected type %T for field size_content", values[i])
			} else if value.Valid {
				c.SizeContent = value.Int64
			}
		case chunk.FieldSizeOutput:
			if value, ok := values[i].(*sql.NullInt64); !ok {
				return fmt.Errorf("unexpected type %T for field size_output", values[i])
			} else if value.Valid {
				c.SizeOutput = value.Int64
			}
		case chunk.FieldSha256Input:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field sha256_input", values[i])
			} else if value.Valid {
				c.Sha256Input = new(string)
				*c.Sha256Input = value.String
			}
		case chunk.FieldSha256Content:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field sha256_content", values[i])
			} else if value.Valid {
				c.Sha256Content = new(string)
				*c.Sha256Content = value.String
			}
		case chunk.FieldSha256Output:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field sha256_output", values[i])
			} else if value.Valid {
				c.Sha256Output = new(string)
				*c.Sha256Output = value.String
			}
		case chunk.ForeignKeys[0]:
			if value, ok := values[i].(*sql.NullScanner); !ok {
				return fmt.Errorf("unexpected type %T for field worker_chunks", values[i])
			} else if value.Valid {
				c.worker_chunks = new(uuid.UUID)
				*c.worker_chunks = *value.S.(*uuid.UUID)
			}
		}
	}
	return nil
}

// QueryWorker queries the "worker" edge of the Chunk entity.
func (c *Chunk) QueryWorker() *WorkerQuery {
	return (&ChunkClient{config: c.config}).QueryWorker(c)
}

// Update returns a builder for updating this Chunk.
// Note that you need to call Chunk.Unwrap() before calling this method if this Chunk
// was returned from a transaction, and the transaction was committed or rolled back.
func (c *Chunk) Update() *ChunkUpdateOne {
	return (&ChunkClient{config: c.config}).UpdateOne(c)
}

// Unwrap unwraps the Chunk entity that was returned from a transaction after it was closed,
// so that all future queries will be executed through the driver which created the transaction.
func (c *Chunk) Unwrap() *Chunk {
	tx, ok := c.config.driver.(*txDriver)
	if !ok {
		panic("ent: Chunk is not a transactional entity")
	}
	c.config.driver = tx.drv
	return c
}

// String implements the fmt.Stringer.
func (c *Chunk) String() string {
	var builder strings.Builder
	builder.WriteString("Chunk(")
	builder.WriteString(fmt.Sprintf("id=%v", c.ID))
	builder.WriteString(", created_at=")
	builder.WriteString(c.CreatedAt.Format(time.ANSIC))
	builder.WriteString(", updated_at=")
	builder.WriteString(c.UpdatedAt.Format(time.ANSIC))
	builder.WriteString(", start=")
	builder.WriteString(c.Start.Format(time.ANSIC))
	builder.WriteString(", lease_expires_at=")
	builder.WriteString(c.LeaseExpiresAt.Format(time.ANSIC))
	builder.WriteString(", state=")
	builder.WriteString(fmt.Sprintf("%v", c.State))
	builder.WriteString(", size_input=")
	builder.WriteString(fmt.Sprintf("%v", c.SizeInput))
	builder.WriteString(", size_content=")
	builder.WriteString(fmt.Sprintf("%v", c.SizeContent))
	builder.WriteString(", size_output=")
	builder.WriteString(fmt.Sprintf("%v", c.SizeOutput))
	if v := c.Sha256Input; v != nil {
		builder.WriteString(", sha256_input=")
		builder.WriteString(*v)
	}
	if v := c.Sha256Content; v != nil {
		builder.WriteString(", sha256_content=")
		builder.WriteString(*v)
	}
	if v := c.Sha256Output; v != nil {
		builder.WriteString(", sha256_output=")
		builder.WriteString(*v)
	}
	builder.WriteByte(')')
	return builder.String()
}

// Chunks is a parsable slice of Chunk.
type Chunks []*Chunk

func (c Chunks) config(cfg config) {
	for _i := range c {
		c[_i].config = cfg
	}
}
