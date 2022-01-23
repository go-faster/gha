package main

import (
	"context"

	"github.com/go-faster/ch"
	"github.com/go-faster/errors"
	"go.uber.org/zap"

	"github.com/go-faster/gha/internal/app"
)

func main() {
	app.Run(func(ctx context.Context, lg *zap.Logger) error {
		conn, err := ch.Dial(ctx, ch.Options{
			Logger: lg,
		})
		if err != nil {
			return errors.Wrap(err, "dial")
		}
		return conn.Close()
	})
}
