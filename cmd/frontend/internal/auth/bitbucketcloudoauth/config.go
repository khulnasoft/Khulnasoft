package bitbucketcloudoauth

import (
	"fmt"

	"github.com/sourcegraph/log"

	"github.com/khulnasoft/khulnasoft/cmd/frontend/internal/auth/providers"
	"github.com/khulnasoft/khulnasoft/internal/collections"
	"github.com/khulnasoft/khulnasoft/internal/conf"
	"github.com/khulnasoft/khulnasoft/internal/conf/conftypes"
	"github.com/khulnasoft/khulnasoft/internal/database"
	"github.com/khulnasoft/khulnasoft/internal/licensing"
	"github.com/khulnasoft/khulnasoft/schema"
)

func Init(logger log.Logger, db database.DB) {
	const pkgName = "bitbucketcloudoauth"
	logger = logger.Scoped(pkgName)
	conf.ContributeValidator(func(cfg conftypes.SiteConfigQuerier) conf.Problems {
		_, problems := parseConfig(logger, cfg, db)
		return problems
	})

	go conf.Watch(func() {
		newProviders, _ := parseConfig(logger, conf.Get(), db)
		if len(newProviders) == 0 {
			providers.Update(pkgName, nil)
			return
		}

		if err := licensing.Check(licensing.FeatureSSO); err != nil {
			logger.Error("Check license for SSO (Bitbucket Cloud OAuth)", log.Error(err))
			providers.Update(pkgName, nil)
			return
		}

		newProvidersList := make([]providers.Provider, 0, len(newProviders))
		for _, p := range newProviders {
			newProvidersList = append(newProvidersList, p.Provider)
		}
		providers.Update(pkgName, newProvidersList)
	})
}

type Provider struct {
	*schema.BitbucketCloudAuthProvider
	providers.Provider
}

func parseConfig(logger log.Logger, cfg conftypes.SiteConfigQuerier, db database.DB) (ps []Provider, problems conf.Problems) {
	existingProviders := make(collections.Set[string])

	for _, pr := range cfg.SiteConfig().AuthProviders {
		if pr.Bitbucketcloud == nil {
			continue
		}

		provider, providerProblems := parseProvider(logger, pr.Bitbucketcloud, db, pr)
		problems = append(problems, conf.NewSiteProblems(providerProblems...)...)
		if provider == nil {
			continue
		}

		if existingProviders.Has(provider.CachedInfo().UniqueID()) {
			problems = append(problems, conf.NewSiteProblems(fmt.Sprintf(`Cannot have more than one Bitbucket Cloud auth provider with url %q and client ID %q, only the first one will be used`, provider.ServiceID, provider.CachedInfo().ClientID))...)
			continue
		}

		ps = append(ps, Provider{
			BitbucketCloudAuthProvider: pr.Bitbucketcloud,
			Provider:                   provider,
		})
		existingProviders.Add(provider.CachedInfo().UniqueID())
	}
	return ps, problems
}
