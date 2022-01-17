package entry

import (
	"bufio"
	"bytes"
	_ "embed"
	"testing"

	"github.com/go-faster/jx"
)

var (
	//go:embed _testdata/2019-11-27T03.single.json
	dataSingle []byte

	//go:embed _testdata/2019-11-27T03.first100.json
	dataMulti []byte
)

func BenchmarkEvent_Decode(b *testing.B) {
	b.Run("Single", func(b *testing.B) {
		d := jx.GetDecoder()
		b.ReportAllocs()
		b.SetBytes(int64(len(dataSingle)))

		var e Event

		for i := 0; i < b.N; i++ {
			e.Reset()
			d.ResetBytes(dataSingle)
			if err := e.Decode(d); err != nil {
				b.Fatal(err)
			}
		}
	})
	b.Run("Concurrent", func(b *testing.B) {
		b.ReportAllocs()
		b.SetBytes(int64(len(dataSingle)))
		b.RunParallel(func(pb *testing.PB) {
			d := jx.GetDecoder()
			var e Event
			for pb.Next() {
				e.Reset()
				d.ResetBytes(dataSingle)
				if err := e.Decode(d); err != nil {
					b.Fatal(err)
				}
			}
		})
	})
	b.Run("Multi", func(b *testing.B) {
		d := jx.GetDecoder()
		b.ReportAllocs()
		b.SetBytes(int64(len(dataMulti)))
		r := bytes.NewReader(dataMulti)
		buf := make([]byte, 1024*1024)

		var e Event

		for i := 0; i < b.N; i++ {
			r.Reset(dataMulti)
			s := bufio.NewScanner(r)
			s.Buffer(buf, len(buf))

			for s.Scan() {
				e.Reset()
				d.ResetBytes(s.Bytes())
				if err := e.Decode(d); err != nil {
					b.Fatal()
				}
			}
			if err := s.Err(); err != nil {
				b.Fatal(err)
			}
		}
	})
	b.Run("MultiConcurrent", func(b *testing.B) {
		b.ReportAllocs()
		b.SetBytes(int64(len(dataMulti)))

		b.RunParallel(func(pb *testing.PB) {
			d := jx.GetDecoder()
			r := bytes.NewReader(dataMulti)
			buf := make([]byte, 1024*1024)

			for pb.Next() {
				r.Reset(dataMulti)
				s := bufio.NewScanner(r)
				s.Buffer(buf, len(buf))

				var e Event

				for s.Scan() {
					e.Reset()
					d.ResetBytes(s.Bytes())
					if err := e.Decode(d); err != nil {
						b.Fatal()
					}
				}
				if err := s.Err(); err != nil {
					b.Fatal(err)
				}
			}
		})
	})
}
