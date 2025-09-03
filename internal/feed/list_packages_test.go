package feed_test

import (
	"testing"

	"github.com/patrickhuber/go-cross/platform"
	"github.com/patrickhuber/go-di"
	"github.com/patrickhuber/wrangle/internal/bootstrap"
	"github.com/patrickhuber/wrangle/internal/feed"
	"github.com/patrickhuber/wrangle/internal/host"
	"github.com/stretchr/testify/require"
)

func TestListPackages(t *testing.T) {
	plat := platform.Linux
	h := host.NewTest(plat, nil, nil)
	defer h.Close()

	container := h.Container()

	// bootstrap needs to run to setup the user and system configurations
	bootstrapService, err := di.Resolve[bootstrap.Service](container)
	require.NoError(t, err)

	err = bootstrapService.Execute(&bootstrap.Request{
		Force: true,
	})
	require.NoError(t, err)

	svc, err := di.Resolve[feed.ListPackages](container)
	require.NoError(t, err)

	request := &feed.ListPackagesRequest{
		Names: []string{"test"},
	}

	response, err := svc.Execute(request)
	require.NoError(t, err)
	require.NotNil(t, response)
	require.Equal(t, 1, len(response.Items))
	require.Equal(t, "test", response.Items[0].Package)
	// we are just listing packages, not versions so no versions are returned
	// maybe we sould return the latest version?
	require.Equal(t, "1.0.0", response.Items[0].Latest)
}
