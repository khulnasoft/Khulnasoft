package azureoauth

import (
	"net/http"

	"github.com/khulnasoft/khulnasoft/cmd/frontend/auth"
	"github.com/khulnasoft/khulnasoft/cmd/frontend/internal/auth/oauth"
	"github.com/khulnasoft/khulnasoft/internal/database"
	"github.com/khulnasoft/khulnasoft/internal/extsvc"
)

func Middleware(db database.DB) *auth.Middleware {
	return &auth.Middleware{
		API: func(next http.Handler) http.Handler {
			return oauth.NewMiddleware(db, extsvc.TypeAzureDevOps, authPrefix, true, next)
		},
		App: func(next http.Handler) http.Handler {
			return oauth.NewMiddleware(db, extsvc.TypeAzureDevOps, authPrefix, false, next)
		},
	}
}
