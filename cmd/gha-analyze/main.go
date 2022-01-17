package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

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

	evTotal atomic.Int64
	evGood  atomic.Int64
	evError atomic.Int64
)

type Analyzer struct {
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

func (a *Analyzer) analyze(ctx context.Context, e *entry.Event) error {
	evTotal.Inc()
	if e.Full() {
		evGood.Inc()
	} else if e.Interesting() && e.Time.Year() > 2014 {
		evError.Inc()
	}
	return nil
}

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
		files := make(chan string)
		g, ctx := errgroup.WithContext(ctx)
		readJSON := speed.NewMetric()
		var last atomic.Value
		last.Store("starting")
		g.Go(func() error {
			t := time.NewTicker(time.Millisecond * 300)
			for {
				select {
				case <-ctx.Done():
					return ctx.Err()
				case <-t.C:
					fmt.Printf("[raw]: %10s [json]: %s [%10d ok=%-10d err=%-10d] ~ %s\n",
						readCompressed.ConsumeSpeed(),
						readJSON.ConsumeSpeed(),
						evTotal.Load(),
						evGood.Load(),
						evError.Load(),
						filepath.Base(last.Load().(string)),
					)
				}
			}
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
		for i := 0; i < arg.Jobs; i++ {
			g.Go(func() error {
				a := &Analyzer{
					reader: entry.NewReader(readJSON),
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
		return g.Wait()
	})
}
