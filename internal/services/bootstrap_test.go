package services_test

import (
	"testing"

	"github.com/patrickhuber/go-di"
	"github.com/patrickhuber/go-xplat/filepath"
	"github.com/patrickhuber/go-xplat/fs"
	"github.com/patrickhuber/go-xplat/os"
	"github.com/patrickhuber/go-xplat/platform"
	"github.com/patrickhuber/wrangle/internal/host"
	"github.com/patrickhuber/wrangle/internal/services"
	"github.com/stretchr/testify/require"
)

func TestLinuxBootstrap(t *testing.T) {
	s := host.NewTest(platform.Linux, nil, nil)
	wrangleFileLocation := "/opt/wrangle/packages/wrangle/1.0.0/wrangle-1.0.0-linux-amd64"
	shimFileLocation := "/opt/wrangle/packages/shim/1.0.0/shim-1.0.0-linux-amd64"
	RunBootstrapTest(t, s, wrangleFileLocation, shimFileLocation)
}

func TestDarwinBootstrap(t *testing.T) {
	s := host.NewTest(platform.Darwin, nil, nil)
	wrangleFileLocation := "/opt/wrangle/packages/wrangle/1.0.0/wrangle-1.0.0-darwin-amd64"
	shimFileLocation := "/opt/wrangle/packages/shim/1.0.0/shim-1.0.0-darwin-amd64"
	RunBootstrapTest(t, s, wrangleFileLocation, shimFileLocation)
}

func TestWindowsBootstrap(t *testing.T) {
	s := host.NewTest(platform.Windows, nil, nil)
	wrangleFileLocation := "C:/ProgramData/wrangle/packages/wrangle/1.0.0/wrangle-1.0.0-windows-amd64.exe"
	shimFileLocation := "C:/ProgramData/wrangle/packages/shim/1.0.0/shim-1.0.0-windows-amd64.exe"
	RunBootstrapTest(t, s, wrangleFileLocation, shimFileLocation)
}

func RunBootstrapTest(t *testing.T,
	s host.Host,
	wrangleFileLocation string,
	shimFileLocation string) {
	defer s.Close()
	container := s.Container()

	bootstrap, err := di.Resolve[services.Bootstrap](container)
	require.Nil(t, err)

	opsys, err := di.Resolve[os.OS](container)
	require.Nil(t, err)

	path, err := di.Resolve[*filepath.Processor](container)
	require.Nil(t, err)

	globalConfigFile := path.Join(opsys.Home(), ".wrangle", "config.yml")
	req := &services.BootstrapRequest{
		ApplicationName: "wrangle",
	}
	err = bootstrap.Execute(req)
	require.Nil(t, err)

	fs, err := di.Resolve[fs.FS](container)
	require.Nil(t, err)

	ok, err := fs.Exists(globalConfigFile)
	require.Nil(t, err)
	require.True(t, ok)

	ok, err = fs.Exists(wrangleFileLocation)
	require.Nil(t, err)
	require.True(t, ok)

	ok, err = fs.Exists(shimFileLocation)
	require.Nil(t, err)
	require.True(t, ok)
}
