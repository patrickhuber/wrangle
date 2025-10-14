package commands_test

import (
	"testing"

	"github.com/patrickhuber/go-cross/console"
	"github.com/patrickhuber/go-cross/filepath"
	"github.com/patrickhuber/go-cross/fs"
	"github.com/patrickhuber/go-cross/os"
	"github.com/patrickhuber/go-cross/platform"
	"github.com/patrickhuber/go-di"
	"github.com/patrickhuber/wrangle/internal/bootstrap"
	"github.com/patrickhuber/wrangle/internal/commands"
	"github.com/patrickhuber/wrangle/internal/diff"
	"github.com/patrickhuber/wrangle/internal/export"
	"github.com/patrickhuber/wrangle/internal/host"
	"github.com/stretchr/testify/require"
)

func TestExport(t *testing.T) {
	h := host.NewTest(platform.Linux, nil, nil)
	container := h.Container()

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

	// create a local toml file with a TEST variable
	// order matters here, this file must be in place before configuration dependent services run
	filePath := path.Join(wd, ".wrangle.toml")

	err = fs.WriteFile(filePath, []byte(`
[spec.env]
TEST="TEST"`), 0644)
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

	exportCmd := &commands.ExportCommand{
		Export: export,
		Diff:   diff,
		Options: commands.ExportOptions{
			Shell: "bash",
		},
	}
	err = exportCmd.Execute()
	require.NoError(t, err)

	memoryConsole := con.(console.Memory)
	out := memoryConsole.OutBuffer().String()
	require.Contains(t, out, "export TEST='TEST'")
}
