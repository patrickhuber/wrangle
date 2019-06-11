package templates

import "fmt"

type encryptionMacro struct {
}

func (m *encryptionMacro) Run(metadata *MacroMetadata) (string, error) {
	if len(metadata.Values) < 1 {
		return "", fmt.Errorf("encryption algorithm is required as first value in metadata")
	}
	algorithm := metadata.Values[0]
	var macro Macro
	switch algorithm {
	case "AES256":
		macro = NewAES256Macro()
	default:
		return "", fmt.Errorf("unable to find algorithm '%s'", algorithm)
	}
	return macro.Run(metadata)
}

func NewEncryptionMacro() Macro {
	return &encryptionMacro{}
}
