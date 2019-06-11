package templates

import (
	"bytes"
	"encoding/base64"
	"fmt"

	"github.com/patrickhuber/wrangle/crypto"
)

type aes256Macro struct {
}

func (m *aes256Macro) Run(metadata *MacroMetadata) (string, error) {
	if len(metadata.Values) != 4 {
		return "", fmt.Errorf("Invalid metadata values, expected AES256, Nonce, Key and Cipher Text")
	}

	algorithm := metadata.Values[0]
	if algorithm != "AES256" {
		return "", fmt.Errorf("Invalid metadata values, expected first value element to be 'AES256', found '%s'", algorithm)
	}

	key := []byte(metadata.Values[1])

	nonce, err := base64.StdEncoding.DecodeString(metadata.Values[2])
	if err != nil {
		return "", err
	}

	cipherText, err := base64.StdEncoding.DecodeString(metadata.Values[3])
	if err != nil {
		return "", err
	}

	crypter, err := crypto.NewAES256GCMCrypter(key, nonce)
	if err != nil {
		return "", err
	}

	reader := bytes.NewBuffer(cipherText)
	writer := &bytes.Buffer{}
	err = crypter.Decrypt(reader, writer)
	if err != nil {
		return "", err
	}

	return writer.String(), nil
}

// NewAES256Macro returns a new instance of the AES256 macro
func NewAES256Macro() Macro {
	return &aes256Macro{}
}
