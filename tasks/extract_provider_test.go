package tasks_test

import (
	"github.com/patrickhuber/wrangle/filesystem"
	"github.com/patrickhuber/wrangle/archiver"
	"github.com/patrickhuber/wrangle/filepath"
	"github.com/patrickhuber/wrangle/tasks"
	"github.com/patrickhuber/wrangle/ui"	
	yaml "gopkg.in/yaml.v2"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("ExtractProvider", func() {
	var (
		fileSystem filesystem.FileSystem
		console    ui.Console
		provider   tasks.Provider
	)
	BeforeEach(func() {
		fileSystem = filesystem.NewMemory()
		console = ui.NewMemoryConsole()
		provider = tasks.NewExtractProvider(fileSystem, console)
	})
	Describe("Execute", func() {
		It("should extract single file", func() {
			taskContext := newFakeTaskContext("/opt/wrangle", "test", "1.0.0")

			filePath := filepath.Join(taskContext.PackageVersionPath(), "test1")
			err := fileSystem.Write(filePath, []byte("this is a test"), 0600)
			Expect(err).To(BeNil())

			archivePath := filepath.Join(taskContext.PackageVersionPath(), "test.tgz")
			tgz := archiver.NewTargz(fileSystem)
			err = tgz.Archive(archivePath, []string{filePath})
			Expect(err).To(BeNil())

			err = fileSystem.Remove(filePath)
			Expect(err).To(BeNil())

			task := tasks.NewExtractTask("test.tgz")
			Expect(task).ToNot(BeNil())

			err = provider.Execute(task, taskContext)
			Expect(err).To(BeNil())

			ok, err := fileSystem.Exists(archivePath)
			Expect(err).To(BeNil())
			Expect(ok).To(BeTrue())

			ok, err = fileSystem.Exists(filePath)
			Expect(err).To(BeNil())
			Expect(ok).To(BeTrue())
		})
	})
	Describe("Decode", func() {
		It("should parse task", func() {

			m := make(map[string]interface{})
			err := yaml.Unmarshal([]byte("extract:\n  archive: /archive\n"), m)
			Expect(err).To(BeNil())

			task, err := provider.Decode(m)
			Expect(err).To(BeNil())
			Expect(task).ToNot(BeNil())

			extractTask, ok := task.(*tasks.ExtractTask)
			Expect(ok).To(BeTrue())
			Expect(extractTask.Details.Archive).To(Equal("/archive"))
		})
	})
})
