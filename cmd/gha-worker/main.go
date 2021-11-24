package main

import (
	"context"
	"flag"
	"time"

	"github.com/go-faster/errors"
	"go.uber.org/zap"

	"github.com/go-faster/gha/internal/app"
	"github.com/go-faster/gha/internal/oas"
)

func main() {
	app.Run(func(ctx context.Context, lg *zap.Logger) error {
		lg.Info("Worker started")
		defer lg.Info("Worker stopped")
		var arg struct {
			Token string
			Addr  string
		}
		flag.StringVar(&arg.Token, "token", "", "worker token")
		flag.StringVar(&arg.Addr, "addr", "http://localhost:8080", "http listen addr")
		flag.Parse()

		client, err := oas.NewClient(arg.Addr)
		if err != nil {
			return errors.Wrap(err, "create crient")
		}

		status, err := client.Status(ctx)
		if err != nil {
			return errors.Wrap(err, "status")
		}

		lg.Info("Got status", zap.String("status", status.Message))

		for {
			if ctx.Err() != nil {
				return ctx.Err()
			}
			job, err := client.Poll(ctx, oas.PollParams{
				Token: arg.Token,
			})
			if err != nil {
				lg.Error("Poll failed", zap.Error(err))
				time.Sleep(time.Second)
				continue
			}

			switch job.Type {
			case oas.JobNothingJob:
				lg.Info("Doing nothing")
				time.Sleep(time.Second * 3)
			case oas.JobDownloadJob:
				key := job.JobDownload.Date
				lg.Info("Downloading",
					zap.String("key", key),
				)
				time.Sleep(time.Second)
			}
		}
	})
}
