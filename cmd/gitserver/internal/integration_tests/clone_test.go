package inttests

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/khulnasoft/khulnasoft/lib/errors"

	mockassert "github.com/derision-test/go-mockgen/v2/testutil/assert"
	mockrequire "github.com/derision-test/go-mockgen/v2/testutil/require"
	"github.com/sourcegraph/log/logtest"
	"github.com/stretchr/testify/require"
	"golang.org/x/time/rate"

	server "github.com/khulnasoft/khulnasoft/cmd/gitserver/internal"
	common "github.com/khulnasoft/khulnasoft/cmd/gitserver/internal/common"
	"github.com/khulnasoft/khulnasoft/cmd/gitserver/internal/git"
	"github.com/khulnasoft/khulnasoft/cmd/gitserver/internal/git/gitcli"
	"github.com/khulnasoft/khulnasoft/cmd/gitserver/internal/gitserverfs"
	"github.com/khulnasoft/khulnasoft/cmd/gitserver/internal/vcssyncer"
	"github.com/khulnasoft/khulnasoft/internal/api"
	"github.com/khulnasoft/khulnasoft/internal/database/dbmocks"
	"github.com/khulnasoft/khulnasoft/internal/extsvc"
	"github.com/khulnasoft/khulnasoft/internal/observation"
	"github.com/khulnasoft/khulnasoft/internal/ratelimit"
	"github.com/khulnasoft/khulnasoft/internal/types"
	"github.com/khulnasoft/khulnasoft/internal/vcs"
	"github.com/khulnasoft/khulnasoft/internal/wrexec"
)

func TestClone(t *testing.T) {
	root := t.TempDir()
	reposDir := filepath.Join(root, "repos")
	remote := createSimpleGitRepo(t, root)

	logger := logtest.Scoped(t)
	db := newMockDB()
	gsStore := dbmocks.NewMockGitserverRepoStore()
	db.GitserverReposFunc.SetDefaultReturn(gsStore)
	ctx := context.Background()
	repo := api.RepoName("github.com/test/repo")

	locker := NewMockRepositoryLocker()
	lock := NewMockRepositoryLock()
	locker.TryAcquireFunc.SetDefaultReturn(lock, true)

	fs := gitserverfs.New(observation.TestContextTB(t), reposDir)
	require.NoError(t, fs.Initialize())
	getRemoteURLFunc := func(_ context.Context, name api.RepoName) (string, error) { //nolint:unparam
		require.Equal(t, repo, name)
		return remote, nil
	}

	s := server.NewServer(&server.ServerOpts{
		Logger: logger,
		FS:     fs,
		GitBackendSource: func(dir common.GitDir, repoName api.RepoName) git.GitBackend {
			return gitcli.NewBackend(logtest.Scoped(t), wrexec.NewNoOpRecordingCommandFactory(), dir, repoName)
		},
		GetRemoteURLFunc: getRemoteURLFunc,
		GetVCSSyncer: func(ctx context.Context, name api.RepoName) (vcssyncer.VCSSyncer, error) {
			getRemoteURLSource := func(ctx context.Context, name api.RepoName) (vcssyncer.RemoteURLSource, error) {
				require.Equal(t, repo, name)
				return vcssyncer.RemoteURLSourceFunc(func(ctx context.Context) (*vcs.URL, error) {
					raw, err := getRemoteURLFunc(ctx, name)
					if err != nil {
						return nil, errors.Wrapf(err, "failed to get remote URL for %s", name)
					}

					u, err := vcs.ParseURL(raw)
					if err != nil {
						return nil, errors.Wrapf(err, "failed to parse remote URL %q", raw)
					}
					return u, nil
				}), nil
			}

			require.Equal(t, repo, name)
			return vcssyncer.NewGitRepoSyncer(logtest.Scoped(t), wrexec.NewNoOpRecordingCommandFactory(), getRemoteURLSource), nil
		},
		DB:                      db,
		RecordingCommandFactory: wrexec.NewNoOpRecordingCommandFactory(),
		Locker:                  locker,
		RPSLimiter:              ratelimit.NewInstrumentedLimiter("GitserverTest", rate.NewLimiter(100, 10)),
		Hostname:                "test-shard",
	})

	// Requesting a repo update should figure out that the repo is not yet
	// cloned and call clone. We expect that clone to succeed.
	_, _, err := s.FetchRepository(ctx, repo)
	require.NoError(t, err)

	// Should have acquired a lock.
	mockassert.CalledOnce(t, locker.TryAcquireFunc)
	// Should have reported status. 24 lines is the output git currently produces.
	// This number might need to be adjusted over time, but before doing so please
	// check that the calls actually use the args you would expect them to use.
	mockassert.CalledN(t, lock.SetStatusFunc, 24)
	// Should have released the lock.
	mockassert.CalledOnce(t, lock.ReleaseFunc)

	// Check it was set to cloning first, then cloned.
	mockassert.CalledN(t, gsStore.SetCloneStatusFunc, 2)
	mockassert.CalledWith(t, gsStore.SetCloneStatusFunc, mockassert.Values(mockassert.Skip, repo, types.CloneStatusCloning, "test-shard"))
	mockassert.CalledWith(t, gsStore.SetCloneStatusFunc, mockassert.Values(mockassert.Skip, repo, types.CloneStatusCloned, "test-shard"))

	// Last output should have been stored for the repo.
	mockrequire.CalledOnce(t, gsStore.SetLastOutputFunc)
	haveLastOutput := gsStore.SetLastOutputFunc.History()[0].Arg2
	require.Contains(t, haveLastOutput, "Creating bare repo\n")
	require.Contains(t, haveLastOutput, "Created bare repo at ")
	require.Contains(t, haveLastOutput, "Fetching remote contents\n")
	// Ensure the path is properly redacted. The redactor just takes the whole
	// remote URL as redacted so this is expected.
	require.Contains(t, haveLastOutput, "From <redacted>")
	// Double newlines should not be part of our standard output. They are not
	// forbidden, but we currently don't use them. So this will guard against
	// regressions in the log processing to make sure we don't accidentally
	// introduce blank newlines for CRLF parsing. (Yes this took a while to get
	// right).
	require.NotContains(t, haveLastOutput, "\n\n")

	// Check that it was called exactly once total.
	mockassert.CalledOnce(t, gsStore.SetLastErrorFunc)
	// And that it was called for the right repo, setting the last error to empty.
	mockassert.CalledWith(t, gsStore.SetLastErrorFunc, mockassert.Values(mockassert.Skip, repo, "", "test-shard"))

	// Check that the repo is in the expected location on disk.
	_, err = os.Stat(fs.RepoDir(repo).Path())
	require.NoError(t, err)
}

