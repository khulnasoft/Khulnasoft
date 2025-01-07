//go:build shell
// +build shell

package util_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/khulnasoft/khulnasoft/cmd/executor/internal/util"
)

func TestHasShellBuildTag_Shell(t *testing.T) {
	assert.True(t, util.HasShellBuildTag())
}
