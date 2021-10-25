package crosspath_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestCrosspath(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Crosspath Suite")
}
