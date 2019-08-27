package archiver_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/patrickhuber/wrangle/archiver"
	"github.com/patrickhuber/wrangle/filesystem"
)

var _ = Describe("Tgz", func() {
	Describe("RoundTrip", func() {
		It("can write and read back a tgz file", func() {
			fileSystem := filesystem.NewMemory()

			err := fileSystem.Write("/tmp/test", []byte("this is a test"), 0666)
			Expect(err).To(BeNil())

			a := archiver.NewTargz(fileSystem)
			err = a.Archive("/tmp/temp.tgz", []string{"/tmp/test"})
			Expect(err).To(BeNil())

			err = fileSystem.Remove("/tmp/test")
			Expect(err).To(BeNil())

			err = a.Extract("/tmp/temp.tgz", "/tmp", []string{".*"})
			Expect(err).To(BeNil())

			ok, err := fileSystem.Exists("/tmp/test")
			Expect(err).To(BeNil())
			Expect(ok).To(BeTrue(), "file /tmp/test not found")
		})
	})
})
