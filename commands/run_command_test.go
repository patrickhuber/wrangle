package commands

import (
	"flag"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/urfave/cli"

	"github.com/patrickhuber/cli-mgr/config"
	"github.com/spf13/afero"
)

func TestRunCommand(t *testing.T) {
	r := require.New(t)
	configStoreManager := config.NewConfigStoreManager()
	fileSystem := afero.NewMemMapFs()
	runCommand := NewRunCommand(configStoreManager, fileSystem)
	app := cli.NewApp()
	set := &flag.FlagSet{}
	context := cli.NewContext(app, set, nil)
	err := runCommand.ExecuteRunCommand(context)
	r.Nil(err)
}
