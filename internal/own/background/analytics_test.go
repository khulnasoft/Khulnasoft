package background

import (
	"bytes"
	"context"
	"io"
	"io/fs"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/khulnasoft/khulnasoft/internal/fileutil"
	"github.com/khulnasoft/khulnasoft/internal/rcache"

	"github.com/khulnasoft/khulnasoft/internal/api"
	"github.com/khulnasoft/khulnasoft/internal/authz"
	"github.com/khulnasoft/khulnasoft/internal/database"
	"github.com/khulnasoft/khulnasoft/internal/database/dbtest"
	"github.com/khulnasoft/khulnasoft/internal/gitserver"
	"github.com/khulnasoft/khulnasoft/internal/observation"
	"github.com/khulnasoft/khulnasoft/internal/types"
)

type fakeGitServer struct {
	gitserver.Client
	files        []string
	fileContents map[string]string
}

func (f fakeGitServer) ReadDir(ctx context.Context, repo api.RepoName, commit api.CommitID, path string, recursive bool) (gitserver.ReadDirIterator, error) {
	fis := make([]fs.FileInfo, 0, len(f.files))
	for _, file := range f.files {
		fis = append(fis, &fileutil.FileInfo{
			Name_: file,
		})
	}
	return gitserver.NewReadDirIteratorFromSlice(fis), nil
}

func (f fakeGitServer) ResolveRevision(ctx context.Context, repo api.RepoName, spec string, opt gitserver.ResolveRevisionOptions) (api.CommitID, error) {
	return api.CommitID(""), nil
}

func (f fakeGitServer) NewFileReader(ctx context.Context, repo api.RepoName, commit api.CommitID, name string) (io.ReadCloser, error) {
	if f.fileContents == nil {
		return nil, os.ErrNotExist
	}
	contents, ok := f.fileContents[name]
	if !ok {
		return nil, os.ErrNotExist
	}
	return io.NopCloser(bytes.NewReader([]byte(contents))), nil
}

func TestAnalyticsIndexerSuccess(t *testing.T) {
	rcache.SetupForTest(t)
	obsCtx := observation.TestContextTB(t)
	logger := obsCtx.Logger
	db := database.NewDB(logger, dbtest.NewDB(t))
	ctx := context.Background()
	user, err := db.Users().Create(ctx, database.NewUser{Username: "test"})
	require.NoError(t, err)
	var repoID api.RepoID = 1
	require.NoError(t, db.Repos().Create(ctx, &types.Repo{Name: "repo", ID: repoID}))
	client := fakeGitServer{
		files: []string{
			"notOwned.go",
			"alsoNotOwned.go",
			"owned/file1.go",
			"owned/file2.go",
			"owned/file3.go",
			"assigned.go",
		},
		fileContents: map[string]string{
			"CODEOWNERS": "/owned/* @owner",
		},
	}
	checker := authz.NewMockSubRepoPermissionChecker()
	checker.EnabledFunc.SetDefaultReturn(true)
	checker.EnabledForRepoIDFunc.SetDefaultReturn(false, nil)
	require.NoError(t, db.AssignedOwners().Insert(ctx, user.ID, repoID, "owned/file1.go", user.ID))
	require.NoError(t, db.AssignedOwners().Insert(ctx, user.ID, repoID, "assigned.go", user.ID))
	require.NoError(t, newAnalyticsIndexer(client, db, logger).indexRepo(ctx, repoID, checker))

	totalFileCount, err := db.RepoPaths().AggregateFileCount(ctx, database.TreeLocationOpts{})
	require.NoError(t, err)
	assert.Equal(t, int32(len(client.files)), totalFileCount)

	gotCounts, err := db.OwnershipStats().QueryAggregateCounts(ctx, database.TreeLocationOpts{})
	require.NoError(t, err)
	// We don't really need to compare time here.
	defaultTime := time.Time{}
	gotCounts.UpdatedAt = defaultTime
	wantCounts := database.PathAggregateCounts{
		CodeownedFileCount:         3,
		AssignedOwnershipFileCount: 2,
		TotalOwnedFileCount:        4,
		UpdatedAt:                  defaultTime,
	}
	assert.Equal(t, wantCounts, gotCounts)
}

func TestAnalyticsIndexerSkipsReposWithSubRepoPerms(t *testing.T) {
	rcache.SetupForTest(t)
	obsCtx := observation.TestContextTB(t)
	logger := obsCtx.Logger
	db := database.NewDB(logger, dbtest.NewDB(t))
	ctx := context.Background()
	var repoID api.RepoID = 1
	err := db.Repos().Create(ctx, &types.Repo{Name: "repo", ID: repoID})
	require.NoError(t, err)
	client := fakeGitServer{
		files: []string{"notOwned.go", "alsoNotOwned.go", "owned/file1.go", "owned/file2.go", "owned/file3.go"},
		fileContents: map[string]string{
			"CODEOWNERS": "/owned/* @owner",
		},
	}
	checker := authz.NewMockSubRepoPermissionChecker()
	checker.EnabledFunc.SetDefaultReturn(true)
	checker.EnabledForRepoIDFunc.SetDefaultReturn(true, nil)
	err = newAnalyticsIndexer(client, db, logger).indexRepo(ctx, repoID, checker)
	require.NoError(t, err)

	totalFileCount, err := db.RepoPaths().AggregateFileCount(ctx, database.TreeLocationOpts{})
	require.NoError(t, err)
	assert.Equal(t, int32(0), totalFileCount)

	codeownedCount, err := db.OwnershipStats().QueryAggregateCounts(ctx, database.TreeLocationOpts{})
	require.NoError(t, err)
	assert.Equal(t, database.PathAggregateCounts{CodeownedFileCount: 0}, codeownedCount)
}

func TestAnalyticsIndexerNoCodeowners(t *testing.T) {
	rcache.SetupForTest(t)
	obsCtx := observation.TestContextTB(t)
	logger := obsCtx.Logger
	db := database.NewDB(logger, dbtest.NewDB(t))
	ctx := context.Background()
	var repoID api.RepoID = 1
	err := db.Repos().Create(ctx, &types.Repo{Name: "repo", ID: repoID})
	require.NoError(t, err)
	client := fakeGitServer{
		files: []string{"notOwned.go", "alsoNotOwned.go", "owned/file1.go", "owned/file2.go", "owned/file3.go"},
	}
	checker := authz.NewMockSubRepoPermissionChecker()
	checker.EnabledFunc.SetDefaultReturn(true)
	checker.EnabledForRepoIDFunc.SetDefaultReturn(false, nil)
	err = newAnalyticsIndexer(client, db, logger).indexRepo(ctx, repoID, checker)
	require.NoError(t, err)

	totalFileCount, err := db.RepoPaths().AggregateFileCount(ctx, database.TreeLocationOpts{})
	require.NoError(t, err)
	assert.Equal(t, int32(5), totalFileCount)

	codeownedCount, err := db.OwnershipStats().QueryAggregateCounts(ctx, database.TreeLocationOpts{})
	defaultTime := time.Time{}
	codeownedCount.UpdatedAt = defaultTime
	require.NoError(t, err)
	assert.Equal(t, database.PathAggregateCounts{CodeownedFileCount: 0, UpdatedAt: defaultTime}, codeownedCount)
}
