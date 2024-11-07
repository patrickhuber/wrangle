package services_test

import (
	"testing"

	"github.com/patrickhuber/go-cross"
	"github.com/patrickhuber/go-cross/arch"
	"github.com/patrickhuber/go-cross/platform"
	"github.com/patrickhuber/go-log"
	"github.com/patrickhuber/wrangle/internal/config"
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
			target := cross.NewTest(test.platform, arch.AMD64)
			fs := target.FS()
			os := target.OS()
			env := target.Env()
			path := target.Path()
			log := log.Memory()
			configuration, err := services.NewConfiguration(os, env, fs, path, log)
			require.NoError(t, err)

			cfg := configuration.GlobalDefault()
			require.NoError(t, err)

			globalConfigPath, err := configuration.DefaultGlobalConfigFilePath()
			require.NoError(t, err)

			err = config.WriteFile(fs, globalConfigPath, cfg)
			require.NoError(t, err)

			err = fs.MkdirAll("/opt/wrangle/packages/test/1.0.0", 0775)
			require.NoError(t, err)

			err = fs.WriteFile("/opt/wrangle/packages/test/1.0.0/"+test.exeName, []byte{}, 0775)
			require.NoError(t, err)

			req := &services.ShimRequest{
				Shell:       "bash",
				Executables: []string{"/opt/wrangle/package/test/1.0.0/" + test.exeName},
			}

			shim := services.NewShim(fs, path, configuration, log)
			err = shim.Execute(req)
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
