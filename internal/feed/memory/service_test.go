package memory_test

import (
	"testing"

	"github.com/patrickhuber/go-log"
	"github.com/patrickhuber/wrangle/internal/feed/conformance"
	"github.com/patrickhuber/wrangle/internal/feed/memory"
)

func TestService(t *testing.T) {
	setup := func(t *testing.T) conformance.ServiceTester {
		items := conformance.GetItemList()
		logger := log.Memory()
		service := memory.NewService("test", logger, items...)
		return conformance.NewServiceTester(service)
	}
	t.Run("can list all packages", func(t *testing.T) {
		tester := setup(t)
		tester.CanListAllPackages(t)
	})
	t.Run("can return latest version", func(t *testing.T) {
		tester := setup(t)
		tester.CanReturnLatestVersion(t)
	})
	t.Run("can return specific version", func(t *testing.T) {
		tester := setup(t)
		tester.CanReturnSpecificVersion(t)
	})
	t.Run("can add version", func(t *testing.T) {
		tester := setup(t)
		tester.CanAddVersion(t)
	})
	t.Run("can update existing version", func(t *testing.T) {
		tester := setup(t)
		tester.CanUpdateExistingVersion(t)
	})
}
