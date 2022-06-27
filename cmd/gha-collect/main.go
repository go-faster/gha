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
	"github.com/mergestat/timediff"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"

	"github.com/go-faster/gha/internal/app"
	"github.com/go-faster/gha/internal/gh"
)

type Service struct {
	lg      *zap.Logger
	batches chan []gh.Event
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
		lg.Info("Fetching events")
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

		rt := res.RateLimit

		// De-duplicating events.
		var newEvents []gh.Event
		currentMet := make(map[int64]struct{})
		for _, ev := range res.Data {
			currentMet[ev.ID] = struct{}{}
			if _, ok := latestMet[ev.ID]; !ok {
				newEvents = append(newEvents, ev)
			}
		}
		if len(newEvents) >= perPage && etag != "" {
			lg.Info("Fetching more events")
			resNext, err := client.Events(ctx, gh.Params{
				PerPage: perPage,
				Page:    2,
			})
			if err != nil {
				return errors.Wrap(err, "failed to fetch events")
			}
			rt = resNext.RateLimit
			for _, ev := range resNext.Data {
				if _, ok := currentMet[ev.ID]; ok {
					continue
				}
				currentMet[ev.ID] = struct{}{}
				if _, ok := latestMet[ev.ID]; !ok {
					newEvents = append(newEvents, ev)
				}
			}
			lg.Info("Additional events loaded",
				zap.Int("new_events_count", len(newEvents)),
			)
			if len(newEvents) >= (perPage * 2) {
				// No sense to continue if we have more than 2 pages of events,
				// request duration will be too long, and we can miss even more.
				lg.Warn("Unable to resolve missing events")
			} else {
				lg.Info("Missing events resolved")
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
		if rt.Remaining < 10 {
			lg.Warn("Rate limit", zap.Int("remaining", rt.Remaining))
			sleep = rt.Reset.Sub(time.Now()) + time.Second
		} else {
			sleep = time.Until(rt.Reset) / time.Duration(rt.Remaining)
		}
		duration := time.Since(start)
		sleep -= duration // don't sleep for more than the rate limit
		if sleep <= 0 {
			sleep = 0
		}
		lg.Info("Events",
			zap.Duration("duration", duration),
			zap.Int("new_count", len(newEvents)),
			zap.Int("remaining", rt.Remaining),
			zap.Int("used", rt.Used),
			zap.Duration("reset", rt.Reset.Sub(time.Now())),
			zap.String("reset_human", timediff.TimeDiff(rt.Reset)),
			zap.String("sleep", sleep.String()),
		)
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
			batches: make(chan []gh.Event, 5),
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
