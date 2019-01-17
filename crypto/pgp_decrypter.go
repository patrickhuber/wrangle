package crypto

import (
	"io"

	"golang.org/x/crypto/openpgp"

	_ "golang.org/x/crypto/ripemd160" // <-- required to run this unit test
)

type pgpDecrypter struct {
	entityList openpgp.EntityList
}

// NewPgpDecrypter creates a new pgp encrypter using the secret key ring
func NewPgpDecrypter(secretKeyRingReader io.Reader) (Decrypter, error) {
	// load the entity list
	entitylist, err := openpgp.ReadKeyRing(secretKeyRingReader)
	if err != nil {
		return nil, err
	}

	return &pgpDecrypter{
		entityList: entitylist}, nil
}

func (d *pgpDecrypter) Decrypt(reader io.Reader, writer io.Writer) error {

	// create the message details
	messageDetails, err := openpgp.ReadMessage(reader, d.entityList, nil, nil)
	if err != nil {
		return err
	}

	// copy the bytes to the file unencrypting along the way
	_, err = io.Copy(writer, messageDetails.UnverifiedBody)
	return err
}
