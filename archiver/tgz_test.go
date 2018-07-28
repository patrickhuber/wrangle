package archiver_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/patrickhuber/wrangle/archiver"
	"github.com/patrickhuber/wrangle/filesystem"
	"github.com/spf13/afero"
)

var _ = Describe("Tgz", func() {
	Describe("RoundTrip", func() {
		It("can write and read back a tgz file", func() {
			fileSystem := filesystem.NewMemMapFsWrapper(afero.NewMemMapFs())

			err := afero.WriteFile(fileSystem, "/tmp/test", []byte("this is a test"), 0666)
			Expect(err).To(BeNil())

			output, err := fileSystem.Create("/tmp/temp.tgz")
			Expect(err).To(BeNil())
			defer output.Close()

			a := archiver.NewTargzArchiver(fileSystem)
			err = a.Archive(output, []string{"/tmp/test"})
			Expect(err).To(BeNil())

			source, err := fileSystem.Open("/tmp/temp.tgz")
			Expect(err).To(BeNil())
			defer source.Close()

			err = a.Extract(source, ".*", "/tmp")
			Expect(err).To(BeNil())
		})
	})
})
