package ilog_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestIlog(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Ilog Suite")
}
