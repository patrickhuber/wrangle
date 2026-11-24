package edit_test

import (
	"testing"

	"github.com/patrickhuber/go-cross"
	"github.com/patrickhuber/go-cross/arch"
	"github.com/patrickhuber/go-cross/platform"
	"github.com/patrickhuber/go-log"
	"github.com/patrickhuber/wrangle/internal/config"
	"github.com/patrickhuber/wrangle/internal/config/edit"
	"github.com/patrickhuber/wrangle/internal/fixtures"
	"github.com/stretchr/testify/require"
)

func TestConfigScopes(t *testing.T) {
	for _, plat := range []platform.Platform{platform.Linux, platform.Windows} {
		t.Run(plat.String(), func(t *testing.T) {
			target := cross.NewTest(plat, arch.AMD64)
			require.NoError(t, fixtures.Apply(target.OS(), target.FS(), target.Env()))
			path, fs := target.Path(), target.FS()
			work, err := target.OS().WorkingDirectory()
			require.NoError(t, err)
			project := path.Join(work, "project")
			child := path.Join(project, "child")
			require.NoError(t, fs.MkdirAll(child, 0755))
			require.NoError(t, target.OS().ChangeDirectory(child))
			root, err := config.GetRoot(target.Env(), path, plat)
			require.NoError(t, err)
			home, err := target.OS().Home()
			require.NoError(t, err)
			system := config.GetDefaultSystemConfigPath(path, root)
			user := config.GetDefaultUserConfigPath(path, home)
			local := path.Join(project, "project.wrangle.yml")
			for _, file := range []string{system, user, local} {
				require.NoError(t, fs.MkdirAll(path.Dir(file), 0755))
				require.NoError(t, config.WriteFile(fs, file, config.Config{ApiVersion: config.ApiVersion, Kind: config.Kind}))
			}
			service := edit.NewService(fs, target.OS(), path, log.Memory(), target.Env(), config.NewMockCliContext(nil), config.NewTestSystemDefaultProvider(path))
			check := func(file, name, version string) {
				t.Helper()
				cfg, err := config.ReadFile(fs, file)
				require.NoError(t, err)
				require.Contains(t, cfg.Spec.Packages, config.Package{Name: name, Version: version})
			}

			require.NoError(t, service.AddPackage(&edit.PackageRequest{Package: "nearest"}))
			check(local, "nearest", "latest")
			require.NoError(t, service.AddPackage(&edit.PackageRequest{Package: "user", Version: "1.2.3", Scope: edit.ScopeUser}))
			check(user, "user", "1.2.3")
			require.NoError(t, service.AddPackage(&edit.PackageRequest{Package: "system", Scope: edit.ScopeSystem}))
			check(system, "system", "latest")
			require.NoError(t, service.AddPackage(&edit.PackageRequest{Package: "local", Scope: edit.ScopeLocal}))
			check(local, "local", "latest")

			require.NoError(t, fs.Remove(local))
			require.NoError(t, service.AddPackage(&edit.PackageRequest{Package: "fallback-user"}))
			check(user, "fallback-user", "latest")
			require.NoError(t, fs.Remove(user))
			require.NoError(t, service.AddPackage(&edit.PackageRequest{Package: "fallback-system"}))
			check(system, "fallback-system", "latest")
			require.NoError(t, service.AddPackage(&edit.PackageRequest{Package: "new-local", Scope: edit.ScopeLocal}))
			check(path.Join(child, ".wrangle.toml"), "new-local", "latest")
		})
	}
}

func TestAmbiguousLocalConfig(t *testing.T) {
	target := cross.NewTest(platform.Linux, arch.AMD64)
	require.NoError(t, fixtures.Apply(target.OS(), target.FS(), target.Env()))
	work, err := target.OS().WorkingDirectory()
	require.NoError(t, err)
	for _, name := range []string{".wrangle.toml", ".wrangle.yml"} {
		require.NoError(t, config.WriteFile(target.FS(), target.Path().Join(work, name), config.Config{}))
	}
	service := edit.NewService(target.FS(), target.OS(), target.Path(), log.Memory(), target.Env(), config.NewMockCliContext(nil), config.NewTestSystemDefaultProvider(target.Path()))
	require.ErrorContains(t, service.AddPackage(&edit.PackageRequest{Package: "test"}), "multiple local config files")
	require.ErrorContains(t, edit.ValidateScope("invalid"), "invalid scope")
}
