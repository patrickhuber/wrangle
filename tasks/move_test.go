package tasks_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/patrickhuber/wrangle/tasks"
	yaml "gopkg.in/yaml.v2"
)

var _ = Describe("Move", func() {
	It("should serialize move", func() {
		moveTask := tasks.NewMoveTask("/source", "/destination")
		data, err := yaml.Marshal(moveTask)
		Expect(err).To(BeNil())
		Expect(string(data)).To(Equal("move:\n  source: /source\n  destination: /destination\n"))
	})
})
