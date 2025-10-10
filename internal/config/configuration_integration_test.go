package config_test

import (
	"testing"

	goconfig "github.com/patrickhuber/go-config"
	cross "github.com/patrickhuber/go-cross"
	"github.com/patrickhuber/go-cross/arch"
	"github.com/patrickhuber/go-cross/platform"
	"github.com/stretchr/testify/require"

	"github.com/patrickhuber/wrangle/internal/config"
)

func TestDefaultConfiguration_LoadsLocalConfigFromParent(t *testing.T) {
	t.Parallel()

	target := cross.NewTest(platform.Linux, arch.AMD64)
	fs := target.FS()
	env := target.Env()
	path := target.Path()
	os := target.OS()

	root := "/workspace"
	child := path.Join(root, "sub")

	err := fs.MkdirAll(child, 0o755)
	require.NoError(t, err)

	localConfigPath := path.Join(root, "project.wrangle.yml")
	localConfig := []byte("spec:\n  env:\n    TEST: TEST\n")
	err = fs.WriteFile(localConfigPath, localConfig, 0o644)
	require.NoError(t, err)

	defaultRoot, err := config.GetRoot(env, os.Platform())
	require.NoError(t, err)

	systemConfigPath := config.GetDefaultSystemConfigPath(path, defaultRoot)
	err = fs.MkdirAll(path.Dir(systemConfigPath), 0o755)
	require.NoError(t, err)
	err = config.WriteFile(fs, systemConfigPath, config.Config{ApiVersion: config.ApiVersion, Kind: config.Kind})
	require.NoError(t, err)

	home, err := os.Home()
	require.NoError(t, err)
	userConfigPath := config.GetDefaultUserConfigPath(path, home)
	err = fs.MkdirAll(path.Dir(userConfigPath), 0o755)
	require.NoError(t, err)
	err = config.WriteFile(fs, userConfigPath, config.Config{ApiVersion: config.ApiVersion, Kind: config.Kind})
	require.NoError(t, err)

	err = os.ChangeDirectory(child)
	require.NoError(t, err)

	resolver := goconfig.DefaultGlobResolver(fs, path)
	cfgSvc, err := config.NewDefaultConfiguration(
		env,
		fs,
		path,
		os,
		config.NewMockCliContext(nil),
		resolver,
		config.NewSystemDefaultProvider(path),
	)
	require.NoError(t, err)

	cfg, err := cfgSvc.Get()
	require.NoError(t, err)

	require.Equal(t, "TEST", cfg.Spec.Environment["TEST"])
}
