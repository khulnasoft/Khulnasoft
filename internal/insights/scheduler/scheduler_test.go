package scheduler

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	edb "github.com/khulnasoft/khulnasoft/internal/database"
	"github.com/khulnasoft/khulnasoft/internal/database/dbmocks"
	"github.com/khulnasoft/khulnasoft/internal/database/dbtest"
	"github.com/khulnasoft/khulnasoft/internal/goroutine"
	"github.com/khulnasoft/khulnasoft/internal/insights/priority"
	"github.com/khulnasoft/khulnasoft/internal/insights/store"
	"github.com/khulnasoft/khulnasoft/internal/insights/types"
	"github.com/khulnasoft/khulnasoft/internal/observation"
	"github.com/khulnasoft/khulnasoft/lib/errors"

	"github.com/sourcegraph/log/logtest"
)

func Test_MonitorStartsAndStops(t *testing.T) {
	logger := logtest.Scoped(t)
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	insightsDB := edb.NewInsightsDB(dbtest.NewInsightsDB(logger, t), logger)
	repos := dbmocks.NewMockRepoStore()
	config := JobMonitorConfig{
		InsightsDB:     insightsDB,
		RepoStore:      repos,
		ObservationCtx: observation.TestContextTB(t),
		CostAnalyzer:   priority.NewQueryAnalyzer(),
	}
	routines := NewBackgroundJobMonitor(ctx, config).Routines()
	err := goroutine.MonitorBackgroundRoutines(ctx, routines...)
	assert.EqualError(t, err, "unable to stop routines gracefully: context deadline exceeded")
}

func TestScheduler_InitialBackfill(t *testing.T) {
	logger := logtest.Scoped(t)
	ctx := context.Background()
	insightsDB := edb.NewInsightsDB(dbtest.NewInsightsDB(logger, t), logger)
	repos := dbmocks.NewMockRepoStore()
	insightsStore := store.NewInsightStore(insightsDB)
	config := JobMonitorConfig{
		InsightsDB:     insightsDB,
		RepoStore:      repos,
		ObservationCtx: observation.TestContextTB(t),
		CostAnalyzer:   priority.NewQueryAnalyzer(),
	}
	monitor := NewBackgroundJobMonitor(ctx, config)

	series, err := insightsStore.CreateSeries(ctx, types.InsightSeries{
		SeriesID:            "series1",
		Query:               "asdf",
		SampleIntervalUnit:  string(types.Month),
		SampleIntervalValue: 1,
		GenerationMethod:    types.Search,
	})
	require.NoError(t, err)

	scheduler := NewScheduler(insightsDB)
	backfill, err := scheduler.InitialBackfill(ctx, series)
	require.NoError(t, err)

	dequeue, found, err := monitor.newBackfillStore.Dequeue(ctx, "test", nil)
	require.NoError(t, err)
	if !found {
		t.Fatal(errors.New("no queued record found"))
	}
	require.Equal(t, backfill.Id, dequeue.backfillId)
}
