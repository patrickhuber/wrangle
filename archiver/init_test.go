package archiver_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestArchiver(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "archiver")
}
