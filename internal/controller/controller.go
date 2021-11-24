package controller

import (
	"context"
	"time"

	"github.com/go-faster/errors"
	"go.uber.org/zap"

	"github.com/go-faster/gha/internal/ent"
	"github.com/go-faster/gha/internal/ent/chunk"
	"github.com/go-faster/gha/internal/ent/worker"
	"github.com/go-faster/gha/internal/oas"
)

func EnsureRollback(tx *ent.Tx) func() {
	var done bool
	tx.OnRollback(func(rollbacker ent.Rollbacker) ent.Rollbacker {
		done = true
		return rollbacker
	})
	tx.OnCommit(func(committer ent.Committer) ent.Committer {
		done = true
		return committer
	})
	return func() {
		if done {
			return
		}
		_ = tx.Rollback()
	}
}

type Handler struct {
	lg *zap.Logger
	db *ent.Client
}

func New(db *ent.Client, lg *zap.Logger) *Handler {
	return &Handler{
		lg: lg,
		db: db,
	}
}

func (h Handler) Run(ctx context.Context) error {
	t := time.NewTicker(time.Second * 10)
	defer t.Stop()

	tick := func(now time.Time) error {
		if err := h.db.Chunk.Update().
			Where(
				chunk.StateEQ(chunk.StateDownloading),
				chunk.LeaseExpiresAtLT(now),
			).
			SetNillableLeaseExpiresAt(nil).
			SetState(chunk.StateNew).Exec(ctx); err != nil {
			return errors.Wrap(err, "update")
		}

		return nil
	}

	if err := tick(time.Now()); err != nil {
		return errors.Wrap(err, "initial tick")
	}

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case now := <-t.C:
			if err := tick(now); err != nil {
				return errors.Wrap(err, "tick")
			}
		}
	}

}

func (h Handler) Poll(ctx context.Context, params oas.PollParams) (oas.Job, error) {
	w, err := h.db.Worker.Query().Where(worker.TokenEQ(params.Token)).Only(ctx)
	if ent.IsNotFound(err) {
		return oas.Job{}, &oas.ErrorStatusCode{
			StatusCode: 401,
			Response: oas.Error{
				Message: "Token not found",
			},
		}
	}

	lg := h.lg.With(
		zap.String("worker_id", w.ID.String()),
	)
	lg.Info("Token found")

	tx, err := h.db.Tx(ctx)
	if err != nil {
		return oas.Job{}, errors.Wrap(err, "tx")
	}
	defer EnsureRollback(tx)()

	ch, err := tx.Chunk.Query().Where(
		chunk.StateIn(chunk.StateNew),
	).Limit(1).First(ctx)
	if err := ch.Update().
		SetLeaseExpiresAt(time.Now().Add(time.Second * 30)).
		SetState(chunk.StateDownloading).
		Exec(ctx); err != nil {
		return oas.Job{}, errors.Wrap(err, "lease")
	}
	if err != nil {
		return oas.Job{}, errors.Wrap(err, "find")
	}
	if err := tx.Commit(); err != nil {
		return oas.Job{}, errors.Wrap(err, "commit")
	}

	return oas.NewJobDownloadJob(oas.JobDownload{
		Type: "download",
		Date: ch.ID,
	}), err
}

func (h Handler) Status(ctx context.Context) (oas.Status, error) {
	return oas.Status{Message: "ok"}, nil
}

func (h Handler) NewError(ctx context.Context, err error) oas.ErrorStatusCode {
	return oas.ErrorStatusCode{
		StatusCode: 500,
		Response: oas.Error{
			Message: err.Error(),
		},
	}
}
