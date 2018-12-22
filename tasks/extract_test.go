package tasks_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/patrickhuber/wrangle/tasks"
	yaml "gopkg.in/yaml.v2"
)

var _ = Describe("Extract", func() {
	It("should serialize extract", func() {
		extractTask := tasks.NewExtractTask("/some/file")
		data, err := yaml.Marshal(extractTask)
		Expect(err).To(BeNil())
		Expect(string(data)).To(Equal("extract:\n  archive: /some/file\n"))
	})
})
