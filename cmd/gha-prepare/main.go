package main

import (
	"bufio"
	"context"
	"database/sql"
	"flag"
	"fmt"
	"io"
	"os"
	"os/signal"
	"path/filepath"
	"runtime"
	"runtime/pprof"
	"strings"
	"time"

	_ "github.com/ClickHouse/clickhouse-go"
	"github.com/cheggaaa/pb/v3"
	"github.com/dustin/go-humanize"
	"github.com/go-faster/errors"
	"github.com/go-faster/jx"
	"github.com/klauspost/compress/zstd"
	"go.uber.org/atomic"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"

	"github.com/go-faster/gha/internal/app"
	"github.com/go-faster/gha/internal/entry"
)

type Metric struct {
	last  time.Time
	read  atomic.Uint64
	total atomic.Uint64
}

func (s *Metric) Report(ctx context.Context, name string) func() error {
	return func() error {
		t := time.NewTicker(time.Millisecond * 300)
		for {
			select {
			case <-ctx.Done():
				return ctx.Err()
			case <-t.C:
				_ = s.ConsumeSpeed()
			}
		}
	}
}

func (s *Metric) Write(p []byte) (n int, err error) {
	n = len(p)
	s.Add(uint64(n))
	return n, nil
}

func (s *Metric) Add(n uint64) {
	s.read.Add(n)
	s.total.Add(n)
}

func (s *Metric) Consume() uint64 {
	n := s.read.Load()
	s.Add(-n)
	return n
}

func (s *Metric) ConsumeSpeed() string {
	now := time.Now()
	delta := now.Sub(s.last)
	s.last = now
	n := readUncompressed.Consume()

	bytesPerSec := float64(n) / delta.Seconds()

	return fmt.Sprintf("%s / sec", humanize.Bytes(uint64(bytesPerSec)))
}

func NewMetric() *Metric {
	return &Metric{last: time.Now()}
}

var (
	readCompressed   = NewMetric()
	readUncompressed = NewMetric()
)

func process(ctx context.Context, lg *zap.Logger, name string, events chan entry.Event) error {
	f, err := os.Open(name)
	if err != nil {
		return errors.Wrap(err, "open")
	}
	defer func() {
		_ = f.Close()
	}()

	r, err := zstd.NewReader(io.TeeReader(f, readCompressed))
	if err != nil {
		return errors.Wrap(err, "zstd")
	}
	defer r.Close()

	s := bufio.NewScanner(io.TeeReader(r, readUncompressed))

	const maxSize = 1024 * 1024 * 150
	s.Buffer(make([]byte, maxSize), maxSize)

	d := jx.GetDecoder()

	for s.Scan() {
		if err := ctx.Err(); err != nil {
			return err
		}

		b := s.Bytes()
		d.ResetBytes(b)

		var event entry.Event
		if err := event.Decode(d); err != nil {
			lg.Error("Invalid entry",
				zap.Error(err),
				zap.String("path", name),
			)
			out, err := os.CreateTemp("", "prepare-*.json")
			if err == nil {
				_, _ = out.WriteString(s.Text())
				name := out.Name()
				_ = out.Close()
				_, _ = fmt.Fprintln(os.Stderr, "Input dumped to", name)
			}
			continue
		}

		if !event.Interesting() {
			continue
		}

		select {
		case <-ctx.Done():
			return ctx.Err()
		case events <- event:
			continue
		}
	}

	if err := s.Err(); err != nil {
		return errors.Wrapf(err, "read: %s", name)
	}

	return nil
}

