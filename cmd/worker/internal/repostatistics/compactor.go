package repostatistics

import (
	"context"
	"time"

	"github.com/sourcegraph/log"

	"github.com/khulnasoft/khulnasoft/cmd/worker/job"
	workerdb "github.com/khulnasoft/khulnasoft/cmd/worker/shared/init/db"
	"github.com/khulnasoft/khulnasoft/internal/database"
	"github.com/khulnasoft/khulnasoft/internal/env"
	"github.com/khulnasoft/khulnasoft/internal/goroutine"
	"github.com/khulnasoft/khulnasoft/internal/observation"
	"github.com/khulnasoft/khulnasoft/lib/errors"
)

// compactor is a worker responsible for compacting rows in the repo_statistics table.
type compactor struct{}

func NewCompactor() job.Job {
	return &compactor{}
}

func (j *compactor) Description() string {
	return ""
}

func (j *compactor) Config() []env.Config {
	return nil
}

func (j *compactor) Routines(_ context.Context, observationCtx *observation.Context) ([]goroutine.BackgroundRoutine, error) {
	db, err := workerdb.InitDB(observationCtx)
	if err != nil {
		return nil, err
	}

	return []goroutine.BackgroundRoutine{
		goroutine.NewPeriodicGoroutine(
			context.Background(),
			&compactorHandler{
				store:  db.RepoStatistics(),
				logger: observationCtx.Logger,
			},
			goroutine.WithName("repomgmt.statistics-compactor"),
			goroutine.WithDescription("compacts repo statistics"),
			goroutine.WithInterval(30*time.Minute),
		),
	}, nil
}

type compactorHandler struct {
	store  database.RepoStatisticsStore
	logger log.Logger
}

var _ goroutine.Handler = &compactorHandler{}

func (h *compactorHandler) Handle(ctx context.Context) error {
	var errs error
	if err := h.store.CompactRepoStatistics(ctx); err != nil {
		errs = errors.Append(errs, errors.Wrap(err, "error compacting repo statistics"))
	}

	if err := h.store.CompactGitserverReposStatistics(ctx); err != nil {
		errs = errors.Append(errs, errors.Wrap(err, "error compacting gitserver repos statistics"))
	}
	return errs
}
