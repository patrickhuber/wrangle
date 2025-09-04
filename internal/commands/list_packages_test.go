package commands_test

import (
	"testing"

	"github.com/patrickhuber/go-cross/platform"
	"github.com/patrickhuber/go-di"
	"github.com/patrickhuber/wrangle/internal/bootstrap"
	"github.com/patrickhuber/wrangle/internal/commands"
	"github.com/patrickhuber/wrangle/internal/host"
	"github.com/stretchr/testify/require"
)

func TestListPackages(t *testing.T) {

	h := host.NewTest(platform.Linux, nil, []string{})
	container := h.Container()

	bootstrapService, err := di.Resolve[bootstrap.Service](container)
	require.NoError(t, err)

	err = bootstrapService.Execute(&bootstrap.Request{})
	require.NoError(t, err)

	cmd := &commands.ListPackagesCommand{
		Options: &commands.ListPackagesOptions{
			Output: ".",
		},
	}

	err = di.Inject(container, cmd)
	require.NoError(t, err)

	err = cmd.Execute()
	require.NoError(t, err)
}
