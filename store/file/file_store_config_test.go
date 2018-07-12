package file

import (
	"testing"

	"github.com/patrickhuber/wrangle/config"
	"github.com/stretchr/testify/require"
)

func TestFileStoreConfig(t *testing.T) {
	t.Run("CanMapNameAndPath", func(t *testing.T) {
		r := require.New(t)
		configSource := &config.Store{
			Name: "name",
			Params: map[string]string{
				"path": "/test",
			},
		}
		cfg, err := NewFileStoreConfig(configSource)
		r.Nil(err)
		r.NotNil(cfg)
		r.Equal("name", cfg.Name)
		r.Equal("/test", cfg.Path)
	})
}
