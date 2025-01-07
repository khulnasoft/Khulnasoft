package visualization

import (
	"github.com/khulnasoft/khulnasoft/lib/codeintel/lsif/reader"
)

type VisualizationContext struct {
	Stasher *reader.Stasher
}

func NewVisualizationContext() *VisualizationContext {
	return &VisualizationContext{
		Stasher: reader.NewStasher(),
	}
}
