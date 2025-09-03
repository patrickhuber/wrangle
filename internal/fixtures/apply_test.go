package fixtures_test

import (
	"testing"

	"github.com/patrickhuber/go-cross"
	"github.com/patrickhuber/go-cross/arch"
	"github.com/patrickhuber/go-cross/platform"
	"github.com/patrickhuber/wrangle/internal/fixtures"
	"github.com/stretchr/testify/require"
)

func TestApply(t *testing.T) {
	tests := []struct {
		name     string
		platform platform.Platform
		arch     arch.Arch
	}{
		{"Windows", platform.Windows, arch.AMD64},
		{"Linux", platform.Linux, arch.AMD64},
		{"Darwin", platform.Darwin, arch.AMD64},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			target := cross.NewTest(tt.platform, tt.arch)

			err := fixtures.Apply(target.OS(), target.FS(), target.Env())
			require.NoError(t, err)

			home, err := target.OS().Home()
			require.NoError(t, err)

			ok, err := target.FS().Exists(home)
			require.NoError(t, err)
			require.True(t, ok, "Home directory '%s' should exist", home)

			working, err := target.OS().WorkingDirectory()
			require.NoError(t, err)

			ok, err = target.FS().Exists(working)
			require.NoError(t, err)
			require.True(t, ok, "Working directory '%s' should exist", working)
		})
	}
}
