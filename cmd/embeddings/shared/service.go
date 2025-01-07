package shared

import (
	"context"

	"github.com/khulnasoft/khulnasoft/internal/debugserver"
	"github.com/khulnasoft/khulnasoft/internal/env"
	"github.com/khulnasoft/khulnasoft/internal/observation"
	"github.com/khulnasoft/khulnasoft/internal/service"
)

type svc struct{}

func (svc) Name() string { return "embeddings" }

func (svc) Configure() (env.Config, []debugserver.Endpoint) {
	var config Config
	config.Load()
	return &config, nil
}

func (svc) Start(ctx context.Context, observationCtx *observation.Context, ready service.ReadyFunc, config env.Config) error {
	return Main(ctx, observationCtx, ready, config.(*Config))
}

var Service service.Service = svc{}
