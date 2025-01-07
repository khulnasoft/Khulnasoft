package auth

import (
	"context"

	"github.com/khulnasoft/khulnasoft/internal/extsvc/github"
	"github.com/khulnasoft/khulnasoft/internal/types"
)

type GitHubClient interface {
	GetRepository(ctx context.Context, owner string, name string) (*github.Repository, error)
	ListInstallationRepositories(ctx context.Context, page int) ([]*github.Repository, bool, int, error)
}

type UserStore interface {
	GetByCurrentAuthUser(context.Context) (*types.User, error)
}
