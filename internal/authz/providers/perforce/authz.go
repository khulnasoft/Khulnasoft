package perforce

import (
	"strings"

	"github.com/sourcegraph/log"

	"github.com/khulnasoft/khulnasoft/internal/database"
	"github.com/khulnasoft/khulnasoft/internal/gitserver"
	"github.com/khulnasoft/khulnasoft/internal/licensing"

	"github.com/khulnasoft/khulnasoft/internal/authz"
	atypes "github.com/khulnasoft/khulnasoft/internal/authz/types"
	"github.com/khulnasoft/khulnasoft/internal/extsvc"
	"github.com/khulnasoft/khulnasoft/internal/types"
	"github.com/khulnasoft/khulnasoft/schema"
)

// NewAuthzProviders returns the set of Perforce authz providers derived from the connections.
//
// It also returns any simple validation problems with the config, separating these into "serious problems"
// and "warnings". "Serious problems" are those that should make Sourcegraph set authz.allowAccessByDefault
// to false. "Warnings" are all other validation problems.
//
// This constructor does not and should not directly check connectivity to external services - if
// desired, callers should use `(*Provider).ValidateConnection` directly to get warnings related
// to connection issues.
func NewAuthzProviders(db database.DB, conns []*types.PerforceConnection) *atypes.ProviderInitResult {
	initResults := &atypes.ProviderInitResult{}
	for _, c := range conns {
		p, err := newAuthzProvider(c.URN, db, c.Authorization, c.P4Port, c.P4User, c.P4Passwd, c.Depots)
		if err != nil {
			initResults.InvalidConnections = append(initResults.InvalidConnections, extsvc.TypePerforce)
			initResults.Problems = append(initResults.Problems, err.Error())
		} else if p != nil {
			initResults.Providers = append(initResults.Providers, p)
		}
	}

	return initResults
}

func newAuthzProvider(
	urn string,
	db database.DB,
	a *schema.PerforceAuthorization,
	host, user, password string,
	depots []string,
) (authz.Provider, error) {
	// Call this function from ValidateAuthz if this function starts returning an error.
	if a == nil {
		return nil, nil
	}

	logger := log.Scoped("authz")
	if err := licensing.Check(licensing.FeatureACLs); err != nil {
		return nil, err
	}

	var depotIDs []extsvc.RepoID
	if a.SubRepoPermissions {
		depotIDs = make([]extsvc.RepoID, len(depots))
		for i, depot := range depots {
			// Force depots as directories
			if strings.HasSuffix(depot, "/") {
				depotIDs[i] = extsvc.RepoID(depot)
			} else {
				depotIDs[i] = extsvc.RepoID(depot + "/")
			}
		}
	}

	return NewProvider(logger, db, gitserver.NewClient("authz.perforce"), urn, host, user, password, depotIDs, a.IgnoreRulesWithHost), nil
}

// ValidateAuthz validates the authorization fields of the given Perforce
// external service config.
func ValidateAuthz(_ *schema.PerforceConnection) error {
	// newAuthzProvider always succeeds, so directly return nil here.
	return nil
}
