package actions_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/patrickhuber/go-log"
	"github.com/patrickhuber/go-xplat/filepath"
	"github.com/patrickhuber/go-xplat/fs"
	"github.com/stretchr/testify/require"

	"github.com/patrickhuber/wrangle/pkg/actions"
	"github.com/patrickhuber/wrangle/pkg/config"
)

func TestDownload(t *testing.T) {
	// start the local http server
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		if strings.HasSuffix(req.URL.Path, "/test-remote") {
			rw.Write([]byte("hello"))
			return
		}
		rw.WriteHeader(http.StatusNotFound)
		rw.Write([]byte("not found"))
	}))

	defer server.Close()

	cfg := &config.Config{
		Paths: &config.Paths{
			Packages: "/wrangle/packages",
			Bin:      "/wrangle/bin",
			Root:     "/wrangle",
		},
	}
	fs := fs.NewMemory()
	path := filepath.NewProcessor()
	provider := actions.NewDownloadProvider(log.Default(), fs, path)
	task := &actions.Action{
		Type: "download",
		Parameters: map[string]any{
			"url": server.URL + "/test-remote",
			"out": "test-local",
		},
	}
	metadata := actions.NewMetadataProvider(path).Get(cfg, "test", "1.0.0")
	err := provider.Execute(task, metadata)
	require.NoError(t, err)

	// verify the file was downloaded
	ok, err := fs.Exists(path.Join(metadata.PackageVersionPath, "test-local"))
	require.NoError(t, err)
	require.True(t, ok)
}
