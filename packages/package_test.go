package packages

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/patrickhuber/wrangle/tasks"
)

var _ = Describe("Package", func() {
	Describe("New", func() {
		It("can replace version in download", func() {
			p := New("a", "1.2",
				tasks.NewDownloadTask("download", "https://((version))", "a_((version)).exe"))
			Expect(len(p.Tasks())).To(Equal(1))

			downloadTask := p.Tasks()[0]

			url, ok := downloadTask.Params().Lookup("url")
			Expect(ok).To(BeTrue())
			Expect(url).To(Equal("https://1.2"))

			out, ok := downloadTask.Params().Lookup("out")
			Expect(ok).To(BeTrue())
			Expect(out).To(Equal("a_1.2.exe"))
		})

		It("can replace version in extract", func() {
			p := New("a", "1.2",
				tasks.NewExtractTask("", "/test/((version))", "ab_((version))"))
			Expect(len(p.Tasks())).To(Equal(1))

			extractTask := p.Tasks()[0]

			archive, ok := extractTask.Params().Lookup("archive")
			Expect(ok).To(BeTrue())
			Expect(archive).To(Equal("/test/1.2"))
		})

	})
})
