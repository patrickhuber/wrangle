package commands

import (
	"flag"
	"testing"

	"github.com/urfave/cli"
)

func TestEnvCommand(t *testing.T) {
	t.Run("CanRunCommand", func(t *testing.T) {
		cmd := NewEnvCommand()

		app := cli.NewApp()
		set := &flag.FlagSet{}
		context := cli.NewContext(app, set, nil)
		cmd.ExecuteCommand(context)
	})
}
