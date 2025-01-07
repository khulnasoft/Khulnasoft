package debugproxies

import (
	"database/sql"
	"fmt"
	"net/http"
	"net/http/httputil"
	"sort"
	"strings"
	"sync"

	"github.com/gorilla/mux"

	"github.com/khulnasoft/khulnasoft/cmd/frontend/internal/app/errorutil"
	"github.com/khulnasoft/khulnasoft/internal/actor"
	"github.com/khulnasoft/khulnasoft/internal/auth"
	"github.com/khulnasoft/khulnasoft/internal/database"
	"github.com/khulnasoft/khulnasoft/internal/errcode"
	"github.com/khulnasoft/khulnasoft/lib/errors"
)

// proxyEndpoint couples the reverse proxy with the endpoint it proxies.
type proxyEndpoint struct {
	reverseProxy http.Handler
	host         string
}

// ReverseProxyHandler handles serving the index page and routing the requests being proxied to their
// respective reverse proxy. proxyEndpoints come from callers calling ReverseProxyHandler.Populate().
// zero value is useful and will provide a "no endpoint found" index until some endpoints get populated.
type ReverseProxyHandler struct {
	// protects the reverseProxies map
	sync.RWMutex
	// keys are the displayNames
	reverseProxies map[string]*proxyEndpoint
}

func (rph *ReverseProxyHandler) AddToRouter(r *mux.Router, db database.DB) {
	r.Handle("/", AdminOnly(db, http.HandlerFunc(rph.serveIndex)))
	r.PathPrefix("/proxies").Handler(http.StripPrefix("/-/debug/proxies", AdminOnly(db, errorutil.Handler(rph.serveReverseProxy))))
}

// serveIndex composes the simple index page with the endpoints sorted by their name.
func (rph *ReverseProxyHandler) serveIndex(w http.ResponseWriter, r *http.Request) {
	rph.RLock()
	displayNames := make([]string, 0, len(rph.reverseProxies))
	for displayName := range rph.reverseProxies {
		displayNames = append(displayNames, displayName)
	}
	rph.RUnlock()

	if len(displayNames) == 0 {
		fmt.Fprintf(w, `Instrumentation: no endpoints found<br>`)
		fmt.Fprintf(w, `<br><br><a href="headers">headers</a><br>`)
		return
	}

	sort.Strings(displayNames)

	for _, displayName := range displayNames {
		fmt.Fprintf(w, `<a href="proxies/%s/">%s</a><br>`, displayName, displayName)
	}
	fmt.Fprintf(w, `<br><br><a href="headers">headers</a><br>`)
}

// serveReverseProxy routes the request to the appropriate reverse proxy by splitting the request path and finding
// the displayName under which the proxy lives.
func (rph *ReverseProxyHandler) serveReverseProxy(w http.ResponseWriter, r *http.Request) error {
	pathParts := strings.Split(r.URL.Path, "/")
	if len(pathParts) < 2 {
		return &errcode.HTTPErr{
			Status: http.StatusNotFound,
			Err:    errors.New("proxy endpoint missing"),
		}
	}

	var pe *proxyEndpoint
	rph.RLock()
	if len(rph.reverseProxies) > 0 {
		pe = rph.reverseProxies[pathParts[1]]
	}
	rph.RUnlock()

	if pe == nil {
		return &errcode.HTTPErr{
			Status: http.StatusNotFound,
			Err:    errors.New("proxy endpoint missing"),
		}
	}

	pe.reverseProxy.ServeHTTP(w, r)
	return nil
}

// Populate declares the proxyEndpoints to use. This method can be called at any time from any goroutine.
// It completely replaces the previous proxied endpoints with the ones specified in the call.
func (rph *ReverseProxyHandler) Populate(db database.DB, peps []Endpoint) {
	rps := make(map[string]*proxyEndpoint, len(peps))
	for _, ep := range peps {
		displayName := displayNameFromEndpoint(ep)
		rps[displayName] = &proxyEndpoint{
			reverseProxy: reverseProxyFromHost(db, ep.Addr, displayName),
			host:         ep.Addr,
		}
	}

	rph.Lock()
	rph.reverseProxies = rps
	rph.Unlock()
}

// Creates a display name from an endpoint suited for using in a URL link.
func displayNameFromEndpoint(ep Endpoint) string {
	host := ep.Hostname
	if host == "" {
		host = ep.Addr
		if idx := strings.Index(host, ":"); idx != -1 {
			host = host[:idx]
		}
	}

	// Stateful Services have unique pod names. Lets use them to avoid stutter
	// in the name (gitserver-gitserver-0 becomes gitserver-0).
	if strings.HasPrefix(host, ep.Service) {
		return host
	}
	return fmt.Sprintf("%s-%s", ep.Service, host)
}

// reverseProxyFromHost creates a reverse proxy from specified host with the path prefix that will be stripped from
// request before it gets sent to the destination endpoint.
func reverseProxyFromHost(db database.DB, host string, pathPrefix string) http.Handler {
	// 🚨 SECURITY: Only admins can create reverse proxies from host
	return AdminOnly(db, &httputil.ReverseProxy{
		Director: func(req *http.Request) {
			req.URL.Scheme = "http"
			req.URL.Host = host
			if i := strings.Index(req.URL.Path, pathPrefix); i >= 0 {
				req.URL.Path = req.URL.Path[i+len(pathPrefix):]
			}
		},
	})
}

// AdminOnly is an HTTP middleware which only allows requests by admins.
func AdminOnly(db database.DB, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ff, err := db.FeatureFlags().GetFeatureFlag(r.Context(), "sourcegraph-operator-site-admin-hide-maintenance")
		if err == nil {
			hide, _ := ff.EvaluateGlobal()
			a := actor.FromContext(r.Context())
			if hide && !a.SourcegraphOperator {
				http.Error(w, "Only Sourcegraph operators are allowed", http.StatusForbidden)
				return
			}
		} else if err != nil && err != sql.ErrNoRows {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if err := auth.CheckCurrentUserIsSiteAdmin(r.Context(), db); err != nil {
			http.Error(w, err.Error(), http.StatusForbidden)
			return
		}

		next.ServeHTTP(w, r)
	})
}
