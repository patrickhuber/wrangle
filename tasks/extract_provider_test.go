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
	var(
		fileSystem afero.Fs
		console ui.Console
		provider tasks.Provider
	)
	BeforeEach(func(){
		fileSystem = afero.NewMemMapFs()
		console = ui.NewMemoryConsole()
		provider = tasks.NewExtractProvider(fileSystem, console)
	})
	Describe("Execute", func() {
		It("should extract single file", func() {

			task := tasks.NewExtractTask("/test/test.tgz", "/destination")
			Expect(task).ToNot(BeNil())

			tgz := archiver.NewTargz(fileSystem)
			err := afero.WriteFile(fileSystem, "/test/test1", []byte("this is a test"), 0600)
			Expect(err).To(BeNil())

			err = tgz.Archive("/test/test.tgz", []string{"/test/test1"})
			Expect(err).To(BeNil())

			err = provider.Execute(task)
			Expect(err).To(BeNil())
		})
	})
	Describe("Unmarshal", func(){
		It("should parse task", func(){
			task, err := provider.Unmarshal("extract:\n  archive: /archive\n  destination: /destination\n")
			Expect(err).To(BeNil())
			Expect(task).ToNot(BeNil())
			extractTask, ok := task.(*tasks.ExtractTask)
			Expect(ok).To(BeTrue())
			Expect(extractTask.Details.Archive).To(Equal("/archive"))
			Expect(extractTask.Details.Destination).To(Equal("/destination"))
		})
	})
})
