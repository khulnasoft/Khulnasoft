package completions

import (
	"context"
	"net/http"

	"github.com/sourcegraph/log"

	"github.com/khulnasoft/khulnasoft/cmd/frontend/enterprise"
	"github.com/khulnasoft/khulnasoft/cmd/frontend/internal/cody"
	"github.com/khulnasoft/khulnasoft/cmd/frontend/internal/completions/resolvers"
	"github.com/khulnasoft/khulnasoft/cmd/frontend/internal/guardrails"
	"github.com/khulnasoft/khulnasoft/cmd/frontend/internal/httpapi/completions"
	"github.com/khulnasoft/khulnasoft/internal/codeintel"
	"github.com/khulnasoft/khulnasoft/internal/conf/conftypes"
	"github.com/khulnasoft/khulnasoft/internal/database"
	"github.com/khulnasoft/khulnasoft/internal/observation"
)

func Init(
	_ context.Context,
	observationCtx *observation.Context,
	db database.DB,
	_ codeintel.Services,
	conf conftypes.UnifiedWatchable,
	enterpriseServices *enterprise.Services,
) error {
	logger := log.Scoped("completions")

	enterpriseServices.NewChatCompletionsStreamHandler = func() http.Handler {
		completionsHandler := completions.NewChatCompletionsStreamHandler(logger, db)
		return requireVerifiedEmailMiddleware(db, observationCtx.Logger, completionsHandler)
	}
	enterpriseServices.NewCodeCompletionsHandler = func() http.Handler {
		codeCompletionsHandler := completions.NewCodeCompletionsHandler(logger, db, guardrails.NewAttributionTest(observationCtx, conf))
		return requireVerifiedEmailMiddleware(db, observationCtx.Logger, codeCompletionsHandler)
	}
	enterpriseServices.CompletionsResolver = resolvers.NewCompletionsResolver(db, observationCtx.Logger)

	return nil
}

func requireVerifiedEmailMiddleware(db database.DB, logger log.Logger, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := cody.CheckVerifiedEmailRequirement(r.Context(), db, logger); err != nil {
			// Report HTTP 403 Forbidden if user has no verified email address.
			http.Error(w, err.Error(), http.StatusForbidden)
			return
		}

		next.ServeHTTP(w, r)
	})
}
