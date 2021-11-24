package main

import (
	"context"
	"flag"
	"os"
	"time"

	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/schema"
	"github.com/go-faster/errors"
	_ "github.com/lib/pq"
	"go.uber.org/zap"

	"github.com/go-faster/gha/internal/app"
	"github.com/go-faster/gha/internal/ent"
	"github.com/go-faster/gha/internal/entry"
)

func ensureRollback(tx *ent.Tx) func() {
	var done bool
	tx.OnRollback(func(rollbacker ent.Rollbacker) ent.Rollbacker {
		done = true
		return rollbacker
	})
	tx.OnCommit(func(committer ent.Committer) ent.Committer {
		done = true
		return committer
	})
	return func() {
		if done {
			return
		}
		_ = tx.Rollback()
	}
}

func main() {
	app.Run(func(ctx context.Context, lg *zap.Logger) error {
		lg.Info("Controller started")
		dsn, ok := os.LookupEnv("DATABASE_DSN")
		if !ok {
			return errors.New(`no DATABASE_DSN provided, example "user=gha password=gha database=gha port=5433"`)
		}
		client, err := ent.Open("postgres", dsn)
		if err != nil {
			return errors.Wrap(err, "db open")
		}
		defer func() {
			if err := client.Close(); err != nil {
				lg.Error("db: close", zap.Error(err))
			}
		}()
		if err := client.Schema.Create(ctx,
			schema.WithDropIndex(true),
			schema.WithDropColumn(true),
		); err != nil {
			return errors.Wrap(err, "db schema")
		}

		lg.Info("Schema ok")

		var arg struct {
			Init bool
		}
		flag.BoolVar(&arg.Init, "init", false, "initialize chunks")
		flag.Parse()

		if arg.Init {
			lg.Info("Initializing chunks")
			now := time.Now()
			var (
				start = entry.Start()
				end   = now.AddDate(0, 0, 2)
			)
			tx, err := client.BeginTx(ctx, &sql.TxOptions{})
			if err != nil {
				return errors.Wrap(err, "begin")
			}
			defer ensureRollback(tx)()
			for current := start; current.Before(end); current = current.Add(entry.Delta) {
				key := current.Format(entry.Layout)

				_, err := tx.Chunk.Create().
					SetID(key).
					SetStart(current).
					SetCreatedAt(now).
					SetUpdatedAt(now).Save(ctx)
				if err != nil {
					return errors.Wrap(err, "save")
				}
			}

			if err := tx.Commit(); err != nil {
				return errors.Wrap(err, "commit")
			}

			lg.Info("Done")
		}

		return nil
	})
}
