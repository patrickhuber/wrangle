package commands

import (
	"testing"

	"github.com/spf13/afero"
)

func TestEnvCommand(t *testing.T) {
	t.Run("CanRunCommand", func(t *testing.T) {
		fileSystem := afero.NewMemMapFs()
		cmd := NewEnvCommand(fileSystem)
		runCommandParams := NewRunCommandParams("", "", "")
		cmd.ExecuteCommand(runCommandParams)
	})
}
