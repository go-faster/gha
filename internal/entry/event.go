package entry

import (
	"bytes"
	"time"
	"unsafe"

	"github.com/go-faster/errors"
	"github.com/go-faster/jx"
)

type Event struct {
	Type []byte
	Repo []byte
	URL  []byte
	Time time.Time
}

func (e *Event) Reset() {
	e.Type = e.Type[:0]
	e.Repo = e.Repo[:0]
	e.Time = time.Time{}
}

func (e *Event) Decode(d *jx.Decoder) error {
	if err := d.ObjBytes(func(d *jx.Decoder, key []byte) error {
		switch string(key) {
		case "type":
			v, err := d.StrAppend(e.Type[:0])
			if err != nil {
				return errors.Wrap(err, "type")
			}
			e.Type = v
			return nil
		case "url":
			v, err := d.StrAppend(e.URL[:0])
			if err != nil {
				return errors.Wrap(err, "url")
			}
			e.URL = v
			return nil
		case "created_at":
			v, err := d.StrBytes()
			if err != nil {
				return errors.Wrap(err, "created at")
			}
			t, err := time.Parse(time.RFC3339, *(*string)(unsafe.Pointer(&v)))
			if err != nil {
				return errors.Wrap(err, "bytes")
			}
			e.Time = t
			return nil
		case "repo", "repository":
			if err := d.ObjBytes(func(d *jx.Decoder, key []byte) error {
				switch string(key) {
				case "url":
					v, err := d.StrAppend(e.Repo[:0])
					if err != nil {
						return errors.Wrap(err, "name")
					}
					e.Repo = v
					return nil
				case "full_name":
					v, err := d.StrAppend(e.Repo[:0])
					if err != nil {
						return errors.Wrap(err, "name")
					}
					e.Repo = v
					return nil
				case "name":
					if len(e.Repo) != 0 {
						return d.Skip()
					}
					v, err := d.StrAppend(e.Repo[:0])
					if err != nil {
						return errors.Wrap(err, "name")
					}
					e.Repo = v
					return nil
				default:
					if err := d.Skip(); err != nil {
						return errors.Wrap(err, "skip")
					}
					return nil
				}
			}); err != nil {
				return errors.Wrap(err, "repo")
			}
			return nil
		default:
			if err := d.Skip(); err != nil {
				return errors.Wrap(err, "skip")
			}
			return nil
		}
	}); err != nil {
		return errors.Wrap(err, "object")
	}

	if !e.Interesting() {
		return nil
	}

	if len(e.Repo) == 0 {
		e.Repo = bytes.TrimPrefix(e.URL, []byte("https://github.com/"))
	}
	if len(e.Repo) == 0 || e.Time.IsZero() || len(e.Type) == 0 {
		return errors.Errorf("missing data: %+v", e)
	}

	return nil
}

func (e Event) Interesting() bool {
	switch string(e.Type) {
	case "WatchEvent", "PushEvent", "IssuesEvent", "PullRequestEvent":
		return true
	default:
		return false
	}
}
