package bootstrap_test

import (
	"testing"

	"github.com/patrickhuber/go-cross/fs"
	"github.com/patrickhuber/go-cross/platform"
	"github.com/patrickhuber/go-di"
	"github.com/patrickhuber/wrangle/internal/bootstrap"
	"github.com/patrickhuber/wrangle/internal/config"
	"github.com/patrickhuber/wrangle/internal/host"
	"github.com/stretchr/testify/require"
)

func TestBootstrap(t *testing.T) {
	type test struct {
		name                 string
		plat                 platform.Platform
		binDirectory         string
		globalConfigFile     string
		versionedPackagePath string
	}
	tests := []test{
		{
			name:                 "Linux",
			plat:                 platform.Linux,
			globalConfigFile:     "/opt/wrangle/config/config.yml",
			binDirectory:         "/opt/wrangle/bin",
			versionedPackagePath: "/opt/wrangle/packages/wrangle/1.0.0/wrangle",
		},
		{
			name:                 "Darwin",
			plat:                 platform.Darwin,
			globalConfigFile:     "/opt/wrangle/config/config.yml",
			binDirectory:         "/opt/wrangle/bin",
			versionedPackagePath: "/opt/wrangle/packages/wrangle/1.0.0/wrangle",
		},
		{
			name:                 "Windows",
			plat:                 platform.Windows,
			globalConfigFile:     "C:/ProgramData/wrangle/config/config.yml",
			binDirectory:         "C:/ProgramData/wrangle/bin",
			versionedPackagePath: "C:/ProgramData/wrangle/packages/wrangle/1.0.0/wrangle.exe",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			RunBootstrapTest(t, test.plat, test.globalConfigFile, test.versionedPackagePath, test.binDirectory)
		})
	}
}

func RunBootstrapTest(
	t *testing.T,
	plat platform.Platform,
	globalConfigFile string,
	versionedPackagePath string,
	binDirectory string) {

	host := host.NewTest(plat, nil, nil)
	container := host.Container()
	defer func() {
		err := host.Close()
		require.NoError(t, err)
	}()

	req := &bootstrap.Request{
		Force: true,
	}

	bootstrap, err := di.Resolve[bootstrap.Service](container)
	require.NoError(t, err)

	err = bootstrap.Execute(req)
	require.NoError(t, err)

	fs, err := di.Resolve[fs.FS](container)
	require.NoError(t, err)

	ok, err := fs.Exists(globalConfigFile)
	require.NoError(t, err)
	require.True(t, ok, globalConfigFile+" does not exist")

	ok, err = fs.Exists(versionedPackagePath)
	require.NoError(t, err)
	require.True(t, ok, "versioned package path: "+versionedPackagePath+" does not exist")

	ok, err = fs.Exists(binDirectory)
	require.NoError(t, err)
	require.True(t, ok, binDirectory+" does not exist")
}

type FakeCliContext struct {
	stringMap map[string]string
}

func NewFakeCliContext(m map[string]string) config.CliContext {
	return &FakeCliContext{
		stringMap: m,
	}
}

func (f FakeCliContext) String(key string) string {
	return f.stringMap[key]
}

func (f FakeCliContext) IsSet(key string) bool {
	_, ok := f.stringMap[key]
	return ok
}
