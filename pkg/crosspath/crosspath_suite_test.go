package crosspath_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestCrosspath(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Crosspath Suite")
}
