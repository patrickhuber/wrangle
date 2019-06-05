package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"fmt"
	"io"
	"io/ioutil"
)

type aes256GCMCrypter struct {
	key    []byte
	nonce  []byte
	cipher cipher.AEAD
}

func (crypter *aes256GCMCrypter) Encrypt(reader io.Reader, writer io.Writer) error {

	if _, err := io.ReadFull(rand.Reader, crypter.nonce); err != nil {
		return err
	}

	data, err := ioutil.ReadAll(reader)
	if err != nil {
		return err
	}

	ciphertext := crypter.cipher.Seal(nil, crypter.nonce, data, nil)

	_, err = writer.Write([]byte(ciphertext))
	return err
}

func (crypter *aes256GCMCrypter) Decrypt(reader io.Reader, writer io.Writer) error {
	ciphertext, err := ioutil.ReadAll(reader)
	if err != nil {
		return err
	}
	plaintext, err := crypter.cipher.Open(nil, crypter.nonce, ciphertext, nil)
	if err != nil {
		return err
	}
	_, err = writer.Write(plaintext)
	return err
}

// NewAES256GCMCrypter creates a new AES 256 GCM encryptor
func NewAES256GCMCrypter(key []byte, nonce []byte) (Crypter, error) {

	// check key length. Optionally pad with some text
	if len(key) != 16 && len(key) != 32 {
		return nil, fmt.Errorf("key length must be 16 or 32 characters")
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	if nonce == nil {
		nonce = make([]byte, aesgcm.NonceSize())
	}

	return &aes256GCMCrypter{
		key:    key,
		nonce:  nonce,
		cipher: aesgcm,
	}, nil
}
