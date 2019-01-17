package crypto

import (
	"io"
)

// Decrypter decrypts a reader into a writer
type Decrypter interface {
	Decrypt(reader io.Reader, writer io.Writer) error
}

// Encrypter encrypts a reader into a writer
type Encrypter interface {
	Encrypt(reader io.Reader, writer io.Writer) error
}

type Crypter interface {
	Decrypter
	Encrypter
}
