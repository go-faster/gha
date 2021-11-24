package main

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/go-faster/errors"
	"go.uber.org/zap"
)

const layout = "2006-01-02-15"

type DateVar struct {
	Date *time.Time
}

func (d DateVar) String() string {
	if d.Date.IsZero() {
		return ""
	}
	return d.Date.Format(layout)
}

func (d DateVar) Set(s string) error {
	v, err := time.ParseInLocation(layout, s, time.UTC)
	if err != nil {
		return errors.Wrap(err, "parse time")
	}

	*d.Date = v
	return nil
}

func dateFlag(date *time.Time, name, usage string) {
	flag.Var(DateVar{Date: date}, name, usage)
}

func GetURL(date time.Time) string {
	return fmt.Sprintf("https://data.gharchive.org/%s.json.gz",
		date.Format(layout),
	)
}

func run(ctx context.Context) error {
	arg := struct {
		Date time.Time
	}{
		Date: time.Now().Add(-time.Hour * 6),
	}
	dateFlag(&arg.Date, "date", "date to download")
	flag.Parse()

	lg, err := zap.NewDevelopment()
	if err != nil {
		return errors.Wrap(err, "log")
	}

	link := GetURL(arg.Date)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, link, http.NoBody)
	if err != nil {
		return errors.Wrap(err, "req")
	}
	start := time.Now()
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return errors.Wrap(err, "get")
	}
	defer func() {
		_ = res.Body.Close()
	}()
	if res.StatusCode != http.StatusOK {
		return errors.Errorf("%s: bad code %d", link, res.StatusCode)
	}

	lg.Info("Found",
		zap.String("method", req.Method),
		zap.String("url", link),
		zap.Duration("duration", time.Since(start).Round(time.Millisecond)),
	)
	return nil
}

func main() {
	if err := run(context.Background()); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "%+v\n", err)
		os.Exit(2)
	}
}
