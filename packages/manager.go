package packages

import (
	"fmt"
	"io"
	"net/http"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/patrickhuber/wrangle/archiver"
	fp "github.com/patrickhuber/wrangle/filepath"

	"github.com/pkg/errors"

	"github.com/patrickhuber/wrangle/filesystem"
)

type manager struct {
	fileSystem filesystem.FsWrapper
}

// Manager defines a manager interface
type Manager interface {
	Download(p Package) error
	Extract(p Package) error
	Link(p Package) error
}

// NewManager creates a new package manager
func NewManager(fileSystem filesystem.FsWrapper) Manager {
	return &manager{fileSystem: fileSystem}
}

func (m *manager) Download(p Package) error {

	if p.Download() == nil {
		return errors.New("package Download() is required")
	}

	// get the file data
	resp, err := http.Get(p.Download().URL())
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return fmt.Errorf("error downloading '%s'. http status code: '%d'. http status: '%s'",
			p.Download().URL(),
			resp.StatusCode,
			resp.Status)
	}

	// create the file
	file, err := m.fileSystem.Create(p.Download().OutPath())
	if err != nil {
		return err
	}
	defer file.Close()

	// Write the body to file
	_, err = io.Copy(file, resp.Body)

	return err
}

func isBinaryPackage(p Package) bool {
	return p.Extract() == nil
}

func (m *manager) Extract(p Package) error {

	if p.Download() == nil {
		return errors.New("package Download() is required")
	}
	if p.Extract() == nil {
		return errors.New("package Extract() is required")
	}

	// open the file for reading
	path := p.Download().OutPath()
	file, err := m.fileSystem.Open(path)

	if err != nil {
		return err
	}

	defer file.Close()

	extension := filepath.Ext(path)
	if strings.HasSuffix(path, ".tar.gz") {
		extension = ".tgz"
	}

	var a archiver.Archiver
	switch extension {
	case ".tgz":
		a = archiver.NewTargzArchiver(m.fileSystem)
		break
	case ".tar":
		a = archiver.NewTarArchiver(m.fileSystem)
		break
	case ".zip":
		a = archiver.NewZipArchiver(m.fileSystem)
		break
	default:
		return fmt.Errorf("unrecoginzed file extension '%s'", extension)
	}

	return a.Extract(file, p.Extract().Filter(), p.Extract().OutPath())
}

func isMatch(name string, filter string) (bool, error) {
	normalizedName := strings.Replace(name, "\\", "/", -1)

	if normalizedName == filter {
		return true, nil
	}
	return regexp.MatchString(filter, normalizedName)
}

func (m *manager) Link(p Package) error {
	// set the permissions of the package output
	if isBinaryPackage(p) {
		return m.postProcessFile(
			p.Alias(),
			p.Download().OutFolder(),
			p.Download().OutFile())
	}
	return m.postProcessFile(
		p.Alias(),
		p.Extract().OutFolder(),
		p.Extract().OutFile())
}

func (m *manager) postProcessFile(alias string, sourceFolder string, sourceFile string) error {
	sourcePath := fp.Join(sourceFolder, sourceFile)
	err := m.fileSystem.Chmod(sourcePath, 0755)
	if err != nil {
		return err
	}

	// the file needs to have a symlink with the alias name
	aliasPath := fp.Join(sourceFolder, alias)
	return m.fileSystem.Symlink(sourcePath, aliasPath)
}
