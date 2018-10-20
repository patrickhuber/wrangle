package tasks_test

import (
	. "github.com/patrickhuber/wrangle/tasks"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Task", func() {
	var (
		task Task
	)
	BeforeEach(func() {
		task = NewTask("name", "type", map[string]string{"key": "value"})
	})
	Describe("Name", func() {
		It("should return name", func() {
			Expect(task.Name()).To(Equal("name"))
		})
	})
	Describe("Type", func() {
		It("should return type", func() {
			Expect(task.Type()).To(Equal("type"))
		})
	})
	Describe("Params", func() {
		It("should return params", func() {
			params := task.Params()
			Expect(params).ToNot(BeNil())
			value, ok := params.Lookup("key")
			Expect(ok).To(BeTrue())
			Expect(value).To(Equal("value"))
		})
	})
})
