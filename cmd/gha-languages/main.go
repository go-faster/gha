package main

import (
	"bufio"
	"encoding/csv"
	"os"

	"github.com/go-faster/errors"
	"github.com/go-faster/jx"
)

func main() {
	s := bufio.NewScanner(os.Stdin)
	d := jx.GetDecoder()
	w := csv.NewWriter(os.Stdout)

	var (
		repo, language string
		languageBytes  int64
	)
	_ = w.Write([]string{"repository", "language"})
	for s.Scan() {
		d.ResetBytes(s.Bytes())

		repo = ""
		language = ""
		languageBytes = 0

		if err := d.ObjBytes(func(d *jx.Decoder, key []byte) error {
			switch string(key) {
			case "repo_name":
				s, err := d.Str()
				if err != nil {
					return err
				}
				repo = s
				return nil
			case "language":
				return d.Arr(func(d *jx.Decoder) error {
					var (
						n string
						b int64
					)

					if err := d.ObjBytes(func(d *jx.Decoder, key []byte) error {
						switch string(key) {
						case "name":
							s, err := d.Str()
							if err != nil {
								return err
							}
							n = s
							return nil
						case "bytes":
							num, err := d.Num()
							if err != nil {
								return err
							}
							v, err := num.Int64()
							if err != nil {
								return err
							}
							b = v
							return nil
						default:
							return errors.Errorf("k: %q", key)
						}
					}); err != nil {
						return err
					}

					if b > languageBytes {
						languageBytes = b
						language = n
					}

					return nil
				})
			default:
				return errors.Errorf("k: %q", key)
			}
		}); err != nil {
			panic(err)
		}

		if err := w.Write([]string{repo, language}); err != nil {
			panic(err)
		}
	}

	w.Flush()
	if err := w.Error(); err != nil {
		panic(err)
	}
}
