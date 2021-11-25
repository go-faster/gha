// Code generated by ogen, DO NOT EDIT.

package oas

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"math"
	"net"
	"net/http"
	"net/url"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/go-faster/errors"
	"github.com/go-faster/jx"
	"github.com/google/uuid"
	"github.com/ogen-go/ogen/conv"
	ht "github.com/ogen-go/ogen/http"
	"github.com/ogen-go/ogen/json"
	"github.com/ogen-go/ogen/otelogen"
	"github.com/ogen-go/ogen/uri"
	"github.com/ogen-go/ogen/validate"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/trace"
)

// No-op definition for keeping imports.
var (
	_ = context.Background()
	_ = fmt.Stringer(nil)
	_ = strings.Builder{}
	_ = errors.Is
	_ = sort.Ints
	_ = http.MethodGet
	_ = io.Copy
	_ = json.Marshal
	_ = bytes.NewReader
	_ = strconv.ParseInt
	_ = time.Time{}
	_ = conv.ToInt32
	_ = uuid.UUID{}
	_ = uri.PathEncoder{}
	_ = url.URL{}
	_ = math.Mod
	_ = validate.Int{}
	_ = ht.NewRequest
	_ = net.IP{}
	_ = otelogen.Version
	_ = trace.TraceIDFromHex
	_ = otel.GetTracerProvider
	_ = metric.NewNoopMeterProvider
	_ = regexp.MustCompile
	_ = jx.Null
	_ = sync.Pool{}
)

// Encode implements json.Marshaler.
func (s Error) Encode(e *jx.Encoder) {
	e.ObjStart()

	e.FieldStart("message")
	e.Str(s.Message)
	e.ObjEnd()
}

// Decode decodes Error from json.
func (s *Error) Decode(d *jx.Decoder) error {
	if s == nil {
		return errors.New(`invalid: unable to decode Error to nil`)
	}
	return d.ObjBytes(func(d *jx.Decoder, k []byte) error {
		switch string(k) {
		case "message":
			v, err := d.Str()
			s.Message = string(v)
			if err != nil {
				return err
			}
		default:
			return d.Skip()
		}
		return nil
	})
}

// Encode implements json.Marshaler.
func (s ErrorStatusCode) Encode(e *jx.Encoder) {
	e.ObjStart()
	e.ObjEnd()
}

// Decode decodes ErrorStatusCode from json.
func (s *ErrorStatusCode) Decode(d *jx.Decoder) error {
	if s == nil {
		return errors.New(`invalid: unable to decode ErrorStatusCode to nil`)
	}
	return d.ObjBytes(func(d *jx.Decoder, k []byte) error {
		switch string(k) {
		default:
			return d.Skip()
		}
		return nil
	})
}

// Encode encodes Job as json.
func (s Job) Encode(e *jx.Encoder) {
	switch s.Type {
	case JobNothingJob:
		s.JobNothing.Encode(e)
	case JobDownloadJob:
		s.JobDownload.Encode(e)
	case JobInventoryJob:
		s.JobInventory.Encode(e)
	}
}

// Decode decodes Job from json.
func (s *Job) Decode(d *jx.Decoder) error {
	if s == nil {
		return errors.New(`invalid: unable to decode Job to nil`)
	}
	// Sum type discriminator.
	if d.Next() != jx.Object {
		return errors.Errorf("unexpected json type %q", d.Next())
	}
	var found bool
	if err := d.Capture(func(d *jx.Decoder) error {
		return d.ObjBytes(func(d *jx.Decoder, key []byte) error {
			if found {
				return d.Skip()
			}
			switch string(key) {
			case "type":
				typ, err := d.Str()
				if err != nil {
					return err
				}
				switch typ {
				case "download":
					s.Type = JobDownloadJob
					found = true
				case "inventory":
					s.Type = JobInventoryJob
					found = true
				case "nothing":
					s.Type = JobNothingJob
					found = true
				default:
					return errors.Errorf("unknown type %s", typ)
				}
				return nil
			}
			return d.Skip()
		})
	}); err != nil {
		return errors.Wrap(err, "capture")
	}
	if !found {
		return errors.New("unable to detect sum type variant")
	}
	switch s.Type {
	case JobNothingJob:
		if err := s.JobNothing.Decode(d); err != nil {
			return err
		}
	case JobDownloadJob:
		if err := s.JobDownload.Decode(d); err != nil {
			return err
		}
	case JobInventoryJob:
		if err := s.JobInventory.Decode(d); err != nil {
			return err
		}
	default:
		return errors.Errorf("inferred invalid type: %s", s.Type)
	}
	return nil
}

