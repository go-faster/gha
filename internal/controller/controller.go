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

func (h Handler) Progress(ctx context.Context, req oas.Progress, params oas.ProgressParams) (oas.Status, error) {
	if err := h.authToken(ctx, params.XToken); err != nil {
		return oas.Status{}, err
	}

	h.lg.Info("Progress",
		zap.String("key", req.Key),
		zap.Bool("done", req.Done),
	)

	u := h.db.Chunk.Update().Where(
		chunk.IDEQ(req.Key),
	)

	if req.Done {
		if err := u.
			SetSha256Input(req.SHA256Input.Value).
			SetSha256Output(req.SHA256Output.Value).
			SetSha256Content(req.SHA256Content.Value).
			SetState(chunk.StateDownloaded).
			SetNillableLeaseExpiresAt(nil).
			Exec(ctx); err != nil {
			return oas.Status{}, errors.Wrap(err, "update")
		}
		return oas.Status{Message: "ack done"}, nil
	}

	if err := u.
		SetLeaseExpiresAt(time.Now().Add(time.Second * 15)).
		Exec(ctx); err != nil {
		return oas.Status{}, errors.Wrap(err, "update lease")
	}

	return oas.Status{Message: "ack"}, nil
}

func New(db *ent.Client, lg *zap.Logger) *Handler {
	return &Handler{
		lg: lg,
		db: db,
	}
}

var _ oas.Handler = (*Handler)(nil)

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

func (h Handler) authToken(ctx context.Context, tok string) error {
	_, err := h.db.Worker.Query().Where(worker.TokenEQ(tok)).Only(ctx)
	if ent.IsNotFound(err) {
		return &oas.ErrorStatusCode{
			StatusCode: 401,
			Response: oas.Error{
				Message: "Token not found",
			},
		}
	}
	if err != nil {
		return errors.Wrap(err, "token")
	}

	return nil
}

func (h Handler) Poll(ctx context.Context, params oas.PollParams) (oas.Job, error) {
	if err := h.authToken(ctx, params.XToken); err != nil {
		return oas.Job{}, err
	}

	tx, err := h.db.Tx(ctx)
	if err != nil {
		return oas.Job{}, errors.Wrap(err, "tx")
	}
	defer EnsureRollback(tx)()

	ch, err := tx.Chunk.Query().Where(
		chunk.StateIn(chunk.StateNew),
	).Limit(1).
		ForUpdate().
		First(ctx)
	if ent.IsNotFound(err) {
		return oas.NewJobNothingJob(oas.JobNothing{
			Type: "nothing",
		}), nil
	}
	if err != nil {
		return oas.Job{}, errors.Wrap(err, "query")
	}

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

	h.lg.Info("Scheduled job",
		zap.String("key", ch.ID),
	)

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
