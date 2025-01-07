package types

import (
	"time"

	"github.com/khulnasoft/khulnasoft/schema"
)

const (
	DequeueCachePrefix = "executor_multihandler_dequeues"
	DequeueTtl         = 5 * time.Minute
	CleanupInterval    = 30 * time.Second
)

var DequeuePropertiesPerQueue = &schema.DequeueCacheConfig{
	Batches: &schema.Batches{
		Limit:  50,
		Weight: 4,
	},
	Codeintel: &schema.Codeintel{
		Limit:  250,
		Weight: 1,
	},
}
