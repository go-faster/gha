package entry

import (
	"bytes"
	"context"
	_ "embed"
	"testing"
)

//go:embed _testdata/2019-11-27T03.first100.json.zst
var dataCompressed []byte

func BenchmarkReader_Decode(b *testing.B) {
	var (
		r   = NewReader()
		br  = bytes.NewReader(dataCompressed)
		ctx = context.Background()
	)
	defer r.Close()
	b.SetBytes(int64(len(dataCompressed)))
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		br.Reset(dataCompressed)
		if err := r.Decode(ctx, br, func(ctx context.Context, e *Event) error {
			return nil
		}); err != nil {
			b.Fatal(err)
		}
	}
}
