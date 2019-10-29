package settings

import (
	"path/filepath"

	"github.com/patrickhuber/wrangle/filesystem"
	"github.com/patrickhuber/wrangle/global"
)

type fsProvider struct {
	fs            filesystem.FileSystem
	platform      string
	homeDirectory string
}

// NewFsProvider defines a constructor for creating a file system based settings provider
func NewFsProvider(fs filesystem.FileSystem, platform string, homeDirectory string) Provider {
	return &fsProvider{
		fs:            fs,
		platform:      platform,
		homeDirectory: homeDirectory,
	}
}

func (provider *fsProvider) Initialize() (*Settings, error) {
	wrangleSettingsPath, err := provider.getSettingsPath()
	if err != nil {
		return nil, err
	}

	// TODO: Remove side effects from get. Should return hard coded defaults unless the settings file exists.
	// if the settings file directory doesn't exist, generate it
	err = provider.ensureSettingsPathExists(wrangleSettingsPath)
	if err != nil {
		return nil, err
	}

	// create the settings file if it doesn't exist
	wrangleSettingsFilePath := provider.getSettingsFilePath(wrangleSettingsPath)
	err = provider.ensureSettingsFileExists(wrangleSettingsFilePath)
	if err != nil {
		return nil, err
	}

	return provider.Get()
}

func (provider *fsProvider) Get() (*Settings, error) {
	wrangleSettingsPath, err := provider.getSettingsPath()
	if err != nil {
		return nil, err
	}

	ok, err := provider.checkSettingsPathExists(wrangleSettingsPath)
	if err != nil {
		return nil, err
	}
	if !ok {
		return provider.generateSettings(), nil
	}

	wrangleSettingsFilePath := provider.getSettingsFilePath(wrangleSettingsPath)

	ok, err = provider.checkSettingsFileExists(wrangleSettingsFilePath)
	if err != nil {
		return nil, err
	}
	if !ok {
		return provider.generateSettings(), nil
	}

	// load the settings file and return the settings
	file, err := provider.fs.Open(wrangleSettingsFilePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := NewReader(file)
	return reader.Read()
}

func (provider *fsProvider) getSettingsPath() (string, error) {
	wrangleSettingsPath := filepath.Join(provider.homeDirectory, ".wrangle")

	return wrangleSettingsPath, nil
}

func (provider *fsProvider) getSettingsFilePath(settingsPath string) string {
	return filepath.Join(settingsPath, "settings.yml")
}

func (provider *fsProvider) checkSettingsPathExists(filePath string) (bool, error) {
	return provider.fs.Exists(filePath)
}

func (provider *fsProvider) ensureSettingsPathExists(path string) error {
	ok, err := provider.checkSettingsPathExists(path)
	if err != nil {
		return err
	}
	if ok {
		return nil
	}
	return provider.fs.Mkdir(path, 0600)
}

func (provider *fsProvider) checkSettingsFileExists(filePath string) (bool, error) {
	return provider.fs.Exists(filePath)
}

func (provider *fsProvider) ensureSettingsFileExists(filePath string) error {
	ok, err := provider.checkSettingsFileExists(filePath)
	if err != nil {
		return err
	}
	if ok {
		return nil
	}
	s := provider.generateSettings()
	return provider.set(filePath, s)
}

func (provider *fsProvider) generateSettings() *Settings {
	if provider.platform == "windows" {
		return provider.generateSettingsWindows()
	}
	return provider.generateSettingsNix()
}

func (provider *fsProvider) generateSettingsWindows() *Settings {
	return &Settings{
		Feeds: []string{global.PackageFeedURL},
		Paths: &Paths{
			Bin:      DefaultWindowsBin,
			Root:     DefaultWindowsRoot,
			Packages: DefaultWindowsPackages,
		},
	}
}

func (provider *fsProvider) generateSettingsNix() *Settings {
	return &Settings{
		Feeds: []string{global.PackageFeedURL},
		Paths: &Paths{
			Bin:      DefaultNixBin,
			Root:     DefaultNixRoot,
			Packages: DefaultNixPackages,
		},
	}
}

func (provider *fsProvider) set(filePath string, s *Settings) error {
	file, err := provider.fs.Create(filePath)
	if err != nil {
		return nil
	}
	defer file.Close()
	writer := NewWriter(file)
	return writer.Write(s)
}

func (provider *fsProvider) Set(s *Settings) error {
	path, err := provider.getSettingsPath()
	if err != nil {
		return err
	}

	err = provider.ensureSettingsPathExists(path)
	if err != nil {
		return err
	}
	wrangleSettingsFilePath := provider.getSettingsFilePath(path)
	return provider.set(wrangleSettingsFilePath, s)
}
