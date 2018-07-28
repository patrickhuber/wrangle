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

			// create the zip file from the test file
			outputFile, err := fs.Create("/test.zip")
			Expect(err).To(BeNil())
			defer outputFile.Close()

			// create the archiver and write out the archive
			arch := archiver.NewZipArchiver(fs)
			err = arch.Archive(outputFile, []string{"/test"})
			Expect(err).To(BeNil())

			// verify the file exists
			ok, err := afero.Exists(fs, "/test.zip")
			Expect(err).To(BeNil())
			Expect(ok).To(BeTrue())

			// remove the old text file
			err = fs.Remove("/test")
			Expect(err).To(BeNil())

			// open the archive
			inputFile, err := fs.Open("/test.zip")
			Expect(err).To(BeNil())
			defer inputFile.Close()

			// extract the archive
			err = arch.Extract(inputFile, ".*", "/")
			Expect(err).To(BeNil())

			// verify the file exists
			ok, err = afero.Exists(fs, "/test")
			Expect(err).To(BeNil())
			Expect(ok).To(BeTrue())
		})
	})
})
