package crypto

import (
	"io"

	"golang.org/x/crypto/openpgp"

	_ "golang.org/x/crypto/ripemd160" // <-- required to run this unit test
)

type pgpDecryptor struct {
	entityList openpgp.EntityList
}

// NewPgpDecryptor creates a new pgp encryptor using the secret key ring
func NewPgpDecryptor(secretKeyRingReader io.Reader) (Decryptor, error) {
	// load the entity list
	entitylist, err := openpgp.ReadKeyRing(secretKeyRingReader)
	if err != nil {
		return nil, err
	}

	return &pgpDecryptor{
		entityList: entitylist}, nil
}

func (d *pgpDecryptor) Decrypt(reader io.Reader, writer io.Writer) error {

	// create the message details
	messageDetails, err := openpgp.ReadMessage(reader, d.entityList, nil, nil)
	if err != nil {
		return err
	}

	// copy the bytes to the file unencrypting along the way
	_, err = io.Copy(writer, messageDetails.UnverifiedBody)
	return err
}
