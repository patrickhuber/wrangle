package commands

import (
	"bytes"
	"testing"

	"github.com/patrickhuber/wrangle/config"
	"github.com/patrickhuber/wrangle/ui"
	"github.com/stretchr/testify/require"
)

func TestPackagesCommand(t *testing.T) {
	t.Run("CanExecutePackagesCommand", func(t *testing.T) {
		r := require.New(t)

		// create console
		console := ui.NewMemoryConsole()

		// create the command
		pkgs := NewPackages(console)

		// create the configuration
		content := `
packages:
- name: one
  version: 0.1.1
- name: two
  version: 2.3.1
`
		configuration, err := config.SerializeString(content)
		r.Nil(err)

		// execute the command
		err = pkgs.Execute(configuration)
		r.Nil(err)

		// verify output
		b, ok := console.Out().(*bytes.Buffer)
		r.True(ok)
		r.NotNil(b)
		r.Equal("one - 0.1.1\ntwo - 2.3.1\n", b.String())
	})
}
