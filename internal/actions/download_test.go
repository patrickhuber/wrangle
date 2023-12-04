package actions_test

import (
	"net/http/httptest"
	"testing"

	"github.com/patrickhuber/go-di"
	"github.com/patrickhuber/go-xplat/filepath"
	"github.com/patrickhuber/go-xplat/fs"
	"github.com/patrickhuber/go-xplat/platform"
	"github.com/stretchr/testify/require"

	"github.com/patrickhuber/wrangle/internal/actions"
	"github.com/patrickhuber/wrangle/internal/host"
	"github.com/patrickhuber/wrangle/internal/services"
)

func TestDownload(t *testing.T) {

	h := host.NewTest(platform.Windows, nil, nil)
	defer h.Close()

	p, err := di.Invoke(h.Container(), actions.NewDownloadProvider)
	require.NoError(t, err)

	provider, ok := p.(actions.Provider)
	require.True(t, ok)

	server, err := di.Resolve[*httptest.Server](h.Container())
	require.NoError(t, err)

	path, err := di.Resolve[*filepath.Processor](h.Container())
	require.NoError(t, err)

	fs, err := di.Resolve[fs.FS](h.Container())
	require.NoError(t, err)

	task := &actions.Action{
		Type: "download",
		Parameters: map[string]any{
			"url": server.URL + "/test",
			"out": "test-local",
		},
	}
	configuration, err := di.Resolve[services.Configuration](h.Container())
	require.NoError(t, err)

	cfg, err := configuration.Get()
	require.NoError(t, err)

	metadata := actions.NewMetadataProvider(path).Get(&cfg, "test", "1.0.0")
	err = provider.Execute(task, metadata)
	require.NoError(t, err)

	// verify the file was downloaded
	ok, err = fs.Exists(path.Join(metadata.PackageVersionPath, "test-local"))
	require.NoError(t, err)
	require.True(t, ok)
}
