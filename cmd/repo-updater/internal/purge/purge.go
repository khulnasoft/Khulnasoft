package purge

import (
	"context"
	"math/rand"
	"time"

	"golang.org/x/time/rate"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/sourcegraph/log"

	"github.com/khulnasoft/khulnasoft/cmd/repo-updater/internal/gitserver"
	"github.com/khulnasoft/khulnasoft/internal/actor"
	"github.com/khulnasoft/khulnasoft/internal/conf/conftypes"
	"github.com/khulnasoft/khulnasoft/internal/database"
	"github.com/khulnasoft/khulnasoft/internal/goroutine"
	"github.com/khulnasoft/khulnasoft/internal/ratelimit"
	"github.com/khulnasoft/khulnasoft/lib/errors"
	"github.com/khulnasoft/khulnasoft/schema"
)

var (
	purgeSuccess = promauto.NewCounter(prometheus.CounterOpts{
		Name: "src_repoupdater_purge_success",
		Help: "Incremented each time we remove a repository clone.",
	})

	purgeFailed = promauto.NewCounter(prometheus.CounterOpts{
		Name: "src_repoupdater_purge_failed",
		Help: "Incremented each time we try and fail to remove a repository clone.",
	})
)

// NewRepositoryPurgeWorker is a worker which deletes repos which are present on
// gitserver, but not enabled/present in our repos table. ttl, should be >= 0 and
// specifies how long ago a repo must be deleted before it is purged.
func NewRepositoryPurgeWorker(ctx context.Context, logger log.Logger, db database.DB, conf conftypes.SiteConfigQuerier) goroutine.BackgroundRoutine {
	limiter := ratelimit.NewInstrumentedLimiter("PurgeRepoWorker", rate.NewLimiter(10, 1))

	var timeToNextPurge time.Duration

	return goroutine.NewPeriodicGoroutine(
		actor.WithInternalActor(ctx),
		goroutine.HandlerFunc(func(ctx context.Context) error {
			purgeConfig := conf.SiteConfig().RepoPurgeWorker
			if purgeConfig == nil {
				purgeConfig = &schema.RepoPurgeWorker{
					// Defaults - align with documentation
					IntervalMinutes:   15,
					DeletedTTLMinutes: 60,
				}
			}
			if purgeConfig.IntervalMinutes <= 0 {
				logger.Debug("purge worker disabled via site config", log.Int("repoPurgeWorker.interval", purgeConfig.IntervalMinutes))
				return nil
			}

			deletedBefore := time.Now().Add(-time.Duration(purgeConfig.DeletedTTLMinutes) * time.Minute)
			purgeLogger := logger.With(log.Time("deletedBefore", deletedBefore))

			timeToNextPurge = time.Duration(purgeConfig.IntervalMinutes) * time.Minute
			purgeLogger.Debug("running repository purge", log.Duration("timeToNextPurge", timeToNextPurge))
			if err := purge(ctx, purgeLogger, db, database.ListPurgableReposOptions{
				Limit:         5000,
				Limiter:       limiter,
				DeletedBefore: deletedBefore,
			}); err != nil {
				return errors.Wrap(err, "failed to run repository clone purge")
			}

			return nil
		}),
		goroutine.WithName("repo-updater.repo-purge-worker"),
		goroutine.WithDescription("deletes repos which are present on gitserver but not in the repos table"),
		goroutine.WithIntervalFunc(func() time.Duration {
			return randSleepDuration(timeToNextPurge, 1*time.Minute)
		}),
	)
}

// PurgeOldestRepos will start a go routine to purge the oldest repos limited by
// limit. The repos are ordered by when they were deleted. limit must be greater
// than zero.
func PurgeOldestRepos(logger log.Logger, db database.DB, limit int, perSecond float64) error {
	if limit <= 0 {
		return errors.Errorf("limit must be greater than zero, got %d", limit)
	}
	sglogError := log.Error

	go func() {
		limiter := ratelimit.NewInstrumentedLimiter("PurgeOldestRepos", rate.NewLimiter(rate.Limit(perSecond), 1))
		// Use a background routine so that we don't time out based on the http context.
		if err := purge(context.Background(), logger, db, database.ListPurgableReposOptions{
			Limit:   limit,
			Limiter: limiter,
		}); err != nil {
			logger.Error("Purging old repos", sglogError(err))
		}
	}()
	return nil
}

// purge purges repos, returning the number of repos that were successfully purged
func purge(ctx context.Context, logger log.Logger, db database.DB, options database.ListPurgableReposOptions) error {
	start := time.Now()
	gitserverClient := gitserver.NewRepositoryServiceClient("repoupdater.purgeworker")
	var (
		total   int
		success int
		failed  int
	)

	repos, err := db.GitserverRepos().ListPurgeableRepos(ctx, options)
	if err != nil {
		return errors.Wrap(err, "listing purgeable repos")
	}

	for _, repo := range repos {
		if options.Limiter != nil {
			if err := options.Limiter.Wait(ctx); err != nil {
				// A rate limit failure is fatal
				return errors.Wrap(err, "waiting for rate limiter")
			}
		}
		total++
		if err := gitserverClient.DeleteRepository(ctx, repo); err != nil {
			// Do not fail at this point, just log so we can remove other repos.
			logger.Error("failed to remove repository", log.String("repo", string(repo)), log.Error(err))
			purgeFailed.Inc()
			failed++
			continue
		}
		success++
		purgeSuccess.Inc()
	}

	// If we did something we log with a higher level.
	statusLogger := logger.Info
	if failed > 0 {
		statusLogger = logger.Warn
	}
	statusLogger("repository purge finished", log.Int("total", total), log.Int("removed", success), log.Int("failed", failed), log.Duration("duration", time.Since(start)))
	return nil
}

// randSleepDuration will sleep for an expected d duration with a jitter in [-jitter /
// 2, jitter / 2].
func randSleepDuration(d, jitter time.Duration) time.Duration {
	delta := time.Duration(rand.Int63n(int64(jitter))) - (jitter / 2)
	return d + delta
}
