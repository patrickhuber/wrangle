package templates

import (
	"fmt"
)

type echoMacro struct {
}

func (m *echoMacro) Run(metadata *MacroMetadata) (string, error) {
	if len(metadata.Values) < 1 {
		return "", fmt.Errorf("echo text is required as the first value in metadata")
	}
	text := metadata.Values[0]
	return text, nil
}

// NewEchoMacro returns a new echo macro
func NewEchoMacro() Macro {
	return &echoMacro{}
}
