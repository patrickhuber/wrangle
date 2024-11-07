package services_test

import (
	"testing"

	"github.com/patrickhuber/go-cross/fs"
	"github.com/patrickhuber/go-cross/platform"
	"github.com/patrickhuber/go-di"
	"github.com/patrickhuber/wrangle/internal/config"
	"github.com/patrickhuber/wrangle/internal/host"
	"github.com/patrickhuber/wrangle/internal/services"
	"github.com/stretchr/testify/require"
)

func TestSecrets(t *testing.T) {
	h := host.NewTest(platform.Linux, nil, nil)
	defer h.Close()

	fs, err := di.Resolve[fs.FS](h.Container())
	require.Nil(t, err)

	configuration, err := di.Resolve[services.Configuration](h.Container())
	require.Nil(t, err)

	// add in the memory store
	cfg := configuration.GlobalDefault()
	cfg.Spec.Stores = append(cfg.Spec.Stores, config.Store{Type: "memory", Name: "test"})

	globalConfigFilePath, err := configuration.DefaultGlobalConfigFilePath()
	require.NoError(t, err)

	err = config.WriteFile(fs, globalConfigFilePath, cfg)
	require.Nil(t, err)

	secrets, err := di.Resolve[services.Secret](h.Container())
	require.NoError(t, err)

	err = secrets.Set("test", "test", "test")
	require.NoError(t, err)

	value, err := secrets.Get("test", "test")
	require.NoError(t, err)
	require.Equal(t, "test", value)
}
