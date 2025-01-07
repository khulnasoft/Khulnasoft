package gitlab

import (
	"net/url"

	"github.com/khulnasoft/khulnasoft/internal/authz"
	atypes "github.com/khulnasoft/khulnasoft/internal/authz/types"
	"github.com/khulnasoft/khulnasoft/internal/database"
	"github.com/khulnasoft/khulnasoft/internal/extsvc"
	"github.com/khulnasoft/khulnasoft/internal/extsvc/gitlab"
	"github.com/khulnasoft/khulnasoft/internal/licensing"
	"github.com/khulnasoft/khulnasoft/internal/types"
	"github.com/khulnasoft/khulnasoft/lib/errors"
	"github.com/khulnasoft/khulnasoft/schema"
)

// NewAuthzProviders returns the set of GitLab authz providers derived from the connections.
//
// It also returns any simple validation problems with the config, separating these into "serious problems"
// and "warnings". "Serious problems" are those that should make Sourcegraph set authz.allowAccessByDefault
// to false. "Warnings" are all other validation problems.
//
// This constructor does not and should not directly check connectivity to external services - if
// desired, callers should use `(*Provider).ValidateConnection` directly to get warnings related
// to connection issues.
func NewAuthzProviders(
	db database.DB,
	conns []*types.GitLabConnection,
	ap []schema.AuthProviders,
) *atypes.ProviderInitResult {
	initResults := &atypes.ProviderInitResult{}
	// Authorization (i.e., permissions) providers
	for _, c := range conns {
		p, err := newAuthzProvider(db, c, ap)
		if err != nil {
			initResults.InvalidConnections = append(initResults.InvalidConnections, extsvc.TypeGitLab)
			initResults.Problems = append(initResults.Problems, err.Error())
		} else if p != nil {
			initResults.Providers = append(initResults.Providers, p)
		}
	}

	return initResults
}

func newAuthzProvider(db database.DB, c *types.GitLabConnection, ps []schema.AuthProviders) (authz.Provider, error) {
	if c.Authorization == nil {
		return nil, nil
	}

	if errLicense := licensing.Check(licensing.FeatureACLs); errLicense != nil {
		return nil, errLicense
	}

	glURL, err := url.Parse(c.Url)
	if err != nil {
		return nil, errors.Errorf("Could not parse URL for GitLab instance %q: %s", c.Url, err)
	}

	switch idp := c.Authorization.IdentityProvider; {
	case idp.Oauth != nil:
		// Check that there is a GitLab authn provider corresponding to this GitLab instance
		foundAuthProvider := false
		syncInternalRepoPermissions := true
		for _, authnProvider := range ps {
			if authnProvider.Gitlab == nil {
				continue
			}
			authnURL := authnProvider.Gitlab.Url
			if authnURL == "" {
				authnURL = "https://gitlab.com"
			}
			authProviderURL, err := url.Parse(authnURL)
			if err != nil {
				// Ignore the error here, because the authn provider is responsible for its own validation
				continue
			}
			if authProviderURL.Hostname() == glURL.Hostname() {
				foundAuthProvider = true
				sirp := authnProvider.Gitlab.SyncInternalRepoPermissions
				syncInternalRepoPermissions = sirp == nil || *sirp
				break
			}
		}
		if !foundAuthProvider {
			return nil, errors.Errorf("Did not find authentication provider matching %q. Check the [**site configuration**](/site-admin/configuration) to verify an entry in [`auth.providers`](https://sourcegraph.com/docs/admin/auth) exists for %s.", c.Url, c.Url)
		}

		return NewOAuthProvider(OAuthProviderOp{
			URN:                         c.URN,
			BaseURL:                     glURL,
			Token:                       c.Token,
			TokenType:                   gitlab.TokenType(c.TokenType),
			DB:                          db,
			SyncInternalRepoPermissions: syncInternalRepoPermissions,
		}), nil
	case idp.Username != nil:
		return NewSudoProvider(SudoProviderOp{
			URN:                         c.URN,
			BaseURL:                     glURL,
			SudoToken:                   c.Token,
			SyncInternalRepoPermissions: !c.MarkInternalReposAsPublic,
		}), nil
	default:
		return nil, errors.Errorf("No identityProvider was specified")
	}
}

// NewOAuthProvider is a mockable constructor for new OAuthProvider instances.
var NewOAuthProvider = func(op OAuthProviderOp) authz.Provider {
	return newOAuthProvider(op, op.CLI)
}

// NewSudoProvider is a mockable constructor for new SudoProvider instances.
var NewSudoProvider = func(op SudoProviderOp) authz.Provider {
	return newSudoProvider(op, nil)
}

// ValidateAuthz validates the authorization fields of the given GitLab external
// service config.
func ValidateAuthz(cfg *schema.GitLabConnection, ps []schema.AuthProviders) error {
	_, err := newAuthzProvider(nil, &types.GitLabConnection{GitLabConnection: cfg}, ps)
	return err
}
