package tasks_test

import (
	. "github.com/patrickhuber/wrangle/tasks"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("ExtractTaskRunner", func() {
	It("", func() {
		task := NewTask("", "", map[string]string{"archive": "", "": ""})
		Expect(task).ToNot(BeNil())
	})
})
