package processes

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestProcessDispatch(t *testing.T) {
	t.Run("CanRunGoVersion", func(t *testing.T) {
		r := require.New(t)
		stdOut := bytes.Buffer{}
		stdErr := bytes.Buffer{}
		stdIn := bytes.Buffer{}
		command := NewProcess("go", []string{"version"}, make(map[string]string), &stdOut, &stdErr, &stdIn)
		err := command.Dispatch()
		r.Nil(err)
	})

	t.Run("CanRunWithNilEnvironmentVariables", func(t *testing.T) {
		r := require.New(t)
		stdOut := bytes.Buffer{}
		stdErr := bytes.Buffer{}
		stdIn := bytes.Buffer{}
		command := NewProcess("go", []string{"version"}, make(map[string]string), &stdOut, &stdErr, &stdIn)
		err := command.Dispatch()
		r.Nil(err)
	})

	t.Run("WritesToStandardOut", func(t *testing.T) {
		r := require.New(t)
		stdOut := bytes.Buffer{}
		stdErr := bytes.Buffer{}
		stdIn := bytes.Buffer{}
		command := NewProcess("go", []string{"version"}, make(map[string]string), &stdOut, &stdErr, &stdIn)
		err := command.Dispatch()
		r.Nil(err)
		r.NotEmpty(stdOut.String())
		r.Empty(stdErr.String())
	})

	t.Run("WritesToStandardError", func(t *testing.T) {
		r := require.New(t)
		stdOut := bytes.Buffer{}
		stdErr := bytes.Buffer{}
		stdIn := bytes.Buffer{}
		command := NewProcess("go", []string{}, make(map[string]string), &stdOut, &stdErr, &stdIn)
		err := command.Dispatch()
		r.NotNil(err)
		r.NotEmpty(stdErr.String())
	})
}
