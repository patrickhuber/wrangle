package commands

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/patrickhuber/cli-mgr/store"
	"github.com/patrickhuber/cli-mgr/store/memory"

	"github.com/patrickhuber/cli-mgr/ui"
	"github.com/spf13/afero"
)

func TestEnvCommand(t *testing.T) {
	t.Run("CanRunCommand", func(t *testing.T) {
		r := require.New(t)

		// create filesystem
		fileSystem := afero.NewMemMapFs()

		// create store manager
		manager := store.NewManager()
		manager.Register(memory.NewMemoryStoreProvider())

		// create console
		console := ui.NewMemoryConsole()

		// create and run command
		cmd := NewEnvCommand(manager, fileSystem, "linux", console)
		runCommandParams := NewRunCommandParams("/config", "echo", "lab")
		cmd.ExecuteCommand(runCommandParams)

		// verify output
		b, ok := console.Out().(*bytes.Buffer)
		r.True(ok)
		r.NotNil(b)
		r.Equal("export CLI_MGR_TEST=value\n", b.String())
	})
}
