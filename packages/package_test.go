package packages

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/patrickhuber/wrangle/tasks"
)

var _ = Describe("Package", func() {
	Describe("New", func() {
		It("can create package", func() {
			p := New("a", "1.2",
				nil,
				tasks.NewDownloadTask("https://1.2", "a_1.2.exe"))
			Expect(p).ToNot(BeNil())
			Expect(p.Name()).To(Equal("a"))
			Expect(p.Version()).To(Equal("1.2"))
			Expect(len(p.Tasks())).To(Equal(1))
		})
	})
})