// Encode implements json.Marshaler.
func (s JobDownload) Encode(e *jx.Encoder) {
	e.ObjStart()

	e.FieldStart("type")
	e.Str(s.Type)

	e.FieldStart("date")
	e.Str(s.Date)
	e.ObjEnd()
}

// Decode decodes JobDownload from json.
func (s *JobDownload) Decode(d *jx.Decoder) error {
	if s == nil {
		return errors.New(`invalid: unable to decode JobDownload to nil`)
	}
	return d.ObjBytes(func(d *jx.Decoder, k []byte) error {
		switch string(k) {
		case "type":
			v, err := d.Str()
			s.Type = string(v)
			if err != nil {
				return err
			}
		case "date":
			v, err := d.Str()
			s.Date = string(v)
			if err != nil {
				return err
			}
		default:
			return d.Skip()
		}
		return nil
	})
}

// Encode implements json.Marshaler.
func (s JobInventory) Encode(e *jx.Encoder) {
	e.ObjStart()

	e.FieldStart("type")
	e.Str(s.Type)

	e.FieldStart("date")
	e.ArrStart()
	for _, elem := range s.Date {
		e.Str(elem)
	}
	e.ArrEnd()
	e.ObjEnd()
}

// Decode decodes JobInventory from json.
func (s *JobInventory) Decode(d *jx.Decoder) error {
	if s == nil {
		return errors.New(`invalid: unable to decode JobInventory to nil`)
	}
	return d.ObjBytes(func(d *jx.Decoder, k []byte) error {
		switch string(k) {
		case "type":
			v, err := d.Str()
			s.Type = string(v)
			if err != nil {
				return err
			}
		case "date":
			s.Date = nil
			if err := d.Arr(func(d *jx.Decoder) error {
				var elem string
				v, err := d.Str()
				elem = string(v)
				if err != nil {
					return err
				}
				s.Date = append(s.Date, elem)
				return nil
			}); err != nil {
				return err
			}
		default:
			return d.Skip()
		}
		return nil
	})
}

// Encode implements json.Marshaler.
func (s JobNothing) Encode(e *jx.Encoder) {
	e.ObjStart()

	e.FieldStart("type")
	e.Str(s.Type)
	e.ObjEnd()
}

// Decode decodes JobNothing from json.
func (s *JobNothing) Decode(d *jx.Decoder) error {
	if s == nil {
		return errors.New(`invalid: unable to decode JobNothing to nil`)
	}
	return d.ObjBytes(func(d *jx.Decoder, k []byte) error {
		switch string(k) {
		case "type":
			v, err := d.Str()
			s.Type = string(v)
			if err != nil {
				return err
			}
		default:
			return d.Skip()
		}
		return nil
	})
}

// Encode encodes int64 as json.
func (o OptInt64) Encode(e *jx.Encoder) {
	e.Int64(int64(o.Value))
}

// Decode decodes int64 from json.
func (o *OptInt64) Decode(d *jx.Decoder) error {
	if o == nil {
		return errors.New(`invalid: unable to decode OptInt64 to nil`)
	}
	switch d.Next() {
	case jx.Number:
		o.Set = true
		v, err := d.Int64()
		if err != nil {
			return err
		}
		o.Value = int64(v)
		return nil
	default:
		return errors.Errorf(`unexpected type %q while reading OptInt64`, d.Next())
	}
}

// Encode encodes string as json.
func (o OptString) Encode(e *jx.Encoder) {
	e.Str(string(o.Value))
}

// Decode decodes string from json.
func (o *OptString) Decode(d *jx.Decoder) error {
	if o == nil {
		return errors.New(`invalid: unable to decode OptString to nil`)
	}
	switch d.Next() {
	case jx.String:
		o.Set = true
		v, err := d.Str()
		if err != nil {
			return err
		}
		o.Value = string(v)
		return nil
	default:
		return errors.Errorf(`unexpected type %q while reading OptString`, d.Next())
	}
}

