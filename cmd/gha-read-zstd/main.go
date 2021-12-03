package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"os"
	"runtime/pprof"
	"time"

	"github.com/dustin/go-humanize"
	"github.com/go-faster/errors"
	"github.com/klauspost/compress/zstd"
	"go.uber.org/zap"

	"github.com/go-faster/gha/internal/app"
)

func summary(n int64, d time.Duration) string {
	seconds := d.Seconds()
	perSecond := float64(n) / seconds

	totalHuman := humanize.Bytes(uint64(n))
	speedHuman := humanize.Bytes(uint64(perSecond))

	return fmt.Sprintf("%s [%s/s]", totalHuman, speedHuman)
}

func main() {
	flag.Parse()
	app.Run(func(ctx context.Context, lg *zap.Logger) error {
		p, err := os.Create("read-zstd.cpu.out")
		if err != nil {
			return err
		}
		defer func() {
			_ = p.Close()
		}()
		if err := pprof.StartCPUProfile(p); err != nil {
			return errors.Wrap(err, "start cpu profile")
		}
		defer pprof.StopCPUProfile()

		start := time.Now()

		f, err := os.Open(flag.Arg(0))
		if err != nil {
			return errors.Wrap(err, "open")
		}
		defer func() {
			_ = f.Close()
		}()

		stat, err := f.Stat()
		if err != nil {
			return errors.Wrap(err, "stat")
		}

		r, err := zstd.NewReader(f)
		if err != nil {
			return errors.Wrap(err, "zstd")
		}
		defer r.Close()

		n, err := io.Copy(io.Discard, r)
		if err != nil {
			return errors.Wrap(err, "copy")
		}

		d := time.Since(start)

		fmt.Printf("%s (%s compressed) %s\n",
			summary(n, d), summary(stat.Size(), d), d.Round(time.Millisecond),
		)

		return nil
	})
}
