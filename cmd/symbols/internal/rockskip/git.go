package rockskip

import (
	"context"

	"github.com/khulnasoft/khulnasoft/cmd/symbols/internal/fetcher"
	"github.com/khulnasoft/khulnasoft/internal/api"
	"github.com/khulnasoft/khulnasoft/internal/gitserver"
	"github.com/khulnasoft/khulnasoft/internal/search"
	"github.com/khulnasoft/khulnasoft/lib/errors"
)

type GitserverClient interface {
	ChangedFiles(context.Context, api.RepoName, string, string) (gitserver.ChangedFilesIterator, error)
	RevList(ctx context.Context, repo string, commit string, onCommit func(commit string) (shouldContinue bool, err error)) error
}

func archiveEach(ctx context.Context, fetcher fetcher.RepositoryFetcher, repo string, commit string, paths []string, onFile func(path string, contents []byte) error) error {
	if len(paths) == 0 {
		return nil
	}

	args := search.SymbolsParameters{Repo: api.RepoName(repo), CommitID: api.CommitID(commit)}
	parseRequestOrErrors := fetcher.FetchRepositoryArchive(ctx, args.Repo, args.CommitID, paths)
	defer func() {
		// Ensure the channel is drained
		for range parseRequestOrErrors {
		}
	}()

	for parseRequestOrError := range parseRequestOrErrors {
		if parseRequestOrError.Err != nil {
			return errors.Wrap(parseRequestOrError.Err, "FetchRepositoryArchive")
		}

		err := onFile(parseRequestOrError.ParseRequest.Path, parseRequestOrError.ParseRequest.Data)
		if err != nil {
			return err
		}
	}

	return nil
}
