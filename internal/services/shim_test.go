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
			h := host.NewTest(test.platform, nil, nil)

			container := h.Container()

			svc, err := di.Resolve[services.Shim](container)
			require.NoError(t, err)

			fs, err := di.Resolve[fs.FS](container)
			require.NoError(t, err)

			configuration, err := di.Resolve[services.Configuration](container)
			require.NoError(t, err)

			cfg := configuration.GlobalDefault()
			require.NoError(t, err)

			globalConfigPath := configuration.DefaultGlobalConfigFilePath()
			err = config.WriteFile(fs, globalConfigPath, cfg)
			require.NoError(t, err)

			err = fs.MkdirAll("/opt/wrangle/packages/test/1.0.0", 0664)
			require.NoError(t, err)

			err = fs.WriteFile("/opt/wrangle/packages/test/1.0.0/"+test.exeName, []byte{}, 0775)
			require.NoError(t, err)

			req := &services.ShimRequest{
				Shell:   "bash",
				Package: "test",
				Version: "1.0.0",
			}
			err = svc.Execute(req)
			if err != nil {
				t.Fatal(err)
			}

			_, err = fs.Stat("/opt/wrangle/bin/" + test.exeName)
			require.NoError(t, err)

			content, err := fs.ReadFile("/opt/wrangle/bin/" + test.exeName)
			require.NoError(t, err)
			require.NotEmpty(t, content)
		})
	}
}
