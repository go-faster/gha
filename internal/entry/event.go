package entry

import (
	"path"
	"strings"
	"time"

	"github.com/go-faster/errors"
	"github.com/go-faster/jx"
)

type Event struct {
	Actor string
	Type  string
	Repo  string
	Time  time.Time

	ID      uint64 // optional
	RepoID  uint64 // optional
	ActorID uint64 // optional
	URL     string // optional
}

func (e *Event) Decode(d *jx.Decoder) error {
	if err := d.ObjBytes(func(d *jx.Decoder, key []byte) error {
		switch string(key) {
		case "type":
			v, err := d.Str()
			if err != nil {
				return errors.Wrap(err, "type")
			}
			e.Type = v
			return nil
		case "url":
			v, err := d.Str()
			if err != nil {
				return errors.Wrap(err, "url")
			}
			e.URL = v
			return nil
		case "id":
			v, err := d.Num()
			if err != nil {
				return errors.Wrap(err, "id")
			}
			id, err := v.Uint64()
			if err != nil {
				return errors.Wrap(err, "uint64")
			}
			e.ID = id
			return nil
		case "created_at":
			v, err := d.StrBytes()
			if err != nil {
				return errors.Wrap(err, "created at")
			}
			t, err := time.Parse(time.RFC3339, string(v))
			if err != nil {
				return errors.Wrap(err, "bytes")
			}
			e.Time = t
			return nil
		case "repo", "repository":
			if err := d.ObjBytes(func(d *jx.Decoder, key []byte) error {
				switch string(key) {
				case "id":
					v, err := d.Uint64()
					if err != nil {
						return errors.Wrap(err, "id")
					}
					e.RepoID = v
					return nil
				case "url":
					v, err := d.Str()
					if err != nil {
						return errors.Wrap(err, "name")
					}
					e.Repo = path.Join(path.Base(path.Dir(v)), path.Base(v))
					return nil
				case "full_name":
					v, err := d.Str()
					if err != nil {
						return errors.Wrap(err, "name")
					}
					e.Repo = v
					return nil
				case "name":
					v, err := d.Str()
					if err != nil {
						return errors.Wrap(err, "name")
					}
					if e.Repo == "" {
						e.Repo = v
					}
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
		case "actor":
			switch d.Next() {
			case jx.String:
				v, err := d.Str()
				if err != nil {
					return errors.Wrap(err, "actor")
				}
				e.Actor = v
				return nil
			case jx.Object:
				if err := d.ObjBytes(func(d *jx.Decoder, key []byte) error {
					switch string(key) {
					case "id":
						v, err := d.Uint64()
						if err != nil {
							return errors.Wrap(err, "id")
						}
						e.ActorID = v
						return nil
					case "login":
						v, err := d.Str()
						if err != nil {
							return errors.Wrap(err, "name")
						}
						e.Actor = v
						return nil
					default:
						if err := d.Skip(); err != nil {
							return errors.Wrap(err, "skip")
						}
						return nil
					}
				}); err != nil {
					return errors.Wrap(err, "skip")
				}
				return nil
			default:
				if err := d.Skip(); err != nil {
					return errors.Wrap(err, "skip")
				}
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

	if e.Repo == "" {
		s := strings.TrimPrefix(e.URL, "https://github.com/")
		elems := strings.Split(s, "/")
		e.Repo = strings.Join(elems[:2], "/")
	}

	if e.Repo == "" || e.Time.IsZero() || e.Type == "" || e.Actor == "" {
		return errors.Errorf("missing data: %+v", e)
	}

	return nil
}

func (e Event) Interesting() bool {
	switch e.Type {
	case "WatchEvent", "PushEvent", "IssuesEvent", "PullRequestEvent":
		return true
	default:
		return false
	}
}
