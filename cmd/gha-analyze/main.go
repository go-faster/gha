package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/go-faster/ch"
	"github.com/go-faster/ch/proto"
	"github.com/go-faster/errors"
	"go.uber.org/atomic"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"

	"github.com/go-faster/gha/internal/app"
	"github.com/go-faster/gha/internal/entry"
	"github.com/go-faster/gha/internal/speed"
)

var (
	readCompressed = speed.NewMetric()
	chSent         = speed.NewMetric()

	evTotal atomic.Int64
	evGood  atomic.Int64
	evError atomic.Int64
	evSent  atomic.Int64
)

type Analyzer struct {
	events chan<- *Events
	buf    *Events
	reader *entry.Reader
}

func (a *Analyzer) Analyze(ctx context.Context, filename string) error {
	f, err := os.Open(filename)
	if err != nil {
		return errors.Wrap(err, "open")
	}
	defer func() { _ = f.Close() }()

	return a.reader.Decode(ctx, entry.Decode{
		Reader:  io.TeeReader(f, readCompressed),
		OnEvent: a.analyze,
		OnError: func(ctx context.Context, data []byte, err error) error {
			evError.Inc()
			return nil
		},
	})
}

func (a *Analyzer) Schedule(ctx context.Context, e *entry.Event) error {
	a.buf.Append(e)
	if a.buf.Date.Rows() < 10_000 {
		return nil
	}
	d := a.buf
	a.buf = getEvents()
	select {
	case a.events <- d:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

func (a *Analyzer) analyze(ctx context.Context, e *entry.Event) error {
	evTotal.Inc()
	if !e.Interesting() {
		return nil
	}
	if e.Full() {
		evGood.Inc()
		if err := a.Schedule(ctx, e); err != nil {
			return errors.Wrap(err, "schedule")
		}
	} else if e.Time.Year() > 2014 {
		evError.Inc()
	}
	return nil
}

type Events struct {
	Date  proto.ColDateTime
	Type  proto.ColEnum8
	Actor proto.ColStr
	Repo  proto.ColStr
}

func (v *Events) Reset() {
	v.Date.Reset()
	v.Type.Reset()
	v.Actor.Reset()
	v.Repo.Reset()
}

var pool = sync.Pool{
	New: func() interface{} {
		return new(Events)
	},
}

func getEvents() *Events {
	return pool.Get().(*Events)
}

func putEvents(e *Events) {
	e.Reset()
	pool.Put(e)
}

const enum = `Enum8('WatchEvent'=0, 'PushEvent'=1, 'IssuesEvent'=2, 'PullRequestEvent'=3)`

func (v *Events) Input() proto.Input {
	return proto.Input{
		{Name: "date", Data: &v.Date},
		{Name: "type", Data: proto.Alias(&v.Type, enum)},
		{Name: "actor", Data: &v.Actor},
		{Name: "repo", Data: &v.Repo},
	}
}

func (v *Events) Append(e *entry.Event) {
	v.Date.Append(proto.ToDateTime(e.Time))
	// "WatchEvent", "PushEvent", "IssuesEvent", "PullRequestEvent"
	// Enum8('WatchEvent'=0, 'PushEvent'=1, 'IssuesEvent'=2, 'PullRequestEvent'=3),
	switch string(e.Type) {
	case "WatchEvent":
		v.Type.Append(0)
	case "PushEvent":
		v.Type.Append(1)
	case "IssuesEvent":
		v.Type.Append(2)
	case "PullRequestEvent":
		v.Type.Append(3)
	default:
		panic(string(e.Type))
	}
	v.Actor.AppendBytes(e.Actor)
	v.Repo.AppendBytes(e.Repo)
}

const ddl = `CREATE TABLE IF NOT EXISTS github.events  (
    date DateTime,
    type Enum8('WatchEvent'=0, 'PushEvent'=1, 'IssuesEvent'=2, 'PullRequestEvent'=3),
    actor String,
    repo LowCardinality(String)
) ENGINE MergeTree() ORDER BY (type, date, repo)`

func main() {
	app.Run(func(ctx context.Context, lg *zap.Logger) error {
		var arg struct {
			Jobs int
		}
		flag.IntVar(&arg.Jobs, "j", 8, "concurrent jobs")
		flag.Parse()
		name := flag.Arg(0)
		if name == "" {
			return errors.New("name argument is required")
		}

		db, err := ch.Dial(ctx, "localhost:9000", ch.Options{
			Logger:      lg.WithOptions(zap.IncreaseLevel(zap.InfoLevel)),
			Compression: ch.CompressionZSTD,
		})
		if err != nil {
			return errors.Wrap(err, "clickhouse dial")
		}
		if err := db.Ping(ctx); err != nil {
			return errors.Wrap(err, "ping")
		}
		if err := db.Do(ctx, ch.Query{Body: `CREATE DATABASE IF NOT EXISTS github`}); err != nil {
			return errors.Wrap(err, "create db")
		}
		if err := db.Do(ctx, ch.Query{Body: ddl}); err != nil {
			return errors.Wrap(err, "ddl")
		}
		if err := db.Do(ctx, ch.Query{Body: `TRUNCATE TABLE github.events`}); err != nil {
			return errors.Wrap(err, "truncate")
		}
		lg.Info("Database OK")

		files := make(chan string)
		done := make(chan struct{})
		start := time.Now()

		defer func() {
			fmt.Println("Done", time.Since(start).Round(time.Second))
		}()

		g, ctx := errgroup.WithContext(ctx)
		readJSON := speed.NewMetric()
		var last atomic.Value
		last.Store("starting")
		g.Go(func() error {
			t := time.NewTicker(time.Second * 3)
			for {
				select {
				case <-ctx.Done():
					return ctx.Err()
				case <-done:
					return nil
				case <-t.C:
					fmt.Printf("[raw %s] [json %s] [read=%-12d err=%-03d] [Δ %7d] [ch %.0f / sec] ~ %s\n",
						readCompressed.ConsumeSpeed(),
						readJSON.ConsumeSpeed(),
						evGood.Load(),
						evError.Load(),
						evGood.Load()-evSent.Load(),
						chSent.Rate(),
						filepath.Base(last.Load().(string)),
					)
				}
			}
		})

		events := make(chan *Events, arg.Jobs)
		g.Go(func() error {
			defer close(done)

			var e Events
			if err := db.Do(ctx, ch.Query{
				Body:  "INSERT INTO github.events VALUES",
				Input: e.Input(),
				OnProgress: func(ctx context.Context, p proto.Progress) error {
					lg.Info("Progress")
					return nil
				},
				OnProfile: func(ctx context.Context, p proto.Profile) error {
					lg.Info("Profile")
					return nil
				},
				OnInput: func(ctx context.Context) error {
					e.Reset()
					select {
					case got, ok := <-events:
						if !ok {
							return io.EOF
						}

						// Copy.
						e.Actor.Pos = append(e.Actor.Pos, got.Actor.Pos...)
						e.Actor.Buf = append(e.Actor.Buf, got.Actor.Buf...)
						e.Repo.Pos = append(e.Repo.Pos, got.Repo.Pos...)
						e.Repo.Buf = append(e.Repo.Buf, got.Repo.Buf...)
						e.Date = append(e.Date, got.Date...)
						e.Type = append(e.Type, got.Type...)
						putEvents(got)

						chSent.Add(uint64(e.Type.Rows()))
						evSent.Add(int64(e.Type.Rows()))

						return nil
					case <-ctx.Done():
						return ctx.Err()
					}
				},
			}); err != nil {
				return errors.Wrap(err, "insert")
			}
			return nil
		})

		g.Go(func() error {
			defer close(files)
			dirEntries, err := os.ReadDir(name)
			if err != nil {
				return errors.Wrap(err, "read dir")
			}
			fmt.Println("dir", name, "entries", len(dirEntries))
			for _, f := range dirEntries {
				n := filepath.Join(name, f.Name())
				if !strings.HasSuffix(n, ".json.zst") {
					continue
				}
				select {
				case files <- n:
				case <-ctx.Done():
					return ctx.Err()
				}
			}
			return nil
		})

		g.Go(func() error {
			defer close(events)
			aG, ctx := errgroup.WithContext(ctx)
			for i := 0; i < arg.Jobs; i++ {
				aG.Go(func() error {
					a := &Analyzer{
						reader: entry.NewReader(readJSON),
						buf:    getEvents(),
						events: events,
					}
					for {
						select {
						case n, ok := <-files:
							if !ok {
								return nil
							}
							last.Store(n)
							if err := a.Analyze(ctx, n); err != nil {
								return errors.Wrap(err, "analyze")
							}
						case <-ctx.Done():
							return ctx.Err()
						}
					}
				})
			}
			return aG.Wait()
		})

		return g.Wait()
	})
}
