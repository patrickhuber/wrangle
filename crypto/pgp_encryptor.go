package crypto

import (
	"io"

	"golang.org/x/crypto/openpgp"

	_ "golang.org/x/crypto/ripemd160" // <-- required to run this unit test
)

type pgpEncryptor struct {
	entityList openpgp.EntityList
}

func NewPgpEncryptor(publicKeyRingReader io.Reader) (Encryptor, error) {
	// load the entity list
	entitylist, err := openpgp.ReadKeyRing(publicKeyRingReader)
	if err != nil {
		return nil, err
	}

	return &pgpEncryptor{
		entityList: entitylist}, nil
}

func (e *pgpEncryptor) Encrypt(reader io.Reader, writer io.Writer) error {

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
