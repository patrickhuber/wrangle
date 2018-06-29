package commands_test

import (
	"bytes"
	"testing"

	"github.com/patrickhuber/cli-mgr/commands"
	"github.com/patrickhuber/cli-mgr/config"
	"github.com/patrickhuber/cli-mgr/ui"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/require"
)

func TestEnvironmentsCommand(t *testing.T) {
	t.Run("CanGetListOfEnvironments", func(t *testing.T) {
		r := require.New(t)
		fileSystem := afero.NewMemMapFs()
		content := `
environments:
- name: one
- name: two
- name: three
`
		afero.WriteFile(fileSystem, "/test", []byte(content), 0644)

		console := ui.NewMemoryConsole()
		command := commands.NewEnvironmentsCommand(fileSystem, console)
		loader := config.NewLoader(fileSystem)
		configuration, err := loader.Load("/test")
		r.Nil(err)
		err = command.ExecuteCommand(configuration)
		r.Nil(err)
		b, ok := console.Out().(*bytes.Buffer)
		r.True(ok)
		r.NotNil(b)

		r.Equal("one\ntwo\nthree\n", b.String())
	})
}
