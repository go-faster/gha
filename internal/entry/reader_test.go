package entry

import (
	"bytes"
	"context"
	_ "embed"
	"testing"
	"time"
)

//go:embed _testdata/2019-11-27T03.first100.json.zst
var dataCompressed []byte

func BenchmarkReader_Decode(b *testing.B) {
	var (
		ctx = context.Background()
		nop = func(ctx context.Context, e *Event) error { return nil }
	)
	b.SetBytes(int64(len(dataCompressed)))
	b.ReportAllocs()
	b.ResetTimer()

	start := time.Now()

	b.RunParallel(func(pb *testing.PB) {
		var (
			r  = NewReader()
			br = bytes.NewReader(dataCompressed)
		)
		defer r.Close()

		for pb.Next() {
			br.Reset(dataCompressed)
			if err := r.Decode(ctx, br, nop); err != nil {
				b.Fatal(err)
			}
		}
	})

	duration := time.Since(start)
	uncompressed := float64(b.N*len(dataMulti)) / duration.Seconds()
	uncompressed /= 1024 * 1024 // MB
	b.ReportMetric(uncompressed, "(json)MB/s")
}
