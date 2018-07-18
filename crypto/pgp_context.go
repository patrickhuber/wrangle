package crypto

import (
	"os/user"
	"path/filepath"
	"strings"

	"github.com/patrickhuber/wrangle/filesystem"
)

type pgpContext struct {
	publicKeyRing filesystem.FilePath
	secureKeyRing filesystem.FilePath
}

type PgpContext interface {
	PublicKeyRing() filesystem.FilePath
	SecureKeyRing() filesystem.FilePath
}

func NewPlatformPgpContext(platform string) (PgpContext, error) {

	// get user directory
	u, err := user.Current()
	if err != nil {
		return nil, err
	}

	gpgDirectory := ""

	if platform == "windows" {
		username := u.Username
		if strings.Contains(username, "\\") {
			split := strings.Split(username, "\\")
			username = split[1]
		}
		gpgDirectory = filepath.Join("c:/Users", username, "/AppData/Roaming/gnupg")
	} else {
		gpgDirectory = filepath.Join("/home", u.Username, ".gnupg")
	}

	return NewPgpContextFromFolder(gpgDirectory), nil
}

func NewPgpContextFromFolder(gnupgFolder string) PgpContext {
	return &pgpContext{
		publicKeyRing: filesystem.NewFilePathFromDirectoryAndFile(gnupgFolder, "pubring.gpg"),
		secureKeyRing: filesystem.NewFilePathFromDirectoryAndFile(gnupgFolder, "secring.gpg"),
	}
}

func NewPgpContext(publicKeyRingFullPath string, secretKeyRingFullPath string) PgpContext {
	return &pgpContext{
		publicKeyRing: filesystem.NewFilePathFromFullPath(publicKeyRingFullPath),
		secureKeyRing: filesystem.NewFilePathFromFullPath(secretKeyRingFullPath),
	}
}

func (ctx *pgpContext) PublicKeyRing() filesystem.FilePath {
	return ctx.publicKeyRing
}

func (ctx *pgpContext) SecureKeyRing() filesystem.FilePath {
	return ctx.secureKeyRing
}
