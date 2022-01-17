package main

import (
	"bufio"
	"context"
	"flag"
	"log"
	"math/rand"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/go-faster/errors"
	"github.com/go-faster/jx"
	"github.com/klauspost/compress/zstd"
	"go.uber.org/atomic"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"

	"github.com/go-faster/gha/internal/app"
)

func main() {
	app.Run(func(ctx context.Context, lg *zap.Logger) error {
		var arg struct {
			Jobs        int
			Samples     int
			SampleFile  float64
			SampleEntry float64
		}
		flag.IntVar(&arg.Jobs, "j", 8, "concurrent jobs")
		flag.IntVar(&arg.Samples, "n", 10, "samples to gather per file")
		flag.Float64Var(&arg.SampleFile, "f", 0.50, "probability of picking file")
		flag.Float64Var(&arg.SampleEntry, "e", 0.99, "probability of picking entry")
		flag.Parse()

		name := flag.Arg(0)
		if name == "" {
			return errors.New("name argument is required")
		}
		files := make(chan string)
		done := make(chan struct{})

		g, ctx := errgroup.WithContext(ctx)
		var (
			last      atomic.Value
			processed atomic.Int64
			write     sync.Mutex
		)
		last.Store("starting")
		dirEntries, err := os.ReadDir(name)
		if err != nil {
			return errors.Wrap(err, "read dir")
		}
		log.Println("dir", name, "entries", len(dirEntries))

		g.Go(func() error {
			t := time.NewTicker(time.Millisecond * 300)
			for {
				select {
				case <-ctx.Done():
					return ctx.Err()
				case <-done:
					return nil
				case <-t.C:
					log.Printf("%10d / %10d ~ %s\n",
						processed.Load(),
						len(dirEntries),
						filepath.Base(last.Load().(string)),
					)
				}
			}
		})
		g.Go(func() error {
			defer close(files)
			defer close(done)

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
				j := jx.GetDecoder()
				buf := make([]byte, 1024*1024*50)
				process := func(filename string) error {
					f, err := os.Open(filename)
					if err != nil {
						return errors.Wrap(err, "open")
					}
					defer func() { _ = f.Close() }()
					z, err := zstd.NewReader(f, zstd.WithDecoderConcurrency(1))
					if err != nil {
						return errors.Wrap(err, "zstd")
					}
					defer z.Close()

					s := bufio.NewScanner(z)
					s.Buffer(buf, len(buf))

					var (
						count int
						i     int
					)
					for s.Scan() {
						i++
						data := s.Bytes()
						j.ResetBytes(data)
						if rand.Float64() < arg.SampleEntry {
							continue
						}
						if j.Validate() == nil {
							write.Lock()
							_, _ = os.Stdout.Write(append(data, '\n'))
							write.Unlock()
							count++
						}
						if err := ctx.Err(); err != nil {
							return err
						}
						if count >= arg.Samples {
							break
						}
					}

					return s.Err()
				}
				for {
					select {
					case n, ok := <-files:
						if !ok {
							return nil
						}
						if rand.Float64() > arg.SampleFile {
							if err := process(n); err != nil {
								return errors.Wrap(err, "analyze")
							}
						}
						last.Store(n)
						processed.Inc()
					case <-ctx.Done():
						return ctx.Err()
					}
				}
			})
		}
		return g.Wait()
	})
}
