package commands

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCommandDispatch(t *testing.T) {
	require := require.New(t)

	command := Process{ExecutableName: "go", Arguments: []string{"version"}}
	err := Dispatch(&command)
	require.Nil(err)
}
