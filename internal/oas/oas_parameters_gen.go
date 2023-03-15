// Code generated by ogen, DO NOT EDIT.

package oas

import (
	"net/http"

	"github.com/ogen-go/ogen/conv"
	"github.com/ogen-go/ogen/middleware"
	"github.com/ogen-go/ogen/ogenerrors"
	"github.com/ogen-go/ogen/uri"
	"github.com/ogen-go/ogen/validate"
)

// PollParams is parameters of poll operation.
type PollParams struct {
	// Worker token.
	XToken string
}

func unpackPollParams(packed middleware.Parameters) (params PollParams) {
	{
		key := middleware.ParameterKey{
			Name: "X-Token",
			In:   "header",
		}
		params.XToken = packed[key].(string)
	}
	return params
}

func decodePollParams(args [0]string, r *http.Request) (params PollParams, _ error) {
	h := uri.NewHeaderDecoder(r.Header)
	// Decode header: X-Token.
	if err := func() error {
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
				return err
			}
		} else {
			return validate.ErrFieldRequired
		}
		return nil
	}(); err != nil {
		return params, &ogenerrors.DecodeParamError{
			Name: "X-Token",
			In:   "header",
			Err:  err,
		}
	}
	return params, nil
}

// ProgressParams is parameters of progress operation.
type ProgressParams struct {
	// Worker token.
	XToken string
}

func unpackProgressParams(packed middleware.Parameters) (params ProgressParams) {
	{
		key := middleware.ParameterKey{
			Name: "X-Token",
			In:   "header",
		}
		params.XToken = packed[key].(string)
	}
	return params
}

func decodeProgressParams(args [0]string, r *http.Request) (params ProgressParams, _ error) {
	h := uri.NewHeaderDecoder(r.Header)
	// Decode header: X-Token.
	if err := func() error {
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
				return err
			}
		} else {
			return validate.ErrFieldRequired
		}
		return nil
	}(); err != nil {
		return params, &ogenerrors.DecodeParamError{
			Name: "X-Token",
			In:   "header",
			Err:  err,
		}
	}
	return params, nil
}
