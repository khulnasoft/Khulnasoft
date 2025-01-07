package backend

import (
	"context"
	"strings"
	"testing"

	"github.com/sourcegraph/log/logtest"

	"github.com/khulnasoft/khulnasoft/internal/api"
	"github.com/khulnasoft/khulnasoft/internal/database/dbmocks"
	"github.com/khulnasoft/khulnasoft/internal/gitserver"
	"github.com/khulnasoft/khulnasoft/lib/errors"
)

func TestRepos_ResolveRev_noRevSpecified_getsDefaultBranch(t *testing.T) {
	logger := logtest.Scoped(t)
	ctx := testContext()

	const wantRepo = "a"
	want := strings.Repeat("a", 40)

	client := gitserver.NewMockClient()
	var calledVCSRepoResolveRevision bool
	client.ResolveRevisionFunc.SetDefaultHook(func(context.Context, api.RepoName, string, gitserver.ResolveRevisionOptions) (api.CommitID, error) {
		calledVCSRepoResolveRevision = true
		return api.CommitID(want), nil
	})

	// (no rev/branch specified)
	commitID, err := NewRepos(logger, dbmocks.NewMockDB(), client).ResolveRev(ctx, "a", "")
	if err != nil {
		t.Fatal(err)
	}
	if !calledVCSRepoResolveRevision {
		t.Error("!calledVCSRepoResolveRevision")
	}
	if string(commitID) != want {
		t.Errorf("got resolved commit %q, want %q", commitID, want)
	}
}

func TestRepos_ResolveRev_noCommitIDSpecified_resolvesRev(t *testing.T) {
	ctx := testContext()
	logger := logtest.Scoped(t)

	const wantRepo = "a"
	want := strings.Repeat("a", 40)

	var calledVCSRepoResolveRevision bool
	client := gitserver.NewMockClient()
	client.ResolveRevisionFunc.SetDefaultHook(func(context.Context, api.RepoName, string, gitserver.ResolveRevisionOptions) (api.CommitID, error) {
		calledVCSRepoResolveRevision = true
		return api.CommitID(want), nil
	})

	commitID, err := NewRepos(logger, dbmocks.NewMockDB(), client).ResolveRev(ctx, "a", "b")
	if err != nil {
		t.Fatal(err)
	}
	if !calledVCSRepoResolveRevision {
		t.Error("!calledVCSRepoResolveRevision")
	}
	if string(commitID) != want {
		t.Errorf("got resolved commit %q, want %q", commitID, want)
	}
}

func TestRepos_ResolveRev_commitIDSpecified_resolvesCommitID(t *testing.T) {
	ctx := testContext()
	logger := logtest.Scoped(t)

	const wantRepo = "a"
	want := strings.Repeat("a", 40)

	var calledVCSRepoResolveRevision bool
	client := gitserver.NewMockClient()
	client.ResolveRevisionFunc.SetDefaultHook(func(context.Context, api.RepoName, string, gitserver.ResolveRevisionOptions) (api.CommitID, error) {
		calledVCSRepoResolveRevision = true
		return api.CommitID(want), nil
	})

	commitID, err := NewRepos(logger, dbmocks.NewMockDB(), client).ResolveRev(ctx, "a", strings.Repeat("a", 40))
	if err != nil {
		t.Fatal(err)
	}
	if !calledVCSRepoResolveRevision {
		t.Error("!calledVCSRepoResolveRevision")
	}
	if string(commitID) != want {
		t.Errorf("got resolved commit %q, want %q", commitID, want)
	}
}

func TestRepos_ResolveRev_commitIDSpecified_failsToResolve(t *testing.T) {
	ctx := testContext()
	logger := logtest.Scoped(t)

	const wantRepo = "a"
	want := errors.New("x")

	var calledVCSRepoResolveRevision bool
	client := gitserver.NewMockClient()
	client.ResolveRevisionFunc.SetDefaultHook(func(context.Context, api.RepoName, string, gitserver.ResolveRevisionOptions) (api.CommitID, error) {
		calledVCSRepoResolveRevision = true
		return "", errors.New("x")
	})

	_, err := NewRepos(logger, dbmocks.NewMockDB(), client).ResolveRev(ctx, "a", strings.Repeat("a", 40))
	if !errors.Is(err, want) {
		t.Fatalf("got err %v, want %v", err, want)
	}
	if !calledVCSRepoResolveRevision {
		t.Error("!calledVCSRepoResolveRevision")
	}
}
