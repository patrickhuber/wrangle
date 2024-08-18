package feed_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/patrickhuber/go-log"
	"github.com/patrickhuber/wrangle/internal/config"
	"github.com/patrickhuber/wrangle/internal/feed"
	"github.com/patrickhuber/wrangle/internal/feed/memory"
)

func TestServiceFactory(t *testing.T) {
	logger := log.Memory()
	provider := memory.NewProvider(logger)
	factory := feed.NewServiceFactory(provider)
	svc, err := factory.Create(config.Feed{
		Name: "test",
		Type: memory.ProviderType,
	})
	require.NoError(t, err)
	require.NotNil(t, svc)
}
