package dependencies

import (
	"context"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/khulnasoft/khulnasoft/dev/sg/internal/check"
	"github.com/khulnasoft/khulnasoft/dev/sg/internal/std"
	"github.com/khulnasoft/khulnasoft/dev/sg/internal/usershell"
)

func TestMacFix(t *testing.T) {
	if !strings.Contains(*sgSetupTests, string(OSMac)) && !strings.Contains(*sgSetupTests, "macos") {
		t.Skip("Skipping Mac sg setup tests")
	}

	// Initialize context with user shell information
	ctx, err := usershell.Context(context.Background())
	require.NoError(t, err)

	// Set up runner with no input and simple output
	runner := check.NewRunner(nil, std.NewSimpleOutput(os.Stdout, true), Mac)

	// automatically fix everything!
	t.Run("Fix", func(t *testing.T) {
		err = runner.Fix(ctx, testArgs)
		require.Nil(t, err)
	})

	// now check that everything was fixed
	t.Run("Check", func(t *testing.T) {
		err = runner.Check(ctx, testArgs)
		assert.Nil(t, err)
	})
}
