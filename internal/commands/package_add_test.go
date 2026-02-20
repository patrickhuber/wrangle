package commands_test

import (
	"testing"

	"github.com/patrickhuber/go-cross/fs"
	"github.com/patrickhuber/go-cross/os"
	"github.com/patrickhuber/go-cross/platform"
	"github.com/patrickhuber/go-di"
	"github.com/patrickhuber/wrangle/internal/bootstrap"
	"github.com/patrickhuber/wrangle/internal/commands"
	"github.com/patrickhuber/wrangle/internal/config"
	"github.com/patrickhuber/wrangle/internal/global"
	"github.com/patrickhuber/wrangle/internal/host"
	"github.com/stretchr/testify/require"
)

func TestPackageAdd(t *testing.T) {
	t.Run("adds package to local config", func(t *testing.T) {
		h := host.NewTest(platform.Linux, nil, []string{})
		container := h.Container()

		bootstrapService, err := di.Resolve[bootstrap.Service](container)
		require.NoError(t, err)

		err = bootstrapService.Execute(&bootstrap.Request{})
		require.NoError(t, err)

		cmd := &commands.PackageAddCommand{
			Options: commands.PackageAddOptions{
				Package: "test",
				Version: "1.0.0",
			},
		}

		err = di.Inject(container, cmd)
		require.NoError(t, err)

		err = cmd.Execute()
		require.NoError(t, err)

		// verify the local config file was created and contains the package
		opsys, err := di.Resolve[os.OS](container)
		require.NoError(t, err)

		workDir, err := opsys.WorkingDirectory()
		require.NoError(t, err)

		filesystem, err := di.Resolve[fs.FS](container)
		require.NoError(t, err)

		path, err := di.Resolve[config.Service](container)
		require.NoError(t, err)
		require.NotNil(t, path)

		localFile := workDir + "/" + global.LocalConfigurationFileName
		exists, err := filesystem.Exists(localFile)
		require.NoError(t, err)
		require.True(t, exists, "local config file should exist")

		cfg, err := config.ReadFile(filesystem, localFile)
		require.NoError(t, err)
		require.Len(t, cfg.Spec.Packages, 1)
		require.Equal(t, "test", cfg.Spec.Packages[0].Name)
		require.Equal(t, "1.0.0", cfg.Spec.Packages[0].Version)
	})

	t.Run("updates existing package version", func(t *testing.T) {
		h := host.NewTest(platform.Linux, nil, []string{})
		container := h.Container()

		bootstrapService, err := di.Resolve[bootstrap.Service](container)
		require.NoError(t, err)

		err = bootstrapService.Execute(&bootstrap.Request{})
		require.NoError(t, err)

		// add first version
		cmd := &commands.PackageAddCommand{
			Options: commands.PackageAddOptions{
				Package: "test",
				Version: "1.0.0",
			},
		}

		err = di.Inject(container, cmd)
		require.NoError(t, err)

		err = cmd.Execute()
		require.NoError(t, err)

		// update to new version
		cmd2 := &commands.PackageAddCommand{
			Options: commands.PackageAddOptions{
				Package: "test",
				Version: "2.0.0",
			},
		}

		err = di.Inject(container, cmd2)
		require.NoError(t, err)

		err = cmd2.Execute()
		require.NoError(t, err)

		// verify the version was updated, not duplicated
		opsys, err := di.Resolve[os.OS](container)
		require.NoError(t, err)

		workDir, err := opsys.WorkingDirectory()
		require.NoError(t, err)

		filesystem, err := di.Resolve[fs.FS](container)
		require.NoError(t, err)

		localFile := workDir + "/" + global.LocalConfigurationFileName
		cfg, err := config.ReadFile(filesystem, localFile)
		require.NoError(t, err)
		require.Len(t, cfg.Spec.Packages, 1)
		require.Equal(t, "test", cfg.Spec.Packages[0].Name)
		require.Equal(t, "2.0.0", cfg.Spec.Packages[0].Version)
	})

	t.Run("adds multiple packages", func(t *testing.T) {
		h := host.NewTest(platform.Linux, nil, []string{})
		container := h.Container()

		bootstrapService, err := di.Resolve[bootstrap.Service](container)
		require.NoError(t, err)

		err = bootstrapService.Execute(&bootstrap.Request{})
		require.NoError(t, err)

		// add first package
		cmd := &commands.PackageAddCommand{
			Options: commands.PackageAddOptions{
				Package: "test",
				Version: "1.0.0",
			},
		}

		err = di.Inject(container, cmd)
		require.NoError(t, err)

		err = cmd.Execute()
		require.NoError(t, err)

		// add second package
		cmd2 := &commands.PackageAddCommand{
			Options: commands.PackageAddOptions{
				Package: "wrangle",
				Version: "0.9.0",
			},
		}

		err = di.Inject(container, cmd2)
		require.NoError(t, err)

		err = cmd2.Execute()
		require.NoError(t, err)

		// verify both packages exist
		opsys, err := di.Resolve[os.OS](container)
		require.NoError(t, err)

		workDir, err := opsys.WorkingDirectory()
		require.NoError(t, err)

		filesystem, err := di.Resolve[fs.FS](container)
		require.NoError(t, err)

		localFile := workDir + "/" + global.LocalConfigurationFileName
		cfg, err := config.ReadFile(filesystem, localFile)
		require.NoError(t, err)
		require.Len(t, cfg.Spec.Packages, 2)
		require.Equal(t, "test", cfg.Spec.Packages[0].Name)
		require.Equal(t, "wrangle", cfg.Spec.Packages[1].Name)
	})
}
