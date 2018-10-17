package processes_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestProcesses(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Processes Suite")
}
