package entry

import (
	"bufio"
	"context"
	"io"

	"github.com/go-faster/errors"
	"github.com/go-faster/jx"
	"github.com/klauspost/compress/zstd"
)

type Reader struct {
	e Event
	j *jx.Decoder
	z *zstd.Decoder

	buf []byte
}

func (r *Reader) Close() {
	r.e.Reset()
	r.j = nil
	r.z.Close()
	r.buf = nil
}

func (r *Reader) Decode(ctx context.Context, rd io.Reader, f func(ctx context.Context, e *Event) error) error {
	if err := r.z.Reset(rd); err != nil {
		return errors.Wrap(err, "zstd reset")
	}

	s := bufio.NewScanner(r.z)
	s.Buffer(r.buf, len(r.buf))

	for s.Scan() {
		if err := ctx.Err(); err != nil {
			return err
		}
		r.j.ResetBytes(s.Bytes())
		r.e.Reset()
		if err := r.e.Decode(r.j); err != nil {
			continue
		}
		if err := f(ctx, &r.e); err != nil {
			return errors.Wrap(err, "f")
		}
	}

	return s.Err()
}

func NewReader() *Reader {
	z, err := zstd.NewReader(nil, zstd.WithDecoderConcurrency(1))
	if err != nil {
		panic(err)
	}
	return &Reader{
		buf: make([]byte, 1024*1024*10),
		z:   z,
		j:   jx.GetDecoder(),
	}
}
