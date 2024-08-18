package services_test

import (
	"fmt"
	"testing"

	"github.com/patrickhuber/go-di"
	"github.com/patrickhuber/go-xplat/fs"
	"github.com/patrickhuber/go-xplat/platform"
	"github.com/patrickhuber/wrangle/internal/config"
	"github.com/patrickhuber/wrangle/internal/host"
	"github.com/patrickhuber/wrangle/internal/services"
	"github.com/stretchr/testify/require"
)

func TestInstall(t *testing.T) {
	type packageTest struct {
		name    string
		version string
	}
	type fileTest struct {
		platform platform.Platform
		file     string
	}
	files := []fileTest{
		{
			platform: platform.Windows,
			file:     `C:\ProgramData\wrangle\packages\test\1.0.0\test-1.0.0-windows-amd64.exe`,
		},
		{
			platform: platform.Linux,
			file:     "/opt/wrangle/packages/test/1.0.0/test-1.0.0-linux-amd64",
		},
		{
			platform: platform.Darwin,
			file:     "/opt/wrangle/packages/test/1.0.0/test-1.0.0-darwin-amd64",
		},
	}
	packages := []packageTest{
		{
			name:    "test",
			version: "latest",
		},
		{
			name:    "test",
			version: "1.0.0",
		},
	}
	for _, f := range files {
		for _, p := range packages {
			t.Run(fmt.Sprintf("%s_%s_%s", f.platform.String(), p.name, p.version), func(t *testing.T) {
				s := host.NewTest(f.platform, nil, nil)
				defer s.Close()
				RunInstallTest(t, f.file, p.name, p.version, s)
			})
		}
	}
}

func RunInstallTest(t *testing.T,
	testFileLocation string,
	packageName string,
	packageVersion string,
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
		Package: packageName,
		Version: packageVersion,
	}

	install, err := di.Resolve[services.Install](container)
	require.NoError(t, err)

	err = install.Execute(req)
	require.Nil(t, err)

	ok, err := fs.Exists(testFileLocation)
	require.Nil(t, err)
	require.True(t, ok, "file '%s' not found", testFileLocation)
}
