package tasks_test

import (
	"github.com/patrickhuber/wrangle/filepath"
	"github.com/patrickhuber/wrangle/filesystem"
	"github.com/patrickhuber/wrangle/tasks"
	"github.com/patrickhuber/wrangle/ui"
	yaml "gopkg.in/yaml.v2"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("MoveProvider", func() {
	var (
		provider   tasks.Provider
		fileSystem filesystem.FileSystem
		console    ui.Console
	)
	BeforeEach(func() {
		fileSystem = filesystem.NewMemory()
		console = ui.NewMemoryConsole()
		provider = tasks.NewMoveProvider(fileSystem, console)

	})
	Describe("Execute", func() {
		var (
			taskContext tasks.TaskContext
		)
		BeforeEach(func() {
			taskContext = newFakeTaskContext("/opt/wrangle", "test", "1.0.0")
		})
		It("can move file", func() {
			sourcePath := filepath.Join(taskContext.PackageVersionPath(), "file")
			fileSystem.Write(sourcePath, []byte("test"), 0666)
			task := tasks.NewMoveTask("file", "renamed")

			err := provider.Execute(task, taskContext)
			Expect(err).To(BeNil())

			destinationPath := filepath.Join(taskContext.PackageVersionPath(), "renamed")
			exists, err := fileSystem.Exists(destinationPath)
			Expect(err).To(BeNil())
			Expect(exists).To(BeTrue())

			isDirectory, err := fileSystem.IsDir(destinationPath)
			Expect(err).To(BeNil())
			Expect(isDirectory).To(BeFalse())
		})
		It("can move directory", func() {

			sourcePath := filepath.Join(taskContext.PackageVersionPath(), "folder/sub/file")
			fileSystem.Write(sourcePath, []byte("test"), 0666)

			task := tasks.NewMoveTask("folder/sub", "folder")
			err := provider.Execute(task, taskContext)
			Expect(err).To(BeNil())

			/* destinationPath := filepath.Join(taskContext.PackageVersionPath(), "folder/file")
			exists, err := fileSystem.Exists(destinationPath)
			Expect(err).To(BeNil())
			Expect(exists).To(BeTrue()) */

			destinationDirectory := filepath.Join(taskContext.PackageVersionPath(), "folder")
			isDirectory, err := fileSystem.IsDir(destinationDirectory)
			Expect(err).To(BeNil())
			Expect(isDirectory).To(BeTrue())

			/*
				exists, err = fileSystem.Exists( "/test1/file")
				Expect(err).To(BeNil())
				Expect(exists).To(BeTrue())
			*/
		})
	})
	Describe("Decode", func() {
		It("should parse task", func() {
			m := make(map[string]interface{})
			err := yaml.Unmarshal([]byte("move:\n  source: /source\n  destination: /destination\n"), m)
			Expect(err).To(BeNil())
			task, err := provider.Decode(m)
			Expect(err).To(BeNil())
			Expect(task).ToNot(BeNil())
			moveTask, ok := task.(*tasks.MoveTask)
			Expect(ok).To(BeTrue())
			Expect(moveTask.Details.Source).To(Equal("/source"))
			Expect(moveTask.Details.Destination).To(Equal("/destination"))
		})
	})
})
