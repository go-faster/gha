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
			Compression: ch.CompressionZSTD,
		})
		if err != nil {
			return errors.Wrap(err, "dial")
		}
		softTimeout := time.Now().Add(time.Minute)
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
			colID  proto.ColInt64    // id Int64
			colTs  proto.ColDateTime // ts DateTime
			colRaw proto.ColBytes    // raw String
		)
		// Stream events to ClickHouse.
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
	const (
		perPage  = 100
		maxPages = 10
	)

	client := gh.NewClient(http.DefaultClient, os.Getenv("GITHUB_TOKEN"))
	latestMet := make(map[int64]struct{})
	lg := c.lg.Named("poll")

	var etag string
Fetch:
	for {
		if err := ctx.Err(); err != nil {
			return err
		}

		var rt gh.RateLimit
		var newEvents []gh.Event
		var start time.Time

		currentMet := make(map[int64]struct{})
		for i := 0; i <= maxPages; i++ {
			lg.Info("Fetching events")
			start = time.Now()
			p := gh.Params{
				Page:    i + 1, // first page is 1
				PerPage: perPage,
			}
			if i == 0 {
				p.Etag = etag
			}
			res, err := client.Events(ctx, p)
			if err != nil {
				return errors.Wrap(err, "failed to fetch events")
			}
			if res.NotModified {
				lg.Info("Not modified", zap.Duration("duration", time.Since(start)))
				continue Fetch
			}
			if res.Unprocessable {
				lg.Warn("Unable to resolve missing events")
				break
			}

			// Updating rate-limit to sleep later.
			rt = res.RateLimit

			// Searching for new events.
			// The currentMet contains events from previous Fetch loop.
			for _, ev := range res.Data {
				if _, ok := currentMet[ev.ID]; ok {
					continue
				}
				currentMet[ev.ID] = struct{}{}
				if _, ok := latestMet[ev.ID]; !ok {
					newEvents = append(newEvents, ev)
				}
			}
			// If first fetch or if there are no new missing events,
			// update etag and finish pagination loop.
			if etag == "" || len(newEvents) < (p.PerPage*p.Page) {
				if i == 0 {
					etag = res.Etag
				} else {
					lg.Info("Missing events resolved",
						zap.Int("pages", p.Page),
					)
				}
				break
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
			// Insert events in background.
		}

		// Calculating next sleep time to avoid rate limit.
		var targetRate time.Duration
		if rt.Remaining < 10 {
			lg.Warn("Rate limit", zap.Int("remaining", rt.Remaining))
			targetRate = rt.Reset.Sub(time.Now()) + time.Second
		} else {
			targetRate = time.Until(rt.Reset) / time.Duration(rt.Remaining)
		}
		duration := time.Since(start)
		sleep := targetRate - duration
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
			zap.Duration("sleep", sleep),
			zap.Duration("target_rate", targetRate),
		)
		select {
		case <-time.After(sleep):
			latestMet = currentMet
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
