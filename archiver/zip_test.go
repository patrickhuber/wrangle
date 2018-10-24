package archiver_test

import (
	"github.com/patrickhuber/wrangle/archiver"
	"github.com/spf13/afero"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Zip", func() {
	Describe("RoundTrip", func() {
		It("can write and read back a zip file", func() {
			fs := afero.NewMemMapFs()

			// create the test file
			err := afero.WriteFile(fs, "/test", []byte("test"), 0666)
			Expect(err).To(BeNil())

			// create the archiver and write out the archive
			arch := archiver.NewZip(fs)
			err = arch.Archive("/test.zip", []string{"/test"})
			Expect(err).To(BeNil())

			// verify the file exists
			ok, err := afero.Exists(fs, "/test.zip")
			Expect(err).To(BeNil())
			Expect(ok).To(BeTrue())

			// remove the old text file
			err = fs.Remove("/test")
			Expect(err).To(BeNil())

			// extract the archive
			err = arch.Extract("/test.zip", "/", []string{".*"})
			Expect(err).To(BeNil())

			// verify the file exists
			ok, err = afero.Exists(fs, "/test")
			Expect(err).To(BeNil())
			Expect(ok).To(BeTrue())
		})
	})
})
