package items_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestItems(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Items Suite")
}
