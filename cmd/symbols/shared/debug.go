package shared

import (
	"github.com/khulnasoft/khulnasoft/internal/debugserver"
)

// GRPCWebUIDebugEndpoint returns a debug endpoint that serves the GRPCWebUI that targets
// this symbols instance.
func GRPCWebUIDebugEndpoint() debugserver.Endpoint {
	return debugserver.NewGRPCWebUIEndpoint("symbols", addr)
}
