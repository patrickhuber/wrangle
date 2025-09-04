package fixtures

import (
	"github.com/patrickhuber/go-cross/env"
	"github.com/patrickhuber/go-cross/fs"
	"github.com/patrickhuber/go-cross/os"
	"github.com/patrickhuber/go-cross/platform"
)

func Apply(os os.OS, fileSystem fs.FS, environment env.Environment) error {
	if platform.IsWindows(os.Platform()) {
		return applyWindows(os, fileSystem, environment)
	} else {
		return applyPosix(os, fileSystem, environment)
	}
}

func applyWindows(os os.OS, fileSystem fs.FS, environment env.Environment) error {
	home, err := os.Home()
	if err != nil {
		return err
	}

	programData := "C:\\ProgramData"

	// Apply Windows-specific configurations
	windowsEnv := map[string]string{
		"USERPROFILE": home,
		"USERNAME":    "fake",
		"ProgramData": programData,
	}

	err = applyEnv(windowsEnv, environment)
	if err != nil {
		return err
	}

	working, err := os.WorkingDirectory()
	if err != nil {
		return err
	}

	directories := []string{
		home,
		working,
		programData,
	}

	for _, dir := range directories {
		err = fileSystem.MkdirAll(dir, 0755)
		if err != nil {
			return err
		}
	}
	return nil
}

func applyPosix(os os.OS, fileSystem fs.FS, environment env.Environment) error {
	home, err := os.Home()
	if err != nil {
		return err
	}

	// Apply Posix-specific configurations
	posixEnv := map[string]string{
		"HOME": home,
		"USER": "fake",
	}

	err = applyEnv(posixEnv, environment)
	if err != nil {
		return err
	}

	working, err := os.WorkingDirectory()
	if err != nil {
		return err
	}

	directories := []string{
		"/opt",
		home,
		working,
	}

	for _, dir := range directories {
		err = fileSystem.MkdirAll(dir, 0755)
		if err != nil {
			return err
		}
	}

	return nil
}

func applyEnv(vars map[string]string, env env.Environment) error {
	for k, v := range vars {
		if err := env.Set(k, v); err != nil {
			return err
		}
	}
	return nil
}
