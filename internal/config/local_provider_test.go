package config_test

import (
	"testing"

	"github.com/patrickhuber/go-xplat/filepath"
	"github.com/patrickhuber/go-xplat/fs"
	"github.com/patrickhuber/go-xplat/os"
	"github.com/patrickhuber/wrangle/internal/config"
	"github.com/stretchr/testify/require"
)

func TestLocalProvider(t *testing.T) {
	var opsys os.OS
	var fileSystem fs.FS
	var path *filepath.Processor
	var localDefault config.Config
	var local config.LocalProvider

	setup := func() {
		localDefault = config.NewLocalDefault()
		opsys = os.NewMock()
		fileSystem = fs.NewMemory()
		path = filepath.NewProcessorWithOS(opsys)

		pwd, err := opsys.WorkingDirectory()
		require.NoError(t, err)

		err = fileSystem.MkdirAll(pwd, 0644)
		require.NoError(t, err)

		local = config.NewLocalProvider(localDefault, opsys, fileSystem, path)
	}

	t.Run("no local config", func(t *testing.T) {
		setup()
		cfgs, err := local.Get()
		require.NoError(t, err)
		require.Equal(t, 0, len(cfgs))
	})

	t.Run("one local config", func(t *testing.T) {
		setup()

		pwd, err := opsys.WorkingDirectory()
		require.NoError(t, err)

		config.NewFile(fileSystem, path.Join(pwd, ".wrangle.yml")).Write(localDefault)

		cfgs, err := local.Get()
		require.NoError(t, err)
		require.Equal(t, 1, len(cfgs))
	})
	t.Run("parent and local config", func(t *testing.T) {
		setup()

		pwd, err := opsys.WorkingDirectory()
		require.NoError(t, err)

		config.NewFile(fileSystem, path.Join(pwd, ".wrangle.yml")).Write(localDefault)
		config.NewFile(fileSystem, path.Join(path.Dir(pwd), ".wrangle.yml")).Write(localDefault)

		cfgs, err := local.Get()
		require.NoError(t, err)
		require.Equal(t, 2, len(cfgs))
	})
}
