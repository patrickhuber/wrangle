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
		fileSystem := afero.NewMemMapFs()
		manager := store.NewManager()
		console := ui.NewMemoryConsole()
		manager.Register(memory.NewMemoryStoreProvider())
		cmd := NewEnvCommand(manager, fileSystem, "linux", console)
		runCommandParams := NewRunCommandParams("/config", "echo", "lab")
		cmd.ExecuteCommand(runCommandParams)
		b, ok := console.Out().(*bytes.Buffer)
		r.True(ok)
		r.NotNil(b)
		r.Equal("export CLI_MGR_TEST=value\n", b.String())
	})
}
