package packages_test

import (
	"os"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/patrickhuber/wrangle/pkg/packages"
)

var _ = Describe("Package", func() {
	Describe("FromYaml", func() {
		It("unmarshalls simple manifest", func() {

			manifestPath := "./fakes/simple.yml"
			content, err := os.ReadFile(manifestPath)
			Expect(err).To(BeNil())
			p, err := packages.FromYaml(string(content))
			Expect(err).To(BeNil())
			Expect(p.Name).To(Equal("bosh"))
			Expect(p.Versions).ToNot(BeNil())
			Expect(len(p.Versions)).To(Equal(1))
			for _, v := range p.Versions {
				Expect(len(v.Targets)).To(Equal(1))
			}
		})
		It("unmarshalls complex manifest", func() {

			manifestPath := "./fakes/complex.yml"
			content, err := os.ReadFile(manifestPath)
			Expect(err).To(BeNil())
			p, err := packages.FromYaml(string(content))
			Expect(err).To(BeNil())
			Expect(p.Name).To(Equal("bbr"))
			Expect(p.Versions).ToNot(BeNil())
			Expect(len(p.Versions)).To(Equal(1))
			for _, v := range p.Versions {
				Expect(len(v.Targets)).To(Equal(2))
			}
		})
	})
})
