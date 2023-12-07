package commands_test

import (
	"testing"

	"github.com/patrickhuber/go-di"
	"github.com/patrickhuber/go-xplat/fs"
	"github.com/patrickhuber/go-xplat/platform"
	"github.com/patrickhuber/wrangle/internal/commands"
	"github.com/patrickhuber/wrangle/internal/config"
	"github.com/patrickhuber/wrangle/internal/host"
	"github.com/patrickhuber/wrangle/internal/services"
	"github.com/stretchr/testify/require"
)

func TestListPackages(t *testing.T) {

	h := host.NewTest(platform.Linux, nil, []string{})
	cmd := &commands.ListPackagesCommand{
		Options: &commands.ListPackagesOptions{
			Output: ".",
		},
	}

	fs, err := di.Resolve[fs.FS](h.Container())
	require.Nil(t, err)

	configuration, err := di.Resolve[services.Configuration](h.Container())
	require.Nil(t, err)

	err = config.WriteFile(fs, configuration.DefaultGlobalConfigFilePath(), configuration.GlobalDefault())
	require.Nil(t, err)

	err = di.Inject(h.Container(), cmd)
	require.NoError(t, err)

	err = cmd.Execute()
	require.NoError(t, err)
}
