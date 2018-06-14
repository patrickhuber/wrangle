package processes

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestProcessDispatch(t *testing.T) {
	t.Run("CanRunGoVersion", func(t *testing.T) {
		r := require.New(t)

		command := NewProcess("go", []string{"version"}, make(map[string]string))
		err := command.Dispatch()
		r.Nil(err)
	})

	t.Run("CanRunWithNilEnvironmentVariables", func(t *testing.T) {
		r := require.New(t)
		command := NewProcess("go", []string{"version"}, nil)
		err := command.Dispatch()
		r.Nil(err)
	})
}
