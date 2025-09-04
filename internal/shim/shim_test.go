package shim_test

import (
	"testing"

	"github.com/patrickhuber/go-cross"
	"github.com/patrickhuber/go-cross/arch"
	"github.com/patrickhuber/go-cross/platform"
	"github.com/patrickhuber/go-log"
	"github.com/patrickhuber/wrangle/internal/config"
	"github.com/patrickhuber/wrangle/internal/global"
	"github.com/patrickhuber/wrangle/internal/shim"
	"github.com/stretchr/testify/require"
)

func TestShim(t *testing.T) {
	type test struct {
		name     string
		platform platform.Platform
		exeName  string
	}
	tests := []test{
		{"file_name_only", platform.Linux, "test"},
		{"exe", platform.Linux, "test.exe"},
		{"bat", platform.Linux, "test.bat"},
		{"cmd", platform.Linux, "test.cmd"},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			target := cross.NewTest(test.platform, arch.AMD64)
			fs := target.FS()
			path := target.Path()
			log := log.Memory()

			cfg := config.Config{
				Spec: config.Spec{
					Environment: map[string]string{
						global.EnvBin:      "/opt/wrangle/bin",
						global.EnvPackages: "/opt/wrangle/packages",
					},
				},
			}
			configuration := config.NewMock(cfg)

			err := fs.MkdirAll(path.Join(cfg.Spec.Environment[global.EnvPackages], "test", "1.0.0"), 0775)
			require.NoError(t, err)

			err = fs.WriteFile(path.Join(cfg.Spec.Environment[global.EnvPackages], "test", "1.0.0", test.exeName), []byte{}, 0775)
			require.NoError(t, err)

			req := &shim.Request{
				Shell: "bash",
				Executables: []string{
					path.Join(cfg.Spec.Environment[global.EnvPackages], "test", "1.0.0", test.exeName),
				},
			}

			shim := shim.NewService(fs, path, configuration, log)
			err = shim.Execute(req)
			if err != nil {
				t.Fatal(err)
			}

			_, err = fs.Stat(path.Join(cfg.Spec.Environment[global.EnvBin], test.exeName))
			require.NoError(t, err)

			content, err := fs.ReadFile(path.Join(cfg.Spec.Environment[global.EnvBin], test.exeName))
			require.NoError(t, err)
			require.NotEmpty(t, content)
		})
	}
}
