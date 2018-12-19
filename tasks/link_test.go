package tasks_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/patrickhuber/wrangle/tasks"
	yaml "gopkg.in/yaml.v2"
)

var _ = Describe("Link", func() {
	It("should serialize link", func() {
		linkTask := tasks.NewLinkTask("/source", "/destination")
		data, err := yaml.Marshal(linkTask)
		Expect(err).To(BeNil())
		Expect(string(data)).To(Equal("link:\n  source: /source\n  destination: /destination\n"))
	})
})