// Encode implements json.Marshaler.
func (s Progress) Encode(e *jx.Encoder) {
	e.ObjStart()

	e.FieldStart("event")
	s.Event.Encode(e)

	e.FieldStart("key")
	e.Str(s.Key)
	if s.InputSizeBytes.Set {
		e.FieldStart("input_size_bytes")
		s.InputSizeBytes.Encode(e)
	}
	if s.ContentSizeBytes.Set {
		e.FieldStart("content_size_bytes")
		s.ContentSizeBytes.Encode(e)
	}
	if s.OutputSizeBytes.Set {
		e.FieldStart("output_size_bytes")
		s.OutputSizeBytes.Encode(e)
	}
	if s.InputReadyBytes.Set {
		e.FieldStart("input_ready_bytes")
		s.InputReadyBytes.Encode(e)
	}
	if s.SHA256Input.Set {
		e.FieldStart("sha256_input")
		s.SHA256Input.Encode(e)
	}
	if s.SHA256Content.Set {
		e.FieldStart("sha256_content")
		s.SHA256Content.Encode(e)
	}
	if s.SHA256Output.Set {
		e.FieldStart("sha256_output")
		s.SHA256Output.Encode(e)
	}
	e.ObjEnd()
}

// Decode decodes Progress from json.
func (s *Progress) Decode(d *jx.Decoder) error {
	if s == nil {
		return errors.New(`invalid: unable to decode Progress to nil`)
	}
	return d.ObjBytes(func(d *jx.Decoder, k []byte) error {
		switch string(k) {
		case "event":
			if err := s.Event.Decode(d); err != nil {
				return err
			}
		case "key":
			v, err := d.Str()
			s.Key = string(v)
			if err != nil {
				return err
			}
		case "input_size_bytes":
			s.InputSizeBytes.Reset()
			if err := s.InputSizeBytes.Decode(d); err != nil {
				return err
			}
		case "content_size_bytes":
			s.ContentSizeBytes.Reset()
			if err := s.ContentSizeBytes.Decode(d); err != nil {
				return err
			}
		case "output_size_bytes":
			s.OutputSizeBytes.Reset()
			if err := s.OutputSizeBytes.Decode(d); err != nil {
				return err
			}
		case "input_ready_bytes":
			s.InputReadyBytes.Reset()
			if err := s.InputReadyBytes.Decode(d); err != nil {
				return err
			}
		case "sha256_input":
			s.SHA256Input.Reset()
			if err := s.SHA256Input.Decode(d); err != nil {
				return err
			}
		case "sha256_content":
			s.SHA256Content.Reset()
			if err := s.SHA256Content.Decode(d); err != nil {
				return err
			}
		case "sha256_output":
			s.SHA256Output.Reset()
			if err := s.SHA256Output.Decode(d); err != nil {
				return err
			}
		default:
			return d.Skip()
		}
		return nil
	})
}

// Encode encodes ProgressEvent as json.
func (s ProgressEvent) Encode(e *jx.Encoder) {
	e.Str(string(s))
}

// Decode decodes ProgressEvent from json.
func (s *ProgressEvent) Decode(d *jx.Decoder) error {
	if s == nil {
		return errors.New(`invalid: unable to decode ProgressEvent to nil`)
	}
	v, err := d.Str()
	if err != nil {
		return err
	}
	*s = ProgressEvent(v)
	return nil
}

// Encode implements json.Marshaler.
func (s Status) Encode(e *jx.Encoder) {
	e.ObjStart()

	e.FieldStart("message")
	e.Str(s.Message)
	e.ObjEnd()
}

// Decode decodes Status from json.
func (s *Status) Decode(d *jx.Decoder) error {
	if s == nil {
		return errors.New(`invalid: unable to decode Status to nil`)
	}
	return d.ObjBytes(func(d *jx.Decoder, k []byte) error {
		switch string(k) {
		case "message":
			v, err := d.Str()
			s.Message = string(v)
			if err != nil {
				return err
			}
		default:
			return d.Skip()
		}
		return nil
	})
}
