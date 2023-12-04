package actions_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/patrickhuber/go-log"
	"github.com/patrickhuber/go-xplat/platform"
	"github.com/patrickhuber/go-xplat/setup"
	"github.com/patrickhuber/wrangle/internal/actions"
)

func TestMove(t *testing.T) {
	h := setup.NewTest(setup.Platform(platform.Linux))
	path := h.Path
	fs := h.FS
	logger := log.Memory()
	provider := actions.NewMoveProvider(fs, path, logger)

	err := fs.WriteFile("/folder/file.txt", []byte("this is a test"), 0644)
	require.Nil(t, err)

	action := &actions.Action{
		Type: "move",
		Parameters: map[string]any{
			"source":      "file.txt",
			"destination": "moved.txt",
		},
	}
	ctx := &actions.Metadata{
		PackageVersionPath: "/folder",
	}
	err = provider.Execute(action, ctx)
	require.Nil(t, err)
	ok, err := fs.Exists("/folder/moved.txt")
	require.Nil(t, err)
	require.True(t, ok)
}
