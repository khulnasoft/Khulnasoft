package gerrit

import (
	atypes "github.com/khulnasoft/khulnasoft/internal/authz/types"
	"github.com/khulnasoft/khulnasoft/internal/extsvc"
	"github.com/khulnasoft/khulnasoft/internal/licensing"
	"github.com/khulnasoft/khulnasoft/internal/types"
)

// NewAuthzProviders returns the set of Gerrit authz providers derived from the connections.
func NewAuthzProviders(conns []*types.GerritConnection) *atypes.ProviderInitResult {
	initResults := &atypes.ProviderInitResult{}

	for _, c := range conns {
		if c.Authorization == nil {
			// No authorization required
			continue
		}

		if err := licensing.Check(licensing.FeatureACLs); err != nil {
			initResults.InvalidConnections = append(initResults.InvalidConnections, extsvc.TypeGerrit)
			initResults.Problems = append(initResults.Problems, err.Error())
			continue
		}

		p, err := NewProvider(c)
		if err != nil {
			initResults.InvalidConnections = append(initResults.InvalidConnections, extsvc.TypeGerrit)
			initResults.Problems = append(initResults.Problems, err.Error())
		}
		if p != nil {
			initResults.Providers = append(initResults.Providers, p)
		}
	}
	return initResults
}
