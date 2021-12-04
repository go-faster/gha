package main

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"strings"

	_ "github.com/ClickHouse/clickhouse-go"
	"github.com/go-faster/errors"
	"go.uber.org/zap"

	"github.com/go-faster/gha/internal/app"
)

func main() {
	app.Run(func(ctx context.Context, lg *zap.Logger) error {
		db, err := sql.Open("clickhouse", os.Getenv("CLICKHOUSE"))
		if err != nil {
			return errors.Wrap(err, "clickhouse")
		}

		fmt.Print("Enum16(")

		q, err := db.QueryContext(ctx, `SELECT DISTINCT (language)
       FROM github_languages`)
		if err != nil {
			return errors.Wrap(err, "query")
		}

		var id int

		for q.Next() {
			var name string
			if err := q.Scan(&name); err != nil {
				return errors.Wrap(err, "scan")
			}
			if id != 0 {
				fmt.Print(", ")
			}
			id++
			name = strings.ReplaceAll(name, "'", "\\'")
			name = fmt.Sprintf("'%s'", name)
			fmt.Printf("%s = %d", name, id)
		}

		fmt.Println(")")

		return nil
	})
}
