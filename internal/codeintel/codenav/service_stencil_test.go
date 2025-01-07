package codenav

import (
	"context"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/sourcegraph/log"

	lsifstoremocks "github.com/khulnasoft/khulnasoft/internal/codeintel/codenav/internal/lsifstore/mocks"
	"github.com/khulnasoft/khulnasoft/internal/codeintel/codenav/shared"
	uploadsshared "github.com/khulnasoft/khulnasoft/internal/codeintel/uploads/shared"
	"github.com/khulnasoft/khulnasoft/internal/gitserver"
	"github.com/khulnasoft/khulnasoft/internal/observation"
	"github.com/khulnasoft/khulnasoft/internal/search/client"
	sgtypes "github.com/khulnasoft/khulnasoft/internal/types"
)

func TestStencil(t *testing.T) {
	// Set up mocks
	fakeRepoStore := AllPresentFakeRepoStore{}
	mockLsifStore := lsifstoremocks.NewMockLsifStore()
	mockUploadSvc := NewMockUploadService()
	mockGitserverClient := gitserver.NewMockClient()
	mockSearchClient := client.NewMockSearchClient()

	mockLsifStore.FindDocumentIDsFunc.SetDefaultHook(findDocumentIDsFuncAllowAny())

	// Init service
	svc := newService(observation.TestContextTB(t), fakeRepoStore, mockLsifStore, mockUploadSvc, mockGitserverClient, mockSearchClient, log.NoOp())

	// Set up request state
	mockRequestState := RequestState{}
	mockRequestState.SetLocalCommitCache(fakeRepoStore, mockGitserverClient)
	mockRequestState.SetLocalGitTreeTranslator(mockGitserverClient, &sgtypes.Repo{})
	uploads := []uploadsshared.CompletedUpload{
		{ID: 50, Commit: "deadbeef", Root: "sub1/"},
		{ID: 51, Commit: "deadbeef", Root: "sub2/"},
		{ID: 52, Commit: "deadbeef", Root: "sub3/"},
		{ID: 53, Commit: "deadbeef", Root: "sub4/"},
	}
	mockRequestState.SetUploadsDataLoader(uploads)

	expectedRanges := []shared.Range{
		{Start: shared.Position{Line: 10, Character: 20}, End: shared.Position{Line: 10, Character: 30}},
		{Start: shared.Position{Line: 11, Character: 20}, End: shared.Position{Line: 11, Character: 30}},
		{Start: shared.Position{Line: 12, Character: 20}, End: shared.Position{Line: 12, Character: 30}},
		{Start: shared.Position{Line: 13, Character: 20}, End: shared.Position{Line: 13, Character: 30}},
		{Start: shared.Position{Line: 14, Character: 20}, End: shared.Position{Line: 14, Character: 30}},
		{Start: shared.Position{Line: 15, Character: 20}, End: shared.Position{Line: 15, Character: 30}},
		{Start: shared.Position{Line: 16, Character: 20}, End: shared.Position{Line: 16, Character: 30}},
		{Start: shared.Position{Line: 17, Character: 20}, End: shared.Position{Line: 17, Character: 30}},
		{Start: shared.Position{Line: 18, Character: 20}, End: shared.Position{Line: 18, Character: 30}},
		{Start: shared.Position{Line: 19, Character: 20}, End: shared.Position{Line: 19, Character: 30}},
	}
	mockLsifStore.GetStencilFunc.PushReturn(nil, nil)
	mockLsifStore.GetStencilFunc.PushReturn(expectedRanges, nil)

	mockRequest := PositionalRequestArgs{
		RequestArgs: RequestArgs{
			RepositoryID: 42,
			Commit:       mockCommit,
			Limit:        50,
		},
		Path:      mockPath,
		Line:      10,
		Character: 20,
	}
	ranges, err := svc.GetStencil(context.Background(), mockRequest, mockRequestState)
	if err != nil {
		t.Fatalf("unexpected error querying hover: %s", err)
	}

	if diff := cmp.Diff(expectedRanges, ranges); diff != "" {
		t.Errorf("unexpected range (-want +got):\n%s", diff)
	}
}

func TestStencilWithDuplicateRanges(t *testing.T) {
	// Set up mocks
	fakeRepoStore := AllPresentFakeRepoStore{}
	mockLsifStore := lsifstoremocks.NewMockLsifStore()
	mockUploadSvc := NewMockUploadService()
	mockGitserverClient := gitserver.NewMockClient()
	mockSearchClient := client.NewMockSearchClient()

	mockLsifStore.FindDocumentIDsFunc.SetDefaultHook(findDocumentIDsFuncAllowAny())

	// Init service
	svc := newService(observation.TestContextTB(t), fakeRepoStore, mockLsifStore, mockUploadSvc, mockGitserverClient, mockSearchClient, log.NoOp())

	// Set up request state
	mockRequestState := RequestState{}
	mockRequestState.SetLocalCommitCache(fakeRepoStore, mockGitserverClient)
	mockRequestState.SetLocalGitTreeTranslator(mockGitserverClient, &sgtypes.Repo{})
	uploads := []uploadsshared.CompletedUpload{
		{ID: 50, Commit: "deadbeef", Root: "sub1/"},
		{ID: 51, Commit: "deadbeef", Root: "sub2/"},
		{ID: 52, Commit: "deadbeef", Root: "sub3/"},
		{ID: 53, Commit: "deadbeef", Root: "sub4/"},
	}
	mockRequestState.SetUploadsDataLoader(uploads)

	expectedRanges := []shared.Range{
		{Start: shared.Position{Line: 10, Character: 20}, End: shared.Position{Line: 10, Character: 30}},
		{Start: shared.Position{Line: 11, Character: 20}, End: shared.Position{Line: 11, Character: 30}},
		{Start: shared.Position{Line: 12, Character: 20}, End: shared.Position{Line: 12, Character: 30}},
		{Start: shared.Position{Line: 13, Character: 20}, End: shared.Position{Line: 13, Character: 30}},
		{Start: shared.Position{Line: 14, Character: 20}, End: shared.Position{Line: 14, Character: 30}},
		{Start: shared.Position{Line: 15, Character: 20}, End: shared.Position{Line: 15, Character: 30}},
		{Start: shared.Position{Line: 16, Character: 20}, End: shared.Position{Line: 16, Character: 30}},
		{Start: shared.Position{Line: 17, Character: 20}, End: shared.Position{Line: 17, Character: 30}},
		{Start: shared.Position{Line: 18, Character: 20}, End: shared.Position{Line: 18, Character: 30}},
		{Start: shared.Position{Line: 19, Character: 20}, End: shared.Position{Line: 19, Character: 30}},
	}
	mockLsifStore.GetStencilFunc.PushReturn(nil, nil)

	// Duplicate the ranges to test that we dedupe them
	mockLsifStore.GetStencilFunc.PushReturn(append(expectedRanges, expectedRanges...), nil)

	mockRequest := PositionalRequestArgs{
		RequestArgs: RequestArgs{
			RepositoryID: 42,
			Commit:       mockCommit,
			Limit:        50,
		},
		Path:      mockPath,
		Line:      10,
		Character: 20,
	}
	ranges, err := svc.GetStencil(context.Background(), mockRequest, mockRequestState)
	if err != nil {
		t.Fatalf("unexpected error querying hover: %s", err)
	}

	if diff := cmp.Diff(expectedRanges, ranges); diff != "" {
		t.Errorf("unexpected range (-want +got):\n%s", diff)
	}
}
