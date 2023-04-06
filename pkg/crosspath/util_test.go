package crosspath_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/patrickhuber/wrangle/pkg/crosspath"
)

var _ = DescribeTable("Util", func(segments []string) {

	result := crosspath.Join(segments...)
	Expect(result).To(Equal("/test/sub"))
},
	Entry("preserves prefix slash", []string{"/test", "/sub"}),
	Entry("converts backslash to forward slash", []string{"\\test", "\\sub"}),
	Entry("adds slashes between elements", []string{"/test", "sub"}),
	Entry("removes right slash from segments", []string{"/test/", "sub/"}),
	Entry("removes duplicate slashes", []string{"/test/", "/sub/"}),
)
