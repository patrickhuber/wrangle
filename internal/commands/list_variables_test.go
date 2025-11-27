package commands_test

import (
	"testing"

	"github.com/patrickhuber/go-cross/platform"
	"github.com/patrickhuber/go-di"
	"github.com/patrickhuber/wrangle/internal/bootstrap"
	"github.com/patrickhuber/wrangle/internal/commands"
	"github.com/patrickhuber/wrangle/internal/config"
	"github.com/patrickhuber/wrangle/internal/host"
	"github.com/stretchr/testify/require"
)

func TestListVariables(t *testing.T) {
	t.Run("lists variables from configuration", func(t *testing.T) {
		h := host.NewTest(platform.Linux, nil, []string{})
		container := h.Container()

		bootstrapService, err := di.Resolve[bootstrap.Service](container)
		require.NoError(t, err)

		err = bootstrapService.Execute(&bootstrap.Request{})
		require.NoError(t, err)

		// Get configuration and add some variables
		configuration, err := di.Resolve[config.Configuration](container)
		require.NoError(t, err)

		// Verify we can get the configuration
		cfg, err := configuration.Get()
		require.NoError(t, err)
		require.NotNil(t, cfg)

		cmd := &commands.ListVariablesCommand{
			Options: &commands.ListVariablesOptions{
				Output: "table",
			},
		}

		err = di.Inject(container, cmd)
		require.NoError(t, err)

		err = cmd.Execute()
		require.NoError(t, err)
	})
}
