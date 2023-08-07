package commands_test

import (
	"testing"

	"github.com/patrickhuber/go-di"
	"github.com/patrickhuber/go-xplat/platform"
	"github.com/patrickhuber/wrangle/internal/commands"
	"github.com/patrickhuber/wrangle/internal/host"
	"github.com/stretchr/testify/require"
)

func TestListPackages(t *testing.T) {

	s := host.NewTest(platform.Linux, nil, []string{})
	cmd := &commands.ListPackagesCommand{
		Options: &commands.ListPackagesOptions{
			Output: ".",
		},
	}

	err := di.Inject(s.Container(), cmd)
	require.NoError(t, err)

	err = cmd.Execute()
	require.NoError(t, err)
}
