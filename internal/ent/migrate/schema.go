// Code generated by entc, DO NOT EDIT.

package migrate

import (
	"entgo.io/ent/dialect/sql/schema"
	"entgo.io/ent/schema/field"
)

var (
	// ChunksColumns holds the columns for the "chunks" table.
	ChunksColumns = []*schema.Column{
		{Name: "id", Type: field.TypeString, Unique: true},
		{Name: "created_at", Type: field.TypeTime},
		{Name: "updated_at", Type: field.TypeTime},
		{Name: "start", Type: field.TypeTime},
		{Name: "lease_expires_at", Type: field.TypeTime, Nullable: true},
		{Name: "state", Type: field.TypeEnum, Enums: []string{"New", "Downloading", "Downloaded", "Inventory", "Ready", "Processing", "Done"}, Default: "New"},
		{Name: "size_input", Type: field.TypeInt64, Nullable: true},
		{Name: "size_content", Type: field.TypeInt64, Nullable: true},
		{Name: "size_output", Type: field.TypeInt64, Nullable: true},
		{Name: "sha256_input", Type: field.TypeString, Nullable: true},
		{Name: "sha256_content", Type: field.TypeString, Nullable: true},
		{Name: "sha256_output", Type: field.TypeString, Nullable: true},
		{Name: "worker_chunks", Type: field.TypeUUID, Nullable: true},
	}
	// ChunksTable holds the schema information for the "chunks" table.
	ChunksTable = &schema.Table{
		Name:       "chunks",
		Columns:    ChunksColumns,
		PrimaryKey: []*schema.Column{ChunksColumns[0]},
		ForeignKeys: []*schema.ForeignKey{
			{
				Symbol:     "chunks_workers_chunks",
				Columns:    []*schema.Column{ChunksColumns[12]},
				RefColumns: []*schema.Column{WorkersColumns[0]},
				OnDelete:   schema.SetNull,
			},
		},
	}
	// WorkersColumns holds the columns for the "workers" table.
	WorkersColumns = []*schema.Column{
		{Name: "id", Type: field.TypeUUID},
		{Name: "created_at", Type: field.TypeTime},
		{Name: "updated_at", Type: field.TypeTime},
		{Name: "name", Type: field.TypeString, Unique: true},
		{Name: "token", Type: field.TypeString},
	}
	// WorkersTable holds the schema information for the "workers" table.
	WorkersTable = &schema.Table{
		Name:       "workers",
		Columns:    WorkersColumns,
		PrimaryKey: []*schema.Column{WorkersColumns[0]},
	}
	// Tables holds all the tables in the schema.
	Tables = []*schema.Table{
		ChunksTable,
		WorkersTable,
	}
)

func init() {
	ChunksTable.ForeignKeys[0].RefTable = WorkersTable
}
