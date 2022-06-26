package main

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/schema"
	"github.com/go-faster/errors"
	_ "github.com/lib/pq"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"

	"github.com/go-faster/gha/internal/app"
	"github.com/go-faster/gha/internal/controller"
	"github.com/go-faster/gha/internal/ent"
	"github.com/go-faster/gha/internal/entry"
	"github.com/go-faster/gha/internal/oas"
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
			Init     bool
			Register string
			Addr     string
		}
		flag.BoolVar(&arg.Init, "init", false, "initialize chunks")
		flag.StringVar(&arg.Addr, "addr", "localhost:8080", "http listen addr")
		flag.StringVar(&arg.Register, "register", "", "name of worker to register")
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
		if arg.Register != "" {
			lg.Info("Registering new worker")
			tokBytes := make([]byte, 20)
			if _, err := io.ReadFull(rand.Reader, tokBytes); err != nil {
				return errors.Wrap(err, "rand")
			}

			var tokRunes []rune
			for i, c := range hex.EncodeToString(tokBytes) {
				if i > 0 && i%10 == 0 {
					tokRunes = append(tokRunes, '-')
				}
				tokRunes = append(tokRunes, c)
			}

			tok := arg.Register + ":" + string(tokRunes)
			w, err := client.Worker.Create().
				SetName(arg.Register).
				SetToken(tok).Save(ctx)
			if err != nil {
				return errors.Wrap(err, "worker create")
			}
			lg.Debug("Created")

			fmt.Println(w.Token)
			return nil
		}

		h := controller.New(client, lg)
		oasServer, err := oas.NewServer(h)
		if err != nil {
			return errors.Wrap(err, "oas server")
		}
		s := &http.Server{
			Addr:    arg.Addr,
			Handler: oasServer,
		}
		g, ctx := errgroup.WithContext(ctx)
		g.Go(func() error {
			defer lg.Info("Stopped")

			lg.Info("Listening", zap.String("http", arg.Addr))
			return s.ListenAndServe()
		})
		g.Go(func() error {
			if err := h.Run(ctx); err != nil {
				return errors.Wrap(err, "controller")
			}
			return nil
		})
		g.Go(func() error {
			<-ctx.Done()

			lg.Info("Shutting down")
			shutCtx, cancel := context.WithTimeout(context.Background(), time.Second)
			defer cancel()

			return s.Shutdown(shutCtx)
		})

		return g.Wait()
	})
}
