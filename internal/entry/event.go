package entry

import (
	"time"

	"github.com/go-faster/errors"
	"github.com/go-faster/jx"
)

type Event struct {
	ID      uint64    `json:"id"`
	RepoID  uint64    `json:"repo_id"`
	Repo    string    `json:"repo"`
	ActorID uint64    `json:"actor_id"`
	Actor   string    `json:"actor"`
	Time    time.Time `json:"time"`
	Type    string    `json:"type"`
}

func (e *Event) Decode(d *jx.Decoder) error {
	return d.ObjBytes(func(d *jx.Decoder, key []byte) error {
		switch string(key) {
		case "type":
			v, err := d.Str()
			if err != nil {
				return errors.Wrap(err, "type")
			}
			e.Type = v
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
	})
}