func TestClone_Fail(t *testing.T) {
	root := t.TempDir()
	reposDir := filepath.Join(root, "repos")
	remote := filepath.Join(root, "non-existing")

	logger := logtest.Scoped(t)
	db := newMockDB()
	gsStore := dbmocks.NewMockGitserverRepoStore()
	db.GitserverReposFunc.SetDefaultReturn(gsStore)
	ctx := context.Background()
	repo := api.RepoName("github.com/test/repo")

	locker := NewMockRepositoryLocker()
	lock := NewMockRepositoryLock()
	locker.TryAcquireFunc.SetDefaultReturn(lock, true)

	fs := gitserverfs.New(observation.TestContextTB(t), reposDir)
	require.NoError(t, fs.Initialize())
	getRemoteURLFunc := func(_ context.Context, name api.RepoName) (string, error) { //nolint:unparam
		require.Equal(t, repo, name)
		return remote, nil
	}

	s := server.NewServer(&server.ServerOpts{
		Logger: logger,
		FS:     fs,
		GitBackendSource: func(dir common.GitDir, repoName api.RepoName) git.GitBackend {
			return gitcli.NewBackend(logtest.Scoped(t), wrexec.NewNoOpRecordingCommandFactory(), dir, repoName)
		},
		GetRemoteURLFunc: getRemoteURLFunc,
		GetVCSSyncer: func(ctx context.Context, name api.RepoName) (vcssyncer.VCSSyncer, error) {
			require.Equal(t, repo, name)
			getRemoteURLSource := func(ctx context.Context, name api.RepoName) (vcssyncer.RemoteURLSource, error) {
				require.Equal(t, repo, name)
				return vcssyncer.RemoteURLSourceFunc(func(ctx context.Context) (*vcs.URL, error) {
					raw, err := getRemoteURLFunc(ctx, name)
					if err != nil {
						return nil, errors.Wrapf(err, "failed to get remote URL for %s", name)
					}

					u, err := vcs.ParseURL(raw)
					if err != nil {
						return nil, errors.Wrapf(err, "failed to parse remote URL %q", raw)
					}
					return u, nil
				}), nil
			}
			return vcssyncer.NewGitRepoSyncer(logtest.Scoped(t), wrexec.NewNoOpRecordingCommandFactory(), getRemoteURLSource), nil
		},
		DB:                      db,
		RecordingCommandFactory: wrexec.NewNoOpRecordingCommandFactory(),
		Locker:                  locker,
		RPSLimiter:              ratelimit.NewInstrumentedLimiter("GitserverTest", rate.NewLimiter(100, 10)),
		Hostname:                "test-shard",
	})

	// Requesting a repo update should figure out that the repo is not yet cloned and call clone.
	// We expect that clone to fail.
	_, _, err := s.FetchRepository(ctx, repo)
	require.Error(t, err)
	require.Contains(t, err.Error(), "failed to clone github.com/test/repo: clone failed. Output: Creating bare repo\nCreated bare repo at")

	// Should have acquired a lock.
	mockassert.CalledN(t, locker.TryAcquireFunc, 1)
	// Should have reported status. 7 lines is the output git currently produces.
	// This number might need to be adjusted over time, but before doing so please
	// check that the calls actually use the args you would expect them to use.
	mockassert.CalledN(t, lock.SetStatusFunc, 7)
	// Should have released the lock.
	mockassert.CalledN(t, lock.ReleaseFunc, 1)

	// Check it was set to cloning first, then uncloned again (since clone failed).
	mockassert.CalledN(t, gsStore.SetCloneStatusFunc, 2)
	mockassert.CalledWith(t, gsStore.SetCloneStatusFunc, mockassert.Values(mockassert.Skip, repo, types.CloneStatusCloning, "test-shard"))
	mockassert.CalledWith(t, gsStore.SetCloneStatusFunc, mockassert.Values(mockassert.Skip, repo, types.CloneStatusNotCloned, "test-shard"))

	// Last output should have been stored for the repo.
	mockrequire.CalledOnce(t, gsStore.SetLastOutputFunc)
	haveLastOutput := gsStore.SetLastOutputFunc.History()[0].Arg2
	require.Contains(t, haveLastOutput, "Creating bare repo\n")
	require.Contains(t, haveLastOutput, "Created bare repo at ")
	require.Contains(t, haveLastOutput, "Fetching remote contents\n")
	// Check that also git output made it here.
	require.Contains(t, haveLastOutput, "does not appear to be a git repository\n")

	// Check that it was called exactly once total.
	mockrequire.CalledOnce(t, gsStore.SetLastErrorFunc)
	// And that it was called for the right repo, setting the last error to empty.
	mockassert.CalledWith(t, gsStore.SetLastErrorFunc, mockassert.Values(mockassert.Skip, repo, mockassert.Skip, "test-shard"))
	require.Contains(t, gsStore.SetLastErrorFunc.History()[0].Arg2, "Creating bare repo\n")
	require.Contains(t, gsStore.SetLastErrorFunc.History()[0].Arg2, "failed to fetch: exit status 128")

	// Check that no repo is in the expected location on disk.
	_, err = os.Stat(fs.RepoDir(repo).Path())
	require.Error(t, err)
	require.True(t, os.IsNotExist(err))
}

func newMockDB() *dbmocks.MockDB {
	db := dbmocks.NewMockDB()
	db.GitserverReposFunc.SetDefaultReturn(dbmocks.NewMockGitserverRepoStore())
	db.FeatureFlagsFunc.SetDefaultReturn(dbmocks.NewMockFeatureFlagStore())

	r := dbmocks.NewMockRepoStore()
	r.GetByNameFunc.SetDefaultHook(func(ctx context.Context, repoName api.RepoName) (*types.Repo, error) {
		return &types.Repo{
			Name: repoName,
			ExternalRepo: api.ExternalRepoSpec{
				ServiceType: extsvc.TypeGitHub,
			},
		}, nil
	})
	db.ReposFunc.SetDefaultReturn(r)

	return db
}
