package commands

import (
	"flag"
	"testing"

	"github.com/spf13/afero"

	"github.com/urfave/cli"
)

func TestEnvCommand(t *testing.T) {
	t.Run("CanRunCommand", func(t *testing.T) {
		fileSystem := afero.NewMemMapFs()
		cmd := NewEnvCommand(fileSystem)

		app := cli.NewApp()
		set := &flag.FlagSet{}
		context := cli.NewContext(app, set, nil)
		cmd.ExecuteCommand(context)
	})
}
