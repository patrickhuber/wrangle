package tasks_test

import (
	. "github.com/patrickhuber/wrangle/tasks"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("DownloadTaskRunner", func() {
	It("", func() {
		task := NewTask("", "", map[string]string{"url": "https://localhost", "out": "/some/path"})
		Expect(task).ToNot(BeNil())
	})
})
