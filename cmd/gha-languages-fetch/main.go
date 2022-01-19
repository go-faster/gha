package main

import (
	"bufio"
	"bytes"
	"context"
	"database/sql"
	"encoding/csv"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"sort"
	"strings"
	"time"

	_ "github.com/ClickHouse/clickhouse-go"
	"github.com/go-faster/errors"
	"github.com/go-faster/jx"
	"github.com/google/go-github/v40/github"
	"github.com/ogen-go/ogen/json"
	"go.etcd.io/bbolt"
	"go.uber.org/atomic"
	"go.uber.org/ratelimit"
	"go.uber.org/zap"
	"golang.org/x/oauth2"
	"golang.org/x/sync/errgroup"

	"github.com/go-faster/gha/internal/app"
	"github.com/go-faster/gha/internal/lang"
)

var ddl = fmt.Sprintf(`CREATE TABLE IF NOT EXISTS github.languages (
    name %s,
	repo String
) ENGINE MergeTree() ORDER BY (name, repo)`, lang.DDL())

func main() {
	app.Run(func(ctx context.Context, lg *zap.Logger) error {
		lg.Info("Starting")

		ts := oauth2.StaticTokenSource(
			&oauth2.Token{AccessToken: os.Getenv("GITHUB_TOKEN")},
		)
		tc := oauth2.NewClient(ctx, ts)
		c := github.NewClient(tc)

		rate := ratelimit.New(4800, ratelimit.Per(time.Hour))

		db, err := sql.Open("clickhouse", os.Getenv("CLICKHOUSE"))
		if err != nil {
			return errors.Wrap(err, "clickhouse")
		}

		cache, err := bbolt.Open("languages", 0666, &bbolt.Options{NoSync: true})
		if err != nil {
			return errors.Wrap(err, "db open")
		}
		defer func() {
			lg.Info("Closing")
			_ = cache.Sync()
			_ = cache.Close()
		}()

		var (
			nothing     = []byte{0}
			unavailable = []byte{1}
			bucket      = []byte("language")
		)
		var arg struct {
			Init bool
			Enum bool
		}

		flag.BoolVar(&arg.Init, "init", false, "init clickhouse language db")
		flag.BoolVar(&arg.Enum, "enum", false, "generate enum")
		flag.Parse()

		f, err := os.Open(flag.Arg(0))
		if err != nil {
			return errors.Wrap(err, "open")
		}
		defer func() {
			_ = f.Close()
		}()

		if arg.Enum {
			languages := map[string]bool{}
			if err := cache.View(func(t *bbolt.Tx) error {
				b := t.Bucket(bucket)
				if b == nil {
					return nil
				}
				return b.ForEach(func(k, v []byte) error {
					switch {
					case bytes.Equal(v, nothing):
						return nil
					case bytes.Equal(v, unavailable):
						return nil
					default:
						languages[string(v)] = true
					}
					return nil
				})
			}); err != nil {
				return errors.Wrap(err, "read")
			}

			var keys []string
			for k := range languages {
				keys = append(keys, k)
			}
			sort.Strings(keys)

			e := jx.GetEncoder()
			e.SetIdent(2)
			e.Obj(func(e *jx.Encoder) {
				for i, k := range keys {
					e.FieldStart(k)
					e.Int(i)
				}
			})

			e.WriteTo(os.Stdout)
			return nil
		}

		if arg.Init {
			if _, err := db.ExecContext(ctx, "DROP TABLE IF EXISTS github.languages"); err != nil {
				return errors.Wrap(err, "truncate")
			}
			if _, err := db.ExecContext(ctx, ddl); err != nil {
				return errors.Wrap(err, "create")
			}
			tx, err := db.BeginTx(ctx, &sql.TxOptions{})
			if err != nil {
				return errors.Wrap(err, "tx")
			}
			var count int
			stmt, err := tx.Prepare("INSERT INTO github.languages (repo, name) VALUES (?, ?)")
			if err != nil {
				return errors.Wrap(err, "prepare")
			}
			f, err := os.Open("bq-results-20211126-160210-t66zllshmlld.json")
			if err != nil {
				return err
			}
			defer func() {
				_ = f.Close()
			}()
			s := bufio.NewScanner(f)

			for s.Scan() {
				var v struct {
					Name     string `json:"repo_name"`
					Language []struct {
						Name  string `json:"name"`
						Bytes int    `json:"bytes,string"`
					}
				}
				if err := json.Unmarshal(s.Bytes(), &v); err != nil {
					return errors.Wrap(err, "json")
				}
				if len(v.Language) == 0 {
					continue
				}
				if !lang.Ok(v.Language[0].Name) {
					continue
				}
				count++
				if _, err := stmt.Exec(v.Name, v.Language[0].Name); err != nil {
					return errors.Wrap(err, "clickhouse")
				}
			}
			if err := s.Err(); err != nil {
				return err
			}

			if err := cache.View(func(t *bbolt.Tx) error {
				b := t.Bucket(bucket)
				if b == nil {
					return nil
				}
				return b.ForEach(func(k, v []byte) error {
					switch {
					case bytes.Equal(v, nothing):
						return nil
					case bytes.Equal(v, unavailable):
						return nil
					default:
						if !lang.Ok(string(v)) {
							return nil
						}
					}
					count++
					if _, err := stmt.Exec(string(k), string(v)); err != nil {
						return errors.Wrap(err, "clickhouse")
					}

					return nil
				})
			}); err != nil {
				return errors.Wrap(err, "ingest")
			}
			if err := tx.Commit(); err != nil {
				return errors.Wrap(err, "commit")
			}
			lg.Info("Initialized", zap.Int("count", count))
			return nil
		}

		r := csv.NewReader(bufio.NewReader(f))
		repos := make(chan []string)

		var count atomic.Int64
		g, ctx := errgroup.WithContext(ctx)
		for i := 0; i < 10; i++ {
			g.Go(func() error {
				for s := range repos {
					repo := s[0]
					if err := ctx.Err(); err != nil {
						return err
					}

					fmt.Println("done", count.Load())
					fmt.Println("repo", s[0], s[1])

					key := []byte(repo)

					var found bool

					if err := cache.View(func(tx *bbolt.Tx) error {
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
						count.Inc()
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
						if err := cache.Update(func(tx *bbolt.Tx) error {
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
						count.Inc()
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
					fmt.Println("set", repo, string(value))
					if err := set(value); err != nil {
						return err
					}

					count.Inc()
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
