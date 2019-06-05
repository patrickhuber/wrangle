package file

import "strings"

// MacroMetadata represents an unstructured macro metadata
type MacroMetadata struct {
	Name   string
	Values []string
}

// ParseMacroMetadata parses marco metadata into a MacroMetadata structure
func ParseMacroMetadata(metadata string) (*MacroMetadata, error) {
	if strings.HasPrefix(metadata, "@") {
		metadata = metadata[1:]
	}
	splits := strings.SplitN(metadata, ":", -1)
	macroMetadata := &MacroMetadata{
		Name:   splits[0],
		Values: splits[1:],
	}
	return macroMetadata, nil
}
