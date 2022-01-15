package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/go-faster/errors"
	"go.uber.org/zap"

	"github.com/go-faster/gha/internal/app"
	"github.com/go-faster/gha/internal/entry"
)

func main() {
	app.Run(func(ctx context.Context, lg *zap.Logger) error {
		flag.Parse()
		name := flag.Arg(0)
		if name == "" {
			return errors.New("name argument is required")
		}
		start := time.Now()
		f, err := os.Open(name)
		if err != nil {
			return errors.Wrap(err, "open")
		}
		defer func() { _ = f.Close() }()

		r := entry.NewReader()
		defer r.Close()

		var count int
		defer func() {
			fmt.Println("count", count, "duration", time.Since(start).Round(time.Millisecond))
		}()

		return r.Decode(ctx, f, func(ctx context.Context, e *entry.Event) error {
			count++
			return nil
		})
	})
}
