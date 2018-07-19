package crypto

import (
	"io"
)

// Decryptor decrypts a reader into a writer
type Decryptor interface {
	Decrypt(reader io.Reader, writer io.Writer) error
}

// Encryptor encrypts a reader into a writer
type Encryptor interface {
	Encrypt(reader io.Reader, writer io.Writer) error
}
