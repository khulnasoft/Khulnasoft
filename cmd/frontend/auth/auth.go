// Package auth contains auth related code for the frontend.
package auth

import (
	"net/http"

	"github.com/khulnasoft/khulnasoft/cmd/frontend/internal/auth/userpasswd"
	"github.com/khulnasoft/khulnasoft/internal/lazyregexp"
)

// AuthURLPrefix is the URL path prefix under which to attach authentication handlers
const AuthURLPrefix = "/.auth"

// Middleware groups two related middlewares (one for the API, one for the app).
type Middleware struct {
	// API is the middleware that performs authentication on the API handler.
	API func(http.Handler) http.Handler

	// App is the middleware that performs authentication on the app handler.
	App func(http.Handler) http.Handler
}

var extraAuthMiddlewares []*Middleware

// RegisterMiddlewares registers additional authentication middlewares. Currently this is used to
// register enterprise-only SSO middleware. This should only be called from an init function.
func RegisterMiddlewares(m ...*Middleware) {
	extraAuthMiddlewares = append(extraAuthMiddlewares, m...)
}

// AuthMiddleware returns the authentication middleware that combines all authentication middlewares
// that have been registered.
func AuthMiddleware() *Middleware {
	m := make([]*Middleware, 0, 1+len(extraAuthMiddlewares))
	m = append(m, RequireAuthMiddleware)
	m = append(m, extraAuthMiddlewares...)
	return composeMiddleware(m...)
}

// composeMiddleware returns a new Middleware that composes the middlewares together.
func composeMiddleware(middlewares ...*Middleware) *Middleware {
	return &Middleware{
		API: func(h http.Handler) http.Handler {
			for _, m := range middlewares {
				h = m.API(h)
			}
			return h
		},
		App: func(h http.Handler) http.Handler {
			for _, m := range middlewares {
				h = m.App(h)
			}
			return h
		},
	}
}

// NormalizeUsername normalizes a proposed username into a format that meets Sourcegraph's
// username formatting rules.
func NormalizeUsername(name string) (string, error) {
	return userpasswd.NormalizeUsername(name)
}

// Equivalent to `^\w(?:\w|[-.](?=\w))*-?$` which we have in the DB constraint, but without a lookahead
var validUsername = lazyregexp.New(`^\w(?:(?:[\w.-]\w|\w)*-?|)$`)

// IsValidUsername returns true if the username matches the constraints in the database.
func IsValidUsername(name string) bool {
	return validUsername.MatchString(name) && len(name) <= 255
}
