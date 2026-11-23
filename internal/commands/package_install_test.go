package commands_test

import (
	"errors"
	"testing"

	"github.com/patrickhuber/go-cross/filepath"
	"github.com/patrickhuber/go-cross/fs"
	"github.com/patrickhuber/go-cross/os"
	"github.com/patrickhuber/go-cross/platform"
	"github.com/patrickhuber/go-di"
	"github.com/patrickhuber/wrangle/internal/add"
	"github.com/patrickhuber/wrangle/internal/bootstrap"
	"github.com/patrickhuber/wrangle/internal/commands"
	"github.com/patrickhuber/wrangle/internal/config"
	"github.com/patrickhuber/wrangle/internal/host"
	"github.com/patrickhuber/wrangle/internal/install"
	"github.com/stretchr/testify/require"
)

type fakeInstaller struct {
	called bool
	err    error
}

func TestInstallAlreadyInstalledWithScope(t *testing.T) {
	h := host.NewTest(platform.Linux, nil, nil)
	defer h.Close()
	container := h.Container()
	filesystem, err := di.Resolve[fs.FS](container)
	require.NoError(t, err)
	path, err := di.Resolve[filepath.Provider](container)
	require.NoError(t, err)
	opsys, err := di.Resolve[os.OS](container)
	require.NoError(t, err)
	work, err := opsys.WorkingDirectory()
	require.NoError(t, err)
	local := path.Join(work, ".wrangle.toml")
	require.NoError(t, config.WriteFile(filesystem, local, config.Config{ApiVersion: config.ApiVersion, Kind: config.Kind}))
	before, err := filesystem.ReadFile(local)
	require.NoError(t, err)

	bootstrapService, err := di.Resolve[bootstrap.Service](container)
	require.NoError(t, err)
	require.NoError(t, bootstrapService.Execute(&bootstrap.Request{}))
	after, err := filesystem.ReadFile(local)
	require.NoError(t, err)
	require.Equal(t, before, after)

	cmd := &commands.PackageInstallCommand{Options: commands.PackageInstallOptions{Package: "wrangle"}}
	require.NoError(t, di.Inject(container, cmd))
	require.NoError(t, cmd.Execute())
	after, err = filesystem.ReadFile(local)
	require.NoError(t, err)
	require.Equal(t, before, after)

	cmd.Options.Scope = "local"
	require.NoError(t, cmd.Execute())
	cfg, err := config.ReadFile(filesystem, local)
	require.NoError(t, err)
	require.Equal(t, []config.Package{{Name: "wrangle", Version: "latest"}}, cfg.Spec.Packages)
}

func (f *fakeInstaller) Execute(_ *install.Request) error {
	f.called = true
	return f.err
}

type fakeAdder struct {
	request *add.Request
}

func (f *fakeAdder) Execute(r *add.Request) error {
	f.request = r
	return nil
}

func TestInstallConfigScope(t *testing.T) {
	for _, tc := range []struct {
		name        string
		scope       string
		version     string
		installErr  error
		wantInstall bool
		wantAdd     bool
		wantError   bool
	}{
		{name: "default does not edit", wantInstall: true},
		{name: "scoped latest", scope: "user", wantInstall: true, wantAdd: true},
		{name: "scoped explicit version", scope: "local", version: "1.2.3", wantInstall: true, wantAdd: true},
		{name: "failed install does not edit", scope: "system", installErr: errors.New("install failed"), wantInstall: true, wantError: true},
		{name: "invalid scope does not install", scope: "invalid", wantError: true},
	} {
		t.Run(tc.name, func(t *testing.T) {
			installer, adder := &fakeInstaller{err: tc.installErr}, &fakeAdder{}
			cmd := commands.PackageInstallCommand{
				Install: installer, Add: adder,
				Options: commands.PackageInstallOptions{Package: "test", Version: tc.version, Scope: tc.scope},
			}
			err := cmd.Execute()
			if tc.wantError {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
			require.Equal(t, tc.wantInstall, installer.called)
			if tc.wantAdd {
				require.Equal(t, &add.Request{Package: "test", Version: tc.version, Scope: tc.scope}, adder.request)
			} else {
				require.Nil(t, adder.request)
			}
		})
	}
}
