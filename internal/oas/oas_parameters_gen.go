// Code generated by ogen, DO NOT EDIT.

package oas

import (
	"net/http"

	"github.com/go-faster/errors"

	"github.com/ogen-go/ogen/conv"
	"github.com/ogen-go/ogen/uri"
)

type PollParams struct {
	// Worker token.
	XToken string
}

func decodePollParams(args [0]string, r *http.Request) (params PollParams, _ error) {
	h := uri.NewHeaderDecoder(r.Header)
	// Decode header: X-Token.
	{
		cfg := uri.HeaderParameterDecodingConfig{
			Name:    "X-Token",
			Explode: false,
		}
		if err := h.HasParam(cfg); err == nil {
			if err := h.DecodeParam(cfg, func(d uri.Decoder) error {
				val, err := d.DecodeValue()
				if err != nil {
					return err
				}

				c, err := conv.ToString(val)
				if err != nil {
					return err
				}

				params.XToken = c
				return nil
			}); err != nil {
				return params, errors.Wrap(err, "header: X-Token: parse")
			}
		} else {
			return params, errors.New("header: X-Token: not specified")
		}
	}
	return params, nil
}

type ProgressParams struct {
	// Worker token.
	XToken string
}

func decodeProgressParams(args [0]string, r *http.Request) (params ProgressParams, _ error) {
	h := uri.NewHeaderDecoder(r.Header)
	// Decode header: X-Token.
	{
		cfg := uri.HeaderParameterDecodingConfig{
			Name:    "X-Token",
			Explode: false,
		}
		if err := h.HasParam(cfg); err == nil {
			if err := h.DecodeParam(cfg, func(d uri.Decoder) error {
				val, err := d.DecodeValue()
				if err != nil {
					return err
				}

				c, err := conv.ToString(val)
				if err != nil {
					return err
				}

				params.XToken = c
				return nil
			}); err != nil {
				return params, errors.Wrap(err, "header: X-Token: parse")
			}
		} else {
			return params, errors.New("header: X-Token: not specified")
		}
	}
	return params, nil
}