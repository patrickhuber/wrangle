package commands_test

import (
	"bytes"
	"testing"

	"github.com/patrickhuber/cli-mgr/commands"
	"github.com/patrickhuber/cli-mgr/ui"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/require"
)

func TestProcessesCommand(t *testing.T) {
	t.Run("CanGetListOfProcesses", func(t *testing.T) {
		r := require.New(t)
		fileSystem := afero.NewMemMapFs()
		content := `
processes:
- name: one
- name: two
- name: three
`
		afero.WriteFile(fileSystem, "/test", []byte(content), 0644)

		console := ui.NewMemoryConsole()
		command := commands.NewProcessesCommand(fileSystem, console)
		err := command.ExecuteCommand("/test")
		r.Nil(err)
		b, ok := console.Out().(*bytes.Buffer)
		r.True(ok)
		r.NotNil(b)

		r.Equal("one\ntwo\nthree\n", b.String())
	})
}
