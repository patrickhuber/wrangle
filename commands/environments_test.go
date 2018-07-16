package commands_test

import (
	"bytes"
	"testing"

	"github.com/patrickhuber/wrangle/commands"
	"github.com/patrickhuber/wrangle/config"
	"github.com/patrickhuber/wrangle/ui"
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
		command := commands.NewEnvironments(fileSystem, console)
		loader := config.NewLoader(fileSystem)
		configuration, err := loader.Load("/test")
		r.Nil(err)
		err = command.Execute(configuration)
		r.Nil(err)
		b, ok := console.Out().(*bytes.Buffer)
		r.True(ok)
		r.NotNil(b)

		r.Equal("name\n----\none\ntwo\nthree\n", b.String())
	})
}
