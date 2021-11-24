package main

import (
	"context"
	"flag"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/go-faster/errors"
	"go.uber.org/zap"

	"github.com/go-faster/gha/internal/app"
	"github.com/go-faster/gha/internal/archive"
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

		dl := archive.New(
			http.DefaultClient,
			filepath.Join(os.TempDir(), "github-archive"),
			lg,
		)

		api, err := oas.NewClient(arg.Addr)
		if err != nil {
			return errors.Wrap(err, "create crient")
		}

		status, err := api.Status(ctx)
		if err != nil {
			return errors.Wrap(err, "status")
		}

		lg.Info("Got status", zap.String("status", status.Message))

		for {
			if ctx.Err() != nil {
				return ctx.Err()
			}
			job, err := api.Poll(ctx, oas.PollParams{
				XToken: arg.Token,
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

				params := oas.ProgressParams{
					XToken: arg.Token,
				}
				dlCtx, dlCancel := context.WithCancel(ctx)
				p := &archive.Progress{
					Cancel: dlCancel,
				}
				go func() {
					t := time.NewTicker(time.Millisecond * 200)
					defer t.Stop()
					defer lg.Info("Done")
					lastProgress := time.Now()
					for {
						select {
						case <-dlCtx.Done():
							return
						case <-t.C:
							wrote := p.Consume()
							lg.Info("Progress",
								zap.Float64("done", p.Ready()),
								zap.Int64("wrote", wrote),
							)
							if wrote == 0 {
								// No progress.
								if time.Since(lastProgress) > time.Second*10 {
									lg.Warn("No progress")
									p.Cancel()
								}
								continue
							}
							lastProgress = time.Now()
							report := oas.Progress{
								Key: key,
							}
							if _, err := api.Progress(ctx, report, params); err != nil {
								lg.Error("Failed to report progress", zap.Error(err))
								dlCancel()
								continue
							}
						}
					}
				}()
				result, err := dl.Download(dlCtx, archive.Options{
					Progress: p,
					Key:      key,
				})
				if err != nil {
					lg.Error("Failed to download", zap.Error(err))
					continue
				}

				lg.Info("Downloaded", zap.String("path", result.Path))
				if _, err := api.Progress(ctx, oas.Progress{
					Done: true,
					Key:  key,

					SHA256Content: oas.NewOptString(result.SHA256Data),
					SHA256Input:   oas.NewOptString(result.SHA256Input),
					SHA256Output:  oas.NewOptString(result.SHA256Output),
				}, params); err != nil {
					lg.Error("Failed to report progress", zap.Error(err))
					continue
				}
			}
		}
	})
}
