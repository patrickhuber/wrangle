package commands_test

import (
	"testing"

	"github.com/patrickhuber/go-cross/console"
	"github.com/patrickhuber/go-cross/env"
	"github.com/patrickhuber/go-cross/filepath"
	"github.com/patrickhuber/go-cross/fs"
	"github.com/patrickhuber/go-cross/os"
	"github.com/patrickhuber/go-cross/platform"
	"github.com/patrickhuber/go-di"
	"github.com/patrickhuber/wrangle/internal/bootstrap"
	"github.com/patrickhuber/wrangle/internal/commands"
	"github.com/patrickhuber/wrangle/internal/config"
	"github.com/patrickhuber/wrangle/internal/diff"
	"github.com/patrickhuber/wrangle/internal/export"
	"github.com/patrickhuber/wrangle/internal/global"
	"github.com/patrickhuber/wrangle/internal/host"
	"github.com/stretchr/testify/require"
)

func TestExport(t *testing.T) {
	h := host.NewTest(platform.Linux, nil, nil)
	container := h.Container()

	processEnv, err := di.Resolve[env.Environment](container)
	require.NoError(t, err)
	require.NoError(t, processEnv.Set("MY_API_TOKEN", "test-token"))

	fs, err := di.Resolve[fs.FS](container)
	require.NoError(t, err)
	require.NotNil(t, fs)

	path, err := di.Resolve[filepath.Provider](container)
	require.NoError(t, err)
	require.NotNil(t, path)

	os, err := di.Resolve[os.OS](container)
	require.NoError(t, err)
	require.NotNil(t, os)

	wd, err := os.WorkingDirectory()
	require.NoError(t, err)
	require.NotEmpty(t, wd)

	// create a local yaml file with variables to export
	filePath := path.Join(wd, ".wrangle.yaml")
	err = fs.WriteFile(filePath, []byte(`
apiVersion: wrangle/v1
kind: Config
spec:
  env:
    TEST: TEST
    API_TOKEN: ((MY_API_TOKEN))
`), 0644)
	require.NoError(t, err)

	bt, err := di.Resolve[bootstrap.Service](container)
	require.NoError(t, err)
	require.NotNil(t, bt)

	export, err := di.Resolve[export.Service](container)
	require.NoError(t, err)
	require.NotNil(t, export)

	diff, err := di.Resolve[diff.Service](container)
	require.NoError(t, err)
	require.NotNil(t, diff)

	con, err := di.Resolve[console.Console](container)
	require.NoError(t, err)
	require.NotNil(t, con)

	// execute bootstrap to make sure we have a correct starting state
	err = bt.Execute(&bootstrap.Request{})
	require.NoError(t, err)

	configuration, err := di.Resolve[config.Service](container)
	require.NoError(t, err)
	cfg, err := configuration.Get()
	require.NoError(t, err)

	userConfigPath := cfg.Spec.Environment[global.EnvUserConfig]
	userCfg, err := config.ReadFile(fs, userConfigPath)
	require.NoError(t, err)
	userCfg.Spec.Stores = []config.Store{{
		Name: "env",
		Type: "env",
		Properties: map[string]any{
			"include": []string{"MY_API_TOKEN"},
		},
	}}
	require.NoError(t, config.WriteFile(fs, userConfigPath, userCfg))

	cfg, err = configuration.Get()
	require.NoError(t, err)
	require.Len(t, cfg.Spec.Stores, 2)

	exportCmd := &commands.ConfigExportCommand{
		Export: export,
		Diff:   diff,
		Options: commands.ConfigExportOptions{
			Shell: "bash",
		},
	}
	err = exportCmd.Execute()
	require.NoError(t, err)

	memoryConsole := con.(console.Memory)
	out := memoryConsole.OutBuffer().String()
	require.Contains(t, out, "export TEST='TEST'")
	require.Contains(t, out, "export API_TOKEN='test-token'")
}