func main() {
	app.Run(func(ctx context.Context, lg *zap.Logger) error {
		var arg struct {
			Jobs int
			CH   string
			Dry  bool

			CPUProfile string
		}
		flag.IntVar(&arg.Jobs, "j", runtime.GOMAXPROCS(-1), "concurrent jobs")
		flag.BoolVar(&arg.Dry, "dry", false, "dry run")
		flag.StringVar(&arg.CH, "clickhouse", "tcp://127.0.0.1:9000", "clickhouse target")
		flag.StringVar(&arg.CPUProfile, "cpuprofile", "", "write cpu profile to `file`")
		flag.Parse()

		start := time.Now()
		defer func() {
			fmt.Println("duration", time.Since(start))
		}()

		if arg.CPUProfile != "" {
			f, err := os.Create(arg.CPUProfile)
			if err != nil {
				return err
			}
			defer func() {
				_ = f.Close()
			}()
			if err := pprof.StartCPUProfile(f); err != nil {
				return errors.Wrap(err, "start cpu profile")
			}
			defer pprof.StopCPUProfile()
		}

		db, err := sql.Open("clickhouse", arg.CH)
		if err != nil {
			return errors.Wrap(err, "clickhouse")
		}

		events := make(chan entry.Event, 1024*10)
		files := make(chan string)

		ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
		defer cancel()

		g, ctx := errgroup.WithContext(ctx)
		g.Go(readUncompressed.Report(ctx, "[json]"))
		g.Go(readCompressed.Report(ctx, "[raw]"))

		if arg.Dry {
			g.Go(func() error {
				defer cancel()

				for range events {
				}
				return nil
			})
		} else {
			g.Go(func() error {
				var (
					tx *sql.Tx
					s  *sql.Stmt
				)
				begin := func() error {
					var err error
					if tx, err = db.Begin(); err != nil {
						return errors.Wrap(err, "begin")
					}
					if s, err = tx.Prepare("INSERT INTO events(event_type, repo_name, created_at) VALUES (?, ?, ?)"); err != nil {
						return errors.Wrap(err, "prepare")
					}

					return nil
				}
				commit := func() error {
					if err := tx.Commit(); err != nil {
						return errors.Wrap(err, "commit")
					}
					return nil
				}
				if err := begin(); err != nil {
					return errors.Wrap(err, "begin")
				}

				var (
					lastInsert time.Time
				)
				for ev := range events {
					if _, err := s.Exec(ev.Type, ev.Repo, ev.Time); err != nil {
						return errors.Wrap(err, "exec")
					}
					if time.Since(lastInsert) < time.Second {
						continue
					}

					lastInsert = time.Now()
					if err := commit(); err != nil {
						return errors.Wrap(err, "commit")
					}

					// Re-init.
					if err := begin(); err != nil {
						return errors.Wrap(err, "commit")
					}
				}
				if err := commit(); err != nil {
					return errors.Wrap(err, "commit")
				}
				return nil
			})
		}
		g.Go(func() error {
			defer close(files)

			dir := flag.Arg(0)
			entries, err := os.ReadDir(dir)
			if err != nil {
				return errors.Wrap(err, "list")
			}

			var (
				paths []string
				total int64
			)
			size := make(map[string]int64, len(entries))
			for _, e := range entries {
				if e.IsDir() {
					continue
				}
				if !strings.HasSuffix(e.Name(), ".json.zst") {
					continue
				}

				p := filepath.Join(dir, e.Name())
				stat, err := os.Stat(p)
				if err != nil {
					return errors.Wrap(err, "stat")
				}
				total += stat.Size()
				size[p] = stat.Size()
				paths = append(paths, p)
			}

			bar := pb.Full.Start64(total)
			bar.Set(pb.Bytes, true)

			for _, p := range paths {
				select {
				case <-ctx.Done():
					return ctx.Err()
				case files <- p:
					bar.Add64(size[p])
					continue
				}
			}
			bar.Finish()

			return nil
		})
		g.Go(func() error {
			defer close(events)

			g, ctx := errgroup.WithContext(ctx)
			for i := 0; i < arg.Jobs; i++ {
				g.Go(func() error {
					for path := range files {
						if err := ctx.Err(); err != nil {
							return ctx.Err()
						}
						if err := process(ctx, lg, path, events); err != nil {
							if errors.Is(err, bufio.ErrTooLong) {
								lg.Warn("Token too long", zap.String("path", path))
								continue
							}
							return errors.Wrap(err, "process")
						}
					}

					return nil
				})
			}

			return g.Wait()
		})

		return g.Wait()
	})
}
