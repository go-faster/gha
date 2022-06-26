package main

import (
	"context"
	"net/http"
	"os"
	"time"

	"github.com/go-faster/errors"
	"github.com/go-faster/jx"
	"github.com/mergestat/timediff"
	"go.uber.org/zap"

	"github.com/go-faster/gha/internal/app"
	"github.com/go-faster/gha/internal/gh"
)

type Event struct {
	Raw       jx.Raw
	CreatedAt time.Time
	ID        string
}

func main() {
	app.Run(func(ctx context.Context, lg *zap.Logger) error {
		const perPage = 100

		client := gh.NewClient(http.DefaultClient, os.Getenv("GITHUB_TOKEN"))
		latestMet := make(map[string]struct{})
		var etag string
		for {
			start := time.Now()
			res, err := client.Events(ctx, gh.Params{
				PerPage: perPage,
				Etag:    etag,
			})
			if err != nil {
				return errors.Wrap(err, "failed to fetch events")
			}
			if res.NotModified {
				lg.Info("Not modified", zap.Duration("duration", time.Since(start)))
				continue
			}

			var current []Event
			for _, v := range res.Data {
				ev := Event{
					Raw: v,
				}
				if err := jx.DecodeBytes(v).ObjBytes(func(d *jx.Decoder, k []byte) error {
					switch string(k) {
					case "created_at":
						v, err := d.Str()
						if err != nil {
							return errors.Wrap(err, "str")
						}
						if ev.CreatedAt, err = time.Parse(time.RFC3339, v); err != nil {
							return err
						}
						return nil
					case "id":
						v, err := d.Str()
						if err != nil {
							return errors.Wrap(err, "id")
						}
						ev.ID = v
						return nil
					default:
						if err := d.Skip(); err != nil {
							return errors.Wrap(err, "skip")
						}
						return nil
					}
				}); err != nil {
					return err
				}
				current = append(current, ev)
			}

			// De-duplicating events.
			currentMet := make(map[string]struct{})
			for _, ev := range current {
				currentMet[ev.ID] = struct{}{}
			}
			var newEvents []Event
			for _, ev := range current {
				if _, ok := latestMet[ev.ID]; !ok {
					newEvents = append(newEvents, ev)
				}
			}

			// Calculating next sleep time to avoid rate limit.
			var sleep time.Duration
			if res.RateLimit.Remaining < 10 {
				lg.Warn("Rate limit", zap.Int("remaining", res.RateLimit.Remaining))
				sleep = res.RateLimit.Reset.Sub(time.Now()) + time.Second
			} else {
				sleep = time.Until(res.RateLimit.Reset) / time.Duration(res.RateLimit.Remaining)
			}
			duration := time.Since(start)
			sleep -= duration // don't sleep for more than the rate limit
			if sleep <= 0 {
				sleep = 0
			}
			lg.Info("Events",
				zap.Duration("duration", duration),
				zap.Int("new_count", len(newEvents)),
				zap.Int("remaining", res.RateLimit.Remaining),
				zap.Int("used", res.RateLimit.Used),
				zap.Duration("reset", res.RateLimit.Reset.Sub(time.Now())),
				zap.String("reset_human", timediff.TimeDiff(res.RateLimit.Reset)),
				zap.String("sleep", sleep.String()),
			)
			if len(newEvents) >= perPage && etag != "" {
				lg.Warn("Missed events")
			}
			select {
			case <-time.After(sleep):
				latestMet = currentMet
				etag = res.Etag
			case <-ctx.Done():
				return ctx.Err()
			}
		}
	})
}
