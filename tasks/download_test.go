package tasks_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/patrickhuber/wrangle/tasks"
	yaml "gopkg.in/yaml.v2"
)

var _ = Describe("Download", func() {
	It("should serialize download", func() {
		downloadTask := tasks.NewDownloadTask("https://www.google.com", "/some/file")
		data, err := yaml.Marshal(downloadTask)
		Expect(err).To(BeNil())
		Expect(string(data)).To(Equal("download:\n  url: https://www.google.com\n  out: /some/file\n"))
	})
	It("should deserialize download", func() {
		data := "download:\n  url: https://www.google.com\n  out: /some/file\n"
		m := make(map[interface{}]interface{})
		err := yaml.Unmarshal([]byte(data), m)
		Expect(err).To(BeNil())
		value, ok := m["download"]
		Expect(ok).To(BeTrue())
		Expect(value).ToNot(BeNil())
	})
})
