package main

import (
	"context"
	"fmt"
	"os"
	"strings"

	_ "github.com/ClickHouse/clickhouse-go"
	"github.com/go-faster/jx"
	"go.uber.org/zap"

	"github.com/go-faster/gha/internal/app"
)

func main() {
	app.Run(func(ctx context.Context, lg *zap.Logger) error {
		d := jx.GetDecoder()
		d.Reset(os.Stdin)

		fmt.Print("Enum16(")

		if err := d.Obj(func(d *jx.Decoder, k string) error {
			v, err := d.Int()
			if err != nil {
				return err
			}

			if v != 0 {
				fmt.Print(", ")
			}

			k = strings.ReplaceAll(k, "'", "\\'")
			k = fmt.Sprintf("'%s'", k)
			fmt.Printf("%s = %d", k, v)

			return nil
		}); err != nil {
			return err
		}

		fmt.Println(")")

		return nil
	})
}
