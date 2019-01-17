package crypto

import (
	"io"

	"golang.org/x/crypto/openpgp"

	_ "golang.org/x/crypto/ripemd160" // <-- required to run this unit test
)

type pgpEncrypter struct {
	entityList openpgp.EntityList
}

// NewPgpEncrypter creates a new pgp encrypter using the given public key ring
func NewPgpEncrypter(publicKeyRingReader io.Reader) (Encrypter, error) {
	// load the entity list
	entitylist, err := openpgp.ReadKeyRing(publicKeyRingReader)
	if err != nil {
		return nil, err
	}

	return &pgpEncrypter{
		entityList: entitylist}, nil
}

func (e *pgpEncrypter) Encrypt(reader io.Reader, writer io.Writer) error {

	// create the encryption writer
	w, err := openpgp.Encrypt(writer, e.entityList, nil, nil, nil)
	if err != nil {
		return err
	}
	defer w.Close()

	// copy the input file to the output buffer
	_, err = io.Copy(w, reader)
	return err
}
