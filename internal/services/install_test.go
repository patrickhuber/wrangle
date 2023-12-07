package services_test

import (
	"testing"

	"github.com/patrickhuber/go-di"
	"github.com/patrickhuber/go-xplat/fs"
	"github.com/patrickhuber/go-xplat/platform"
	"github.com/patrickhuber/wrangle/internal/config"
	"github.com/patrickhuber/wrangle/internal/host"
	"github.com/patrickhuber/wrangle/internal/services"
	"github.com/stretchr/testify/require"
)

func TestWindowsInstall(t *testing.T) {
	s := host.NewTest(platform.Windows, nil, nil)
	defer s.Close()

	testFileLocation := `C:\ProgramData\wrangle\packages\test\1.0.0\test-1.0.0-windows-amd64.exe`
	RunInstallTest(t, testFileLocation, s)
}

func TestLinuxInstall(t *testing.T) {
	s := host.NewTest(platform.Linux, nil, nil)
	defer s.Close()

	testFileLocation := "/opt/wrangle/packages/test/1.0.0/test-1.0.0-linux-amd64"
	RunInstallTest(t, testFileLocation, s)
}

func TestDarwinInstall(t *testing.T) {
	s := host.NewTest(platform.Darwin, nil, nil)
	defer s.Close()

	testFileLocation := "/opt/wrangle/packages/test/1.0.0/test-1.0.0-darwin-amd64"
	RunInstallTest(t, testFileLocation, s)
}

func RunInstallTest(t *testing.T,
	testFileLocation string,
	s host.Host) {
	container := s.Container()

	fs, err := di.Resolve[fs.FS](container)
	require.NoError(t, err)

	configuration, err := di.Resolve[services.Configuration](container)
	require.NoError(t, err)

	cfg := configuration.GlobalDefault()
	require.NoError(t, err)

	globalConfigPath := configuration.DefaultGlobalConfigFilePath()
	err = config.WriteFile(fs, globalConfigPath, cfg)
	require.NoError(t, err)

	req := &services.InstallRequest{
		Package: "test",
	}

	install, err := di.Resolve[services.Install](container)
	require.NoError(t, err)

	err = install.Execute(req)
	require.Nil(t, err)

	ok, err := fs.Exists(testFileLocation)
	require.Nil(t, err)
	require.True(t, ok, "file '%s' not found", testFileLocation)
}
