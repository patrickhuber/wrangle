package config

import (
	"fmt"
	"path"

	"github.com/patrickhuber/go-cross/env"
	"github.com/patrickhuber/go-cross/filepath"
	"github.com/patrickhuber/go-cross/platform"
)

func GetRoot(env env.Environment, plat platform.Platform) (string, error) {
	root := "/opt/wrangle"
	if platform.IsWindows(plat) {
		programData := env.Get("ProgramData")
		if len(programData) == 0 {
			return "", fmt.Errorf("failed to get ProgramData environment variable")
		}
		root = path.Join(programData, "wrangle")
	}
	return root, nil
}

func GetAppName(appNameBase string, plat platform.Platform) (string, error) {
	appName := appNameBase
	if platform.IsWindows(plat) {
		appName = appName + ".exe"
	}
	return appName, nil
}

func GetDefaultSystemConfigPath(path filepath.Provider, root string) string {
	return path.Join(root, "config", "config.yml")
}

func GetDefaultUserConfigPath(path filepath.Provider, home string) string {
	return path.Join(home, ".wrangle", "config.yml")
}

func GetDefaultBinPath(path filepath.Provider, root string) string {
	return path.Join(root, "bin")
}

func GetDefaultPackagesPath(path filepath.Provider, root string) string {
	return path.Join(root, "packages")
}
