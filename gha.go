package main

import (
	"bufio"
	"context"
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"encoding/hex"
	"flag"
	"fmt"
	"hash"
	"io"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"time"

	"github.com/go-faster/errors"
	"github.com/klauspost/compress/gzip"
	"github.com/klauspost/compress/zstd"
	"go.uber.org/zap"
)

// adjusted original layout to be a prefix of RFC3339.
const layout = "2006-01-02T15"

type DateVar struct {
	Date *time.Time
}

func (d DateVar) String() string {
	if d.Date == nil || d.Date.IsZero() {
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
	// You can't describe 24-hour without leading zero.
	// Correct format is yyyy-MM-dd-H.
	return fmt.Sprintf("https://data.gharchive.org/%s-%d.json.gz",
		date.Format("2006-01-02"), date.Hour(),
	)
}

type Hash struct {
	io.Writer

	SHA1   hash.Hash
	SHA256 hash.Hash
	MD5    hash.Hash
}

func hexHashField(name string, h hash.Hash) zap.Field {
	return zap.String(name, hex.EncodeToString(h.Sum(nil)))
}

func (h Hash) Fields() []zap.Field {
	return []zap.Field{
		hexHashField("sha256", h.SHA256),
		hexHashField("sha1", h.SHA1),
		hexHashField("md5", h.MD5),
	}
}

func NewHash() Hash {
	h := Hash{
		SHA1:   sha1.New(),
		SHA256: sha256.New(),
		MD5:    md5.New(),
	}
	h.Writer = io.MultiWriter(
		h.SHA1,
		h.SHA256,
		h.MD5,
	)
	return h
}

func run(ctx context.Context) (err error) {
	arg := struct {
		Date time.Time
		Dir  string
	}{
		// Lags ~4 hours from realtime.
		Date: time.Now().Add(-time.Hour * 6),
	}
	dateFlag(&arg.Date, "date", "date to download")
	flag.StringVar(&arg.Dir, "dir", filepath.Join(os.TempDir(), "github-archive"), "directory to use")
	flag.Parse()

	lg, err := zap.NewDevelopment()
	if err != nil {
		return errors.Wrap(err, "log")
	}

	if err := os.MkdirAll(arg.Dir, 0755); err != nil {
		return errors.Wrap(err, "mkdir")
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

	var (
		hInput  = NewHash() // hash for input
		hData   = NewHash() // hash for json
		hOutput = NewHash() // hash for output
	)

	// Reading body as gzip.
	bodyBuf := bufio.NewReader(res.Body)
	reader, err := gzip.NewReader(io.TeeReader(bodyBuf, hInput))
	if err != nil {
		return errors.Wrap(err, "gzip")
	}

	// Output. Change to file.
	outName := fmt.Sprintf("%s.json.zst", arg.Date.Format(layout))
	outPath := filepath.Join(arg.Dir, outName)
	out, err := os.Create(outPath)
	if err != nil {
		return errors.Wrap(err, "create file")
	}
	defer func() {
		_ = out.Close()
	}()
	defer func() {
		if err == nil {
			// Cleaning up only on failure.
			return
		}

		lg.Warn("Cleaning up file")
		if err := os.Remove(outPath); err != nil {
			lg.Error("Failed to cleanup", zap.Error(err))
		}
	}()

	outWriter, err := zstd.NewWriter(io.MultiWriter(out, hOutput))
	if err != nil {
		return errors.Wrap(err, "zstd writer init")
	}

	// Up to 1 MiB for copy.
	buf := make([]byte, 1024*1024)
	total, err := io.CopyBuffer(io.MultiWriter(outWriter, hData), reader, buf)
	if err != nil {
		return errors.Wrap(err, "copy")
	}

	if err := outWriter.Close(); err != nil {
		return errors.Wrap(err, "zstd")
	}

	// Calculating re-compression ratio for fun and profit.
	stat, err := out.Stat()
	if err != nil {
		return errors.Wrap(err, "stat")
	}
	ratio := float64(res.ContentLength) / float64(stat.Size())
	totalRatio := float64(total) / float64(stat.Size())

	if err := out.Close(); err != nil {
		return errors.Wrap(err, "close file")
	}

	lg.Info("Processed",
		zap.String("path", outPath),
		zap.String("date", arg.Date.Format(layout)),
		zap.Int64("bytes_output", stat.Size()),
		zap.Int64("bytes_total", total),
		zap.Int64("bytes_input", res.ContentLength),
		zap.String("relative_ratio", fmt.Sprintf("%.0f%%", ratio*100)),
		zap.String("absolute_ratio", fmt.Sprintf("%.0f%%", totalRatio*100)),
		zap.Duration("duration", time.Since(start).Round(time.Millisecond)),
	)

	lg.Info("Input", hInput.Fields()...)
	lg.Info("Data", hData.Fields()...)
	lg.Info("Output", hOutput.Fields()...)

	return nil
}

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	if err := run(ctx); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Failed: %+v\n", err)
		os.Exit(2)
	}
}
