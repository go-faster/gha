package main

import (
	"bufio"
	"bytes"
	"context"
	"encoding/csv"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/go-faster/errors"
	"github.com/google/go-github/v40/github"
	"go.etcd.io/bbolt"
	"go.uber.org/ratelimit"
	"go.uber.org/zap"
	"golang.org/x/oauth2"
	"golang.org/x/sync/errgroup"

	"github.com/go-faster/gha/internal/app"
)

func main() {
	app.Run(func(ctx context.Context, lg *zap.Logger) error {
		lg.Info("Starting")

		ts := oauth2.StaticTokenSource(
			&oauth2.Token{AccessToken: os.Getenv("GITHUB_TOKEN")},
		)
		tc := oauth2.NewClient(ctx, ts)
		c := github.NewClient(tc)

		rate := ratelimit.New(4800, ratelimit.Per(time.Hour))

		db, err := bbolt.Open("languages", 0666, &bbolt.Options{NoSync: true})
		if err != nil {
			return errors.Wrap(err, "db open")
		}
		defer func() {
			lg.Info("Closing")
			_ = db.Sync()
			_ = db.Close()
		}()

		flag.Parse()
		f, err := os.Open(flag.Arg(0))
		if err != nil {
			return errors.Wrap(err, "open")
		}
		defer func() {
			_ = f.Close()
		}()

		r := csv.NewReader(bufio.NewReader(f))
		repos := make(chan []string)
		bucket := []byte("language")

		var (
			nothing     = []byte{0}
			unavailable = []byte{1}

			count int
		)

		g, ctx := errgroup.WithContext(ctx)
		for i := 0; i < 10; i++ {
			g.Go(func() error {
				for s := range repos {
					repo := s[0]
					if err := ctx.Err(); err != nil {
						return err
					}

					fmt.Println("done", count)
					fmt.Println("repo", s[0], s[1])

					key := []byte(repo)

					var found bool

					if err := db.View(func(tx *bbolt.Tx) error {
						b := tx.Bucket(bucket)
						if b == nil {
							return nil
						}
						v := b.Get(key)
						found = v != nil

						if !found {
							return nil
						}

						switch {
						case bytes.Equal(v, nothing):
							fmt.Println(repo, "-")
						case bytes.Equal(v, unavailable):
							fmt.Println(repo, "N/A")
						default:
							fmt.Println(repo, string(v))
						}

						return nil
					}); err != nil {
						return err
					}

					if found {
						count++
						continue
					}

					v := strings.Split(repo, "/")
					rate.Take()
					languages, res, err := c.Repositories.ListLanguages(ctx, v[0], v[1])
					if re, ok := err.(*github.RateLimitError); ok {
						d := time.Until(re.Rate.Reset.Time) + time.Second*10
						fmt.Println("sleeping", d)
						select {
						case <-ctx.Done():
							return ctx.Err()
						case <-time.After(d):
							continue
						}
					}
					if res == nil {
						if err == nil {
							err = errors.New("nil body")
						}
						return err
					}

					set := func(v []byte) error {
						if err := db.Update(func(tx *bbolt.Tx) error {
							b, err := tx.CreateBucketIfNotExists(bucket)
							if err != nil {
								return errors.Wrap(err, "bucket")
							}

							if err := b.Put(key, v); err != nil {
								return errors.Wrap(err, "put")
							}

							return nil
						}); err != nil {
							return errors.Wrap(err, "update")
						}
						return nil
					}

					switch res.StatusCode {
					case http.StatusNotFound, http.StatusForbidden, 451:
						if err := set(unavailable); err != nil {
							return err
						}
						count++
						continue
					case http.StatusOK:
					default:
						time.Sleep(time.Second)
						continue
					}
					if err != nil {
						return err
					}

					value := nothing
					var max int
					for l, n := range languages {
						if n > max {
							max = n
							value = []byte(l)
						}
					}
					if err := set(value); err != nil {
						return err
					}

					count++
				}

				return nil
			})
		}

		g.Go(func() error {
			defer close(repos)

			for {
				if err := ctx.Err(); err != nil {
					return err
				}
				rec, err := r.Read()
				if errors.Is(err, io.EOF) {
					return nil
				}
				if err != nil {
					return err
				}
				if len(rec) != 2 {
					return errors.New("bad format")
				}
				select {
				case <-ctx.Done():
					return ctx.Err()
				case repos <- rec:
				}
			}
		})

		return g.Wait()
	})
}
