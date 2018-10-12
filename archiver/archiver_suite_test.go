package archiver_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestArchiver(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Archiver Suite")
}
