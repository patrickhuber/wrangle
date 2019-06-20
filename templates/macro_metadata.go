package templates

import (
	"fmt"
	"strings"
)

// MacroMetadata represents an unstructured macro metadata
type MacroMetadata struct {
	Name   string
	Values []string
}

// ParseMacroMetadata parses marco metadata into a MacroMetadata structure
func ParseMacroMetadata(metadata string) (*MacroMetadata, error) {
	if IsMacroMetadata(metadata) {
		metadata = metadata[1:]
	} else {
		return nil, fmt.Errorf("invalid macro metadata")
	}
	splits := strings.SplitN(metadata, ":", -1)
	macroMetadata := &MacroMetadata{
		Name:   splits[0],
		Values: splits[1:],
	}
	return macroMetadata, nil
}

// IsMacroMetadata determiens if a given string is a macro metadata string
func IsMacroMetadata(metadata string) bool {
	return strings.HasPrefix(metadata, "@")
}
