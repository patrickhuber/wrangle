package file

import (
	"testing"

	"github.com/patrickhuber/cli-mgr/config"
	"github.com/stretchr/testify/require"
)

func TestFileConfigStoreConfig(t *testing.T) {
	t.Run("CanMapNameAndPath", func(t *testing.T) {
		r := require.New(t)
		configSource := &config.ConfigSource{
			Name: "name",
			Params: map[string]string{
				"path": "/test",
			},
		}
		cfg, err := NewFileConfigStoreConfig(configSource)
		r.Nil(err)
		r.NotNil(cfg)
		r.Equal("name", cfg.Name)
		r.Equal("/test", cfg.Path)
	})
}
