package main

import (
	"bufio"
	"context"
	"database/sql"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	_ "github.com/ClickHouse/clickhouse-go"
	"github.com/cheggaaa/pb/v3"
	"github.com/go-faster/errors"
	"github.com/go-faster/jx"
	"github.com/klauspost/compress/zstd"
	"golang.org/x/sync/errgroup"

	"github.com/go-faster/gha/internal/entry"
)

func process(ctx context.Context, name string, events chan entry.Event) error {
	f, err := os.Open(name)
	if err != nil {
		return errors.Wrap(err, "open")
	}
	defer func() {
		_ = f.Close()
	}()

	r, err := zstd.NewReader(bufio.NewReader(f))
	if err != nil {
		return errors.Wrap(err, "zstd")
	}
	defer r.Close()

	s := bufio.NewScanner(r)
	d := jx.GetDecoder()

	for s.Scan() {
		d.ResetBytes(s.Bytes())

		var event entry.Event
		if err := event.Decode(d); err != nil {
			out, err := os.CreateTemp("", "prepare-*.json")
			if err == nil {
				_, _ = out.WriteString(s.Text())
				name := out.Name()
				_ = out.Close()
				_, _ = fmt.Fprintln(os.Stderr, "Input dumped to", name)
			}
			return errors.Wrap(err, "decode")
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

	return nil
}

func run() error {
	var arg struct {
		Jobs int
		CH   string
	}
	flag.IntVar(&arg.Jobs, "j", runtime.GOMAXPROCS(-1), "concurrent jobs")
	flag.StringVar(&arg.CH, "clickhouse", "tcp://127.0.0.1:9000", "clickhouse target")
	flag.Parse()

	db, err := sql.Open("clickhouse", arg.CH)
	if err != nil {
		return errors.Wrap(err, "clickhouse")
	}

	events := make(chan entry.Event, 1024*10)
	files := make(chan string)
	g, ctx := errgroup.WithContext(context.Background())
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
			if s, err = tx.Prepare("INSERT INTO events(event_type, actor_login, repo_name, created_at) VALUES (?, ?, ?, ?)"); err != nil {
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
			if _, err := s.Exec(ev.Type, ev.Actor, ev.Repo, ev.Time); err != nil {
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
					if err := process(ctx, path, events); err != nil {
						return errors.Wrap(err, "process")
					}
				}

				return nil
			})
		}

		return g.Wait()
	})

	return g.Wait()
}

func main() {
	if err := run(); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Failed: %+v\n", err)
		os.Exit(2)
	}
}
