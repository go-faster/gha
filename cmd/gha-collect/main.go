package main

import (
	"context"
	"io"
	"net/http"
	_ "net/http/pprof"
	"os"
	"sort"
	"time"

	"github.com/ClickHouse/ch-go"
	"github.com/ClickHouse/ch-go/proto"
	"github.com/go-faster/errors"
	"github.com/go-faster/jx"
	"github.com/mergestat/timediff"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"

	"github.com/go-faster/gha/internal/app"
	"github.com/go-faster/gha/internal/gh"
)

type Event struct {
	Raw       jx.Raw
	CreatedAt time.Time
	ID        int64
}

type Service struct {
	lg      *zap.Logger
	batches chan []Event
}

func (c *Service) Send(ctx context.Context) error {
	for {
		if err := ctx.Err(); err != nil {
			return err
		}
		db, err := ch.Dial(ctx, ch.Options{
			Database:    "faster",
			Address:     os.Getenv("CLICKHOUSE_ADDR"),
			Compression: ch.CompressionLZ4,
		})
		if err != nil {
			return errors.Wrap(err, "dial")
		}
		softTimeout := time.Now().Add(time.Second * 30)
		/*
			CREATE TABLE github_events_raw
			(
			    id  Int64,
			    ts  DateTime32,
			    raw String CODEC (ZSTD(16))
			) ENGINE = ReplacingMergeTree
			      PARTITION BY toYYYYMMDD(ts)
			      ORDER BY (ts, id);
		*/
		var (
			colID  proto.ColInt64
			colTs  proto.ColDateTime
			colRaw proto.ColBytes
		)
		q := ch.Query{
			Body: "INSERT INTO github_events_raw VALUES",
			Input: proto.Input{
				{Name: "id", Data: &colID},
				{Name: "ts", Data: &colTs},
				{Name: "raw", Data: &colRaw},
			},
			OnInput: func(ctx context.Context) error {
				colID.Reset()
				colTs.Reset()
				colRaw.Reset()
				if time.Now().After(softTimeout) {
					return io.EOF
				}
				select {
				case <-ctx.Done():
					return ctx.Err()
				case <-time.After(time.Second * 5):
					return io.EOF
				case batch := <-c.batches:
					for _, e := range batch {
						colID.Append(e.ID)
						colTs.Append(e.CreatedAt)
						colRaw.Append(e.Raw)
					}
					return nil
				}
			},
		}
		if err := db.Do(ctx, q); err != nil {
			return errors.Wrap(err, "do")
		}
	}
}

func (c *Service) Poll(ctx context.Context) error {
	const perPage = 100

	client := gh.NewClient(http.DefaultClient, os.Getenv("GITHUB_TOKEN"))
	latestMet := make(map[int64]struct{})
	lg := c.lg.Named("poll")

	var etag string
	for {
		if err := ctx.Err(); err != nil {
			return err
		}
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
					v, err := d.Num()
					if err != nil {
						return errors.Wrap(err, "id")
					}
					id, err := v.Int64()
					if err != nil {
						return err
					}
					ev.ID = id
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
		currentMet := make(map[int64]struct{})
		for _, ev := range current {
			currentMet[ev.ID] = struct{}{}
		}
		var newEvents []Event
		for _, ev := range current {
			if _, ok := latestMet[ev.ID]; !ok {
				newEvents = append(newEvents, ev)
			}
		}
		sort.SliceStable(newEvents, func(i, j int) bool {
			a, b := newEvents[i], newEvents[j]
			// Sort ascending by created_at, id.
			if a.CreatedAt.Equal(b.CreatedAt) {
				return a.ID < b.ID
			}
			return a.CreatedAt.Before(b.CreatedAt)
		})

		select {
		case <-ctx.Done():
			return ctx.Err()
		case c.batches <- newEvents:
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
}

func main() {
	app.Run(func(ctx context.Context, lg *zap.Logger) error {
		g, ctx := errgroup.WithContext(ctx)
		s := &Service{
			batches: make(chan []Event, 1),
			lg:      lg,
		}
		g.Go(func() error {
			return s.Poll(ctx)
		})
		g.Go(func() error {
			return s.Send(ctx)
		})
		srv := &http.Server{
			Handler: http.DefaultServeMux,
			Addr:    ":8090",
		}
		g.Go(func() error {
			<-ctx.Done()
			return srv.Shutdown(ctx)
		})
		g.Go(func() error {
			if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
				return errors.Wrap(err, "http")
			}
			return nil
		})
		return g.Wait()
	})
}
