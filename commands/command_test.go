package commands

import (
	"testing"
)

func TestCommandDispatch(t *testing.T) {
	process := "go"
	arguments := []string{"version"}
	command := Process{ExecutableName: process, Arguments: arguments}
	err := Dispatch(&command)
	if err != nil {
		panic(err)
	}
}
