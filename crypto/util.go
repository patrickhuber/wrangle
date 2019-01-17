package crypto

import "github.com/spf13/afero"

// EncryptFile encrypts a file for the given filesystem, encrypter and files
func EncryptFile(fs afero.Fs, encrypter Encrypter, plainTextFile string, encryptedFile string) error {
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

	return encrypter.Encrypt(plaintext, encrypted)
}
