package githubrelease_test

import (
	"regexp"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/patrickhuber/wrangle/pkg/resource/githubrelease"
)

var _ = Describe("Source", func() {
	It("regex find must parse and match version", func() {
		source := githubrelease.Source{
			TagFilter: "[0-9]([.][0-9]){2}",
		}
		re := regexp.MustCompile(source.TagFilter)
		findString := re.FindString("v1.2.3")
		Expect(findString).To(Equal("1.2.3"))
		matchString := re.MatchString("v1.2.3")
		Expect(matchString).To(BeTrue())
	})
})
