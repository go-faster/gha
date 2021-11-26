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
	"golang.org/x/sync/errgroup"

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
			Jobs  int
			Dir   string
		}
		flag.IntVar(&arg.Jobs, "j", 4, "maximum concurrent jobs")
		flag.StringVar(&arg.Token, "token", "", "worker token")
		flag.StringVar(&arg.Dir, "dir", filepath.Join(os.TempDir(), "github-archive"), "output dir")
		flag.StringVar(&arg.Addr, "addr", "http://localhost:8080", "http listen addr")
		flag.Parse()

		dl := archive.New(
			&http.Client{
				Timeout: time.Second * 30,
				Transport: &http.Transport{
					TLSHandshakeTimeout: time.Second * 3,
					MaxConnsPerHost:     10,
					MaxIdleConnsPerHost: 100,
					MaxIdleConns:        100,
				},
			},
			arg.Dir,
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

		jobs := make(chan oas.Job)

		g, ctx := errgroup.WithContext(ctx)
		g.Go(func() error {
			defer close(jobs)
			for {
				if err := ctx.Err(); err != nil {
					return err
				}
				job, err := api.Poll(ctx, oas.PollParams{
					XToken: arg.Token,
				})
				if err != nil {
					lg.Error("Poll failed", zap.Error(err))
					time.Sleep(time.Second)
					continue
				}
				if !job.IsJobInventory() {
					// Simple jobs that does not require split.
					jobs <- job
					continue
				}
				// Split inventory into sub-jobs.
				for _, k := range job.JobInventory.Date {
					jobs <- oas.NewJobInventoryJob(oas.JobInventory{
						Type: job.JobInventory.Type,
						Date: []string{k},
					})
				}
			}
		})

		handleInventory := func(ctx context.Context, j oas.JobInventory) error {
			if len(j.Date) == 0 {
				// Nothing to do.
				lg.Warn("Got blank inventory")
				return nil
			}
			if len(j.Date) != 1 {
				lg.Warn("Got inventory with unexpected key count")
			}

			key := j.Date[0]
			lg.Info("Inventory",
				zap.String("key", key),
			)
			params := oas.ProgressParams{
				XToken: arg.Token,
			}

			res, err := dl.Inventory(ctx, key)
			if errors.Is(err, archive.ErrNotFound) {
				return nil
			}
			if err != nil {
				return errors.Wrap(err, "inventory")
			}
			if _, err := api.Progress(ctx, oas.Progress{
				Event: oas.ProgressEventInventory,
				Key:   key,

				OutputSizeBytes:  oas.NewOptInt64(res.SizeOutput),
				ContentSizeBytes: oas.NewOptInt64(res.SizeContent),
				InputSizeBytes:   oas.NewOptInt64(res.SizeInput),
			}, params); err != nil {
				return errors.Wrap(err, "report")
			}

			return nil
		}

		handleDownload := func(ctx context.Context, j oas.JobDownload) error {
			key := j.Date
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
				t := time.NewTicker(time.Second * 3)
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
							zap.String("key", key),
						)
						if wrote == 0 {
							// No progress.
							if time.Since(lastProgress) > time.Second*5 {
								lg.Warn("No progress")
							}
							if time.Since(lastProgress) > time.Second*10 {
								p.Cancel()
							}
							continue
						}
						lastProgress = time.Now()
						if _, err := api.Progress(ctx, oas.Progress{
							Event:           oas.ProgressEventDownloading,
							Key:             key,
							InputSizeBytes:  oas.NewOptInt64(p.Total()),
							InputReadyBytes: oas.NewOptInt64(p.ReadyBytes()),
						}, params); err != nil {
							lg.Error("Failed to report progress", zap.Error(err))
							dlCancel()
							continue
						}
					}
				}
			}()

			res, err := dl.Download(dlCtx, archive.Options{
				Progress: p,
				Key:      key,
			})
			if err != nil {
				return errors.Wrap(err, "download")
			}

			lg.Info("Downloaded", zap.String("path", res.Path))
			if _, err := api.Progress(ctx, oas.Progress{
				Event: oas.ProgressEventDone,
				Key:   key,

				SHA256Content: oas.NewOptString(res.SHA256Content),
				SHA256Input:   oas.NewOptString(res.SHA256Input),
				SHA256Output:  oas.NewOptString(res.SHA256Output),

				OutputSizeBytes:  oas.NewOptInt64(res.SizeOutput),
				ContentSizeBytes: oas.NewOptInt64(res.SizeContent),
				InputReadyBytes:  oas.NewOptInt64(res.SizeInput),
			}, params); err != nil {
				return errors.Wrap(err, "report")
			}

			return nil
		}
		handleJob := func(ctx context.Context, j oas.Job) error {
			if j.IsJobNothing() {
				lg.Info("Doing nothing")
				time.Sleep(time.Second * 3)
				return nil
			}
			if j, ok := j.GetJobInventory(); ok {
				return handleInventory(ctx, j)
			}
			if j, ok := j.GetJobDownload(); ok {
				return handleDownload(ctx, j)
			}
			return errors.Errorf("unknown job: %v", j.Type)
		}

		for i := 0; i < arg.Jobs; i++ {
			g.Go(func() error {
				for job := range jobs {
					if err := handleJob(ctx, job); err != nil {
						lg.Warn("Job failed", zap.Error(err))
						continue
					}
				}
				return nil
			})
		}
		return g.Wait()
	})
}
