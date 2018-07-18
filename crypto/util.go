package crypto

import "github.com/spf13/afero"

func EncryptFile(fs afero.Fs, encryptor Encryptor, plainTextFile string, encryptedFile string) error {
	plaintext, err := fs.Open(plainTextFile)
	if err != nil {
		return err
	}
	defer plaintext.Close()

	encrypted, err := fs.Create(encryptedFile)
	if err != nil {
		return err
	}
	defer encrypted.Close()

	return encryptor.Encrypt(plaintext, encrypted)
}
