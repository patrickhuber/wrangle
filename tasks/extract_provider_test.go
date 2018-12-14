package tasks_test

import (
	"github.com/patrickhuber/wrangle/archiver"
	"github.com/patrickhuber/wrangle/tasks"
	"github.com/patrickhuber/wrangle/ui"
	"github.com/spf13/afero"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("ExtractProvider", func() {
	Describe("Execute", func() {
		It("should extract single file", func() {

			task := tasks.NewExtractTask("/test/test.tgz", "/destination")
			Expect(task).ToNot(BeNil())

			fileSystem := afero.NewMemMapFs()
			console := ui.NewMemoryConsole()
			provider := tasks.NewExtractProvider(fileSystem, console)

			tgz := archiver.NewTargz(fileSystem)
			err := afero.WriteFile(fileSystem, "/test/test1", []byte("this is a test"), 0600)
			Expect(err).To(BeNil())

			err = tgz.Archive("/test/test.tgz", []string{"/test/test1"})
			Expect(err).To(BeNil())

			err = provider.Execute(task)
			Expect(err).To(BeNil())
		})
	})
})
