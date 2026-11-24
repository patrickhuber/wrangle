package commands_test

import (
	"testing"

	"github.com/patrickhuber/go-cross/fs"
	"github.com/patrickhuber/go-cross/os"
	"github.com/patrickhuber/go-cross/platform"
	"github.com/patrickhuber/go-di"
	"github.com/patrickhuber/wrangle/internal/bootstrap"
	"github.com/patrickhuber/wrangle/internal/commands"
	"github.com/patrickhuber/wrangle/internal/config"
	"github.com/patrickhuber/wrangle/internal/global"
	"github.com/patrickhuber/wrangle/internal/host"
	"github.com/stretchr/testify/require"
)

func TestStoreAdd(t *testing.T) {
	h := host.NewTest(platform.Linux, nil, []string{})
	container := h.Container()

	bootstrapService, err := di.Resolve[bootstrap.Service](container)
	require.NoError(t, err)
	require.NoError(t, bootstrapService.Execute(&bootstrap.Request{}))

	cmd := &commands.StoreAddCommand{
		Options: commands.StoreAddOptions{
			Name:         "github",
			Type:         "env",
			Include:      []string{"GITHUB_TOKEN"},
			Properties:   []string{"description=GitHub credentials"},
			Dependencies: []string{"shared"},
			Scope:        "local",
		},
	}
	require.NoError(t, di.Inject(container, cmd))
	require.NoError(t, cmd.Execute())

	opsys, err := di.Resolve[os.OS](container)
	require.NoError(t, err)
	workDir, err := opsys.WorkingDirectory()
	require.NoError(t, err)
	filesystem, err := di.Resolve[fs.FS](container)
	require.NoError(t, err)

	cfg, err := config.ReadFile(filesystem, workDir+"/"+global.LocalConfigurationFileName)
	require.NoError(t, err)
	require.Len(t, cfg.Spec.Stores, 1)
	require.Equal(t, "github", cfg.Spec.Stores[0].Name)
	require.Equal(t, "env", cfg.Spec.Stores[0].Type)
	require.Equal(t, []any{"GITHUB_TOKEN"}, cfg.Spec.Stores[0].Properties["include"])
	require.Equal(t, "GitHub credentials", cfg.Spec.Stores[0].Properties["description"])
	require.Equal(t, []string{"shared"}, cfg.Spec.Stores[0].Dependencies)
}

func TestStoreAddUpdatesStoreWithSameName(t *testing.T) {
	h := host.NewTest(platform.Linux, nil, []string{})
	container := h.Container()

	bootstrapService, err := di.Resolve[bootstrap.Service](container)
	require.NoError(t, err)
	require.NoError(t, bootstrapService.Execute(&bootstrap.Request{}))

	for _, include := range [][]string{{"OLD_TOKEN"}, {"GITHUB_TOKEN"}} {
		cmd := &commands.StoreAddCommand{
			Options: commands.StoreAddOptions{
				Name:    "github",
				Type:    "env",
				Include: include,
				Scope:   "local",
			},
		}
		require.NoError(t, di.Inject(container, cmd))
		require.NoError(t, cmd.Execute())
	}

	opsys, err := di.Resolve[os.OS](container)
	require.NoError(t, err)
	workDir, err := opsys.WorkingDirectory()
	require.NoError(t, err)
	filesystem, err := di.Resolve[fs.FS](container)
	require.NoError(t, err)

	cfg, err := config.ReadFile(filesystem, workDir+"/"+global.LocalConfigurationFileName)
	require.NoError(t, err)
	require.Len(t, cfg.Spec.Stores, 1)
	require.Equal(t, []any{"GITHUB_TOKEN"}, cfg.Spec.Stores[0].Properties["include"])
}
