package shim_test

import (
	"testing"

	"github.com/patrickhuber/go-cross"
	"github.com/patrickhuber/go-cross/arch"
	"github.com/patrickhuber/go-cross/platform"
	"github.com/patrickhuber/go-log"
	"github.com/patrickhuber/go-shellhook"
	"github.com/patrickhuber/wrangle/internal/config"
	"github.com/patrickhuber/wrangle/internal/fixtures"
	"github.com/patrickhuber/wrangle/internal/global"
	"github.com/patrickhuber/wrangle/internal/shim"
	"github.com/stretchr/testify/require"
)

func TestShim(t *testing.T) {
	type test struct {
		name     string
		platform platform.Platform
		shell    string
		exeName  string
		expected string
	}
	tests := []test{
		{"file_name_only", platform.Linux, shellhook.Bash, "test", "test"},
		{"exe", platform.Windows, shellhook.Powershell, "test.exe", "test.ps1"},
		{"bat", platform.Windows, shellhook.Powershell, "test.bat", "test.ps1"},
		{"cmd", platform.Windows, shellhook.Powershell, "test.cmd", "test.ps1"},
		{"ps1", platform.Windows, shellhook.Powershell, "test.exe", "test.ps1"},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			target := cross.NewTest(test.platform, arch.AMD64)
			fs := target.FS()
			path := target.Path()
			log := log.Memory()
			env := target.Env()

			fixtures.Apply(target.OS(), fs, env)

			root, err := config.GetRoot(env, path, test.platform)
			require.NoError(t, err)

			binDirectory := config.GetDefaultBinPath(path, root)
			packagesDirectory := config.GetDefaultPackagesPath(path, root)

			cfg := config.Config{
				Spec: config.Spec{
					Environment: map[string]string{
						global.EnvBin:      binDirectory,
						global.EnvPackages: packagesDirectory,
					},
				},
			}
			configuration := config.NewMock(cfg)

			packageName := "test"
			packageVersion := "1.0.0"

			err = fs.MkdirAll(path.Join(packagesDirectory, packageName, packageVersion), 0775)
			require.NoError(t, err)

			err = fs.WriteFile(path.Join(packagesDirectory, packageName, packageVersion, test.exeName), []byte{}, 0775)
			require.NoError(t, err)

			req := &shim.Request{
				Shell: test.shell,
				Executables: []string{
					path.Join(packagesDirectory, packageName, packageVersion, test.exeName),
				},
			}

			shim := shim.NewService(fs, path, configuration, log)
			err = shim.Execute(req)
			if err != nil {
				t.Fatal(err)
			}

			_, err = fs.Stat(path.Join(binDirectory, test.expected))
			require.NoError(t, err)

			content, err := fs.ReadFile(path.Join(binDirectory, test.expected))
			require.NoError(t, err)
			require.NotEmpty(t, content)
		})
	}
}
