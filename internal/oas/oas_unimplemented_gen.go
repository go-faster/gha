// Code generated by ogen, DO NOT EDIT.

package oas

import (
	"context"

	ht "github.com/ogen-go/ogen/http"
)

var _ Handler = UnimplementedHandler{}

// UnimplementedHandler is no-op Handler which returns http.ErrNotImplemented.
type UnimplementedHandler struct{}

// Poll implements poll operation.
//
// Request job from coordinator.
//
// POST /job/poll
func (UnimplementedHandler) Poll(ctx context.Context, params PollParams) (r Job, _ error) {
	return r, ht.ErrNotImplemented
}

// Progress implements progress operation.
//
// Report progress.
//
// POST /progress
func (UnimplementedHandler) Progress(ctx context.Context, req Progress, params ProgressParams) (r Status, _ error) {
	return r, ht.ErrNotImplemented
}

// Status implements status operation.
//
// Get status.
//
// GET /status
func (UnimplementedHandler) Status(ctx context.Context) (r Status, _ error) {
	return r, ht.ErrNotImplemented
}

// NewError creates ErrorStatusCode from error returned by handler.
//
// Used for common default response.
func (UnimplementedHandler) NewError(ctx context.Context, err error) (r ErrorStatusCode) {
	return r
}
