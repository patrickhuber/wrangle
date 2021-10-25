package crosspath_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/patrickhuber/wrangle/pkg/crosspath"
)

var _ = Describe("Util", func() {
	It("preserves prefix slash", func() {
		segments := []string{"/test", "/sub"}
		result := crosspath.Join(segments...)
		Expect(result).To(Equal("/test/sub"))
	})
	It("converts backslash to forward slash", func() {
		segments := []string{"\\test", "\\sub"}
		result := crosspath.Join(segments...)
		Expect(result).To(Equal("/test/sub"))
	})
	It("adds slashes between elements", func() {
		segments := []string{"/test", "sub"}
		result := crosspath.Join(segments...)
		Expect(result).To(Equal("/test/sub"))
	})
	It("removes right slash from segments", func() {
		segments := []string{"/test/", "sub/"}
		result := crosspath.Join(segments...)
		Expect(result).To(Equal("/test/sub"))
	})
	It("removes duplicate slashes", func() {
		segments := []string{"/test/", "/sub/"}
		result := crosspath.Join(segments...)
		Expect(result).To(Equal("/test/sub"))
	})
})
