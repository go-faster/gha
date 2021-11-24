package main

import (
	"context"

	"go.uber.org/zap"

	"github.com/go-faster/gha/internal/app"
)

func main() {
	app.Run(func(ctx context.Context, lg *zap.Logger) error {
		lg.Info("Controller started")
		return nil
	})
}
