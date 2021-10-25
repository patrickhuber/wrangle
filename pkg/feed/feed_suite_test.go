package feed_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestFeed(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Feed Suite")
}
