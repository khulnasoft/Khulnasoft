package shared

import (
	"context"

	server "github.com/khulnasoft/khulnasoft/cmd/gitserver/internal"
	"github.com/khulnasoft/khulnasoft/cmd/gitserver/internal/common"
	"github.com/khulnasoft/khulnasoft/cmd/gitserver/internal/git"
	"github.com/khulnasoft/khulnasoft/cmd/gitserver/internal/git/gitcli"
	"github.com/khulnasoft/khulnasoft/cmd/gitserver/internal/gitserverfs"
	"github.com/khulnasoft/khulnasoft/internal/api"
	"github.com/khulnasoft/khulnasoft/internal/database"
	"github.com/khulnasoft/khulnasoft/internal/goroutine"
	"github.com/khulnasoft/khulnasoft/internal/observation"
	"github.com/khulnasoft/khulnasoft/internal/wrexec"
	"github.com/khulnasoft/khulnasoft/lib/errors"
)

// TestAPIServer returns a new gitserver API server for testing. Do not use this
// in a production workload.
func TestAPIServer(ctx context.Context, observationCtx *observation.Context, db database.DB, config *Config, getRemoteURLFunc func(ctx context.Context, repo api.RepoName) (string, error)) (goroutine.BackgroundRoutine, error) {
	logger := observationCtx.Logger

	// Load and validate configuration.
	if err := config.Validate(); err != nil {
		return nil, errors.Wrap(err, "failed to validate configuration")
	}

	// Prepare the file system.
	fs := gitserverfs.New(observationCtx, config.ReposDir)
	if err := fs.Initialize(); err != nil {
		return nil, err
	}

	backendSource := func(dir common.GitDir, repoName api.RepoName) git.GitBackend {
		return git.NewObservableBackend(gitcli.NewBackend(logger, wrexec.NewNoOpRecordingCommandFactory(), dir, repoName))
	}
	gitserver := makeServer(observationCtx, fs, db, wrexec.NewNoOpRecordingCommandFactory(), backendSource, config.ExternalAddress, config.CoursierCacheDir, server.NewRepositoryLocker(), getRemoteURLFunc)
	httpServer := makeHTTPServer(logger, fs, makeGRPCServer(logger, gitserver, config), config.ListenAddress)

	return &testServerRoutine{start: httpServer.Start, stop: func() {
		_ = httpServer.Stop(context.Background())
		gitserver.Stop()
	}}, nil
}

type testServerRoutine struct {
	start func()
	stop  func()
}

func (t *testServerRoutine) Name() string {
	return "gitserver-test"
}

func (t *testServerRoutine) Start() {
	t.start()
}

func (t *testServerRoutine) Stop(context.Context) error {
	t.stop()
	return nil
}
