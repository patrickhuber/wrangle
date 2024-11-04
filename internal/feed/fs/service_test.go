package fs_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/patrickhuber/go-cross"
	"github.com/patrickhuber/go-cross/arch"
	"github.com/patrickhuber/go-cross/platform"
	"github.com/patrickhuber/go-log"
	"github.com/patrickhuber/wrangle/internal/feed"
	"github.com/patrickhuber/wrangle/internal/feed/conformance"
	feedfs "github.com/patrickhuber/wrangle/internal/feed/fs"
)

func TestService(t *testing.T) {
	t.Run("can list all packages", func(t *testing.T) {
		tester := SetupServiceTest(t)
		tester.CanListAllPackages(t)
	})
	t.Run("can return latest version", func(t *testing.T) {
		tester := SetupServiceTest(t)
		tester.CanReturnLatestVersion(t)
	})
	t.Run("can return specific version", func(t *testing.T) {
		tester := SetupServiceTest(t)
		tester.CanReturnSpecificVersion(t)
	})
	t.Run("can add version", func(t *testing.T) {
		tester := SetupServiceTest(t)
		tester.CanAddVersion(t)
	})
	t.Run("can update existing version", func(t *testing.T) {
		tester := SetupServiceTest(t)
		tester.CanUpdateExistingVersion(t)
	})
}

func SetupServiceTest(t *testing.T) conformance.ServiceTester {
	h := cross.NewTest(platform.Linux, arch.AMD64)
	fs := h.FS()
	path := h.Path()
	logger := log.Memory()
	svc := feedfs.NewService("test", fs, path, "/opt/wrangle/feed", logger)
	items := conformance.GetItemList()

	response, err := svc.Update(&feed.UpdateRequest{
		Items: &feed.ItemUpdate{
			Add: items,
		},
	})
	require.NoError(t, err)
	require.Equal(t, conformance.TotalItemCount, response.Changed)
	return conformance.NewServiceTester(svc)
}
