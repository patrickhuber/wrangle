package feed_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/patrickhuber/wrangle/pkg/config"
	"github.com/patrickhuber/wrangle/pkg/feed"
)

func TestManager(t *testing.T) {
	cfg := &config.Config{
		Feeds: []*config.Feed{
			{
				Name: "local",
			},
			{
				Name: "remote",
			},
		},
	}
	mgr := feed.NewManager(cfg)
	feeds := mgr.List()
	require.Equal(t, 2, len(feeds))
}
