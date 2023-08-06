package services_test

import (
	"testing"

	"github.com/patrickhuber/go-di"
	"github.com/patrickhuber/go-xplat/filepath"
	"github.com/patrickhuber/go-xplat/fs"
	"github.com/patrickhuber/go-xplat/os"
	"github.com/patrickhuber/go-xplat/platform"
	"github.com/patrickhuber/wrangle/internal/services"
	"github.com/patrickhuber/wrangle/internal/setup"
	"github.com/patrickhuber/wrangle/pkg/config"
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v3"
)

func TestWindowsInstall(t *testing.T) {
	s := setup.NewTest(platform.Windows)
	testFileLocation := `C:\ProgramData\wrangle\packages\test\1.0.0\test-1.0.0-windows-amd64.exe`
	RunInstallTest(t, testFileLocation, s)
}

func TestLinuxInstall(t *testing.T) {
	s := setup.NewTest(platform.Linux)
	testFileLocation := "/opt/wrangle/packages/test/1.0.0/test-1.0.0-linux-amd64"
	RunInstallTest(t, testFileLocation, s)
}

func TestDarwinInstall(t *testing.T) {
	s := setup.NewTest(platform.Darwin)
	testFileLocation := "/opt/wrangle/packages/test/1.0.0/test-1.0.0-darwin-amd64"
	RunInstallTest(t, testFileLocation, s)
}

func RunInstallTest(t *testing.T,
	testFileLocation string,
	s setup.Setup) {
	defer s.Close()
	container := s.Container()

	fs, err := di.Resolve[fs.FS](container)
	require.Nil(t, err)

	install, err := di.Resolve[services.Install](container)
	require.Nil(t, err)

	path, err := di.Resolve[*filepath.Processor](container)
	require.Nil(t, err)

	opsys, err := di.Resolve[os.OS](container)
	require.Nil(t, err)

	reader, err := di.Resolve[config.Provider](container)
	require.Nil(t, err)

	cfg, err := reader.Get()
	require.Nil(t, err)

	globalConfigPath := path.Join(opsys.Home(), ".wrangle", "config.yml")
	cfgBytes, err := yaml.Marshal(cfg)
	require.Nil(t, err)

	err = fs.WriteFile(globalConfigPath, cfgBytes, 0644)
	require.Nil(t, err)

	req := &services.InstallRequest{
		Package: "test",
	}

	err = install.Execute(req)
	require.Nil(t, err)

	ok, err := fs.Exists(testFileLocation)
	require.Nil(t, err)
	require.True(t, ok, "file '%s' not found", testFileLocation)
}
