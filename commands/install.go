package commands

import (
	"fmt"
	"strings"

	"github.com/patrickhuber/wrangle/tasks"

	"github.com/spf13/afero"

	"github.com/patrickhuber/wrangle/filepath"

	"github.com/patrickhuber/wrangle/filesystem"

	"github.com/patrickhuber/wrangle/config"
	"github.com/patrickhuber/wrangle/packages"
)

type install struct {
	platform     string
	packagesPath string
	fileSystem   filesystem.FsWrapper
	manager      packages.Manager
	loader       config.Loader
}

// Install defines an install package command
type Install interface {
	Execute(packageName string, packageVersion string) error
}

// NewInstall creates a new install package command
func NewInstall(
	platform string,
	packagesPath string,
	fileSystem filesystem.FsWrapper,
	manager packages.Manager,
	loader config.Loader) (Install, error) {
	if packagesPath == "" {
		return nil, fmt.Errorf("packages path can not be empty")
	}
	return &install{
			platform:     platform,
			packagesPath: packagesPath,
			fileSystem:   fileSystem,
			manager:      manager,
			loader:       loader},
		nil
}

func (cmd *install) Execute(packageName string, packageVersion string) error {
	// load the configuration for the package into the struct
	configPkg, err := cmd.findConfigPackage(packageName, packageVersion)
	if err != nil {
		return err
	}

	// turn the config into package object
	pkg, err := cmd.createPackageFromConfig(configPkg, cmd.platform)
	if err != nil {
		return err
	}
	return cmd.manager.Install(pkg)
}

func (cmd *install) findConfigPackage(packageName, packageVersion string) (*config.Package, error) {
	if strings.TrimSpace(packageName) == "" {
		return nil, fmt.Errorf("package name is required")
	}

	packagePath := filepath.Join(cmd.packagesPath, packageName)
	useLatestVersion := len(strings.TrimSpace(packageVersion)) == 0

	var packageManifestPath string
	if !useLatestVersion {
		packageManifestFileName := fmt.Sprintf("%s.%s.yml", packageName, packageVersion)
		packageManifestPath = filepath.Join(packagePath, packageVersion, packageManifestFileName)

	} else {
		files, err := afero.ReadDir(cmd.fileSystem, packagePath)
		if err != nil {
			return nil, err
		}
		for _, file := range files {
			if !file.IsDir() {
				continue
			}
			packageVersion = file.Name()
			break
		}
	}

	configPackage, err := cmd.loader.LoadPackage(packageManifestPath)
	if err != nil {
		return nil, err
	}
	return configPackage, nil
}

func (cmd *install) createPackageFromConfig(
	configPackage *config.Package, platformName string) (packages.Package, error) {

	taskList := make([]tasks.Task, 0)

	for _, configPlatform := range configPackage.Platforms {
		if configPlatform.Name != platformName {
			continue
		}
		for _, configTask := range configPlatform.Tasks {
			params := map[string]string{}
			for key, value := range configTask.Params {
				params[key] = value.(string)
			}
			task := tasks.NewTask(configTask.Name, configTask.Type, params)
			taskList = append(taskList, task)
		}
	}
	// create the package with the download and extract params set
	pkg := packages.New(
		configPackage.Name,
		configPackage.Version,
		taskList...)

	return pkg, nil
}
