package crypto

import (
	"io"
)

type Decryptor interface {
	Decrypt(reader io.Reader, writer io.Writer) error
}

type Encryptor interface {
	Encrypt(reader io.Reader, writer io.Writer) error
}
