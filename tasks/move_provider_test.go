package tasks_test

import (
	"gopkg.in/yaml.v2"
	"github.com/patrickhuber/wrangle/tasks"
	. "github.com/patrickhuber/wrangle/tasks"
	"github.com/patrickhuber/wrangle/ui"
	"github.com/spf13/afero"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("MoveProvider", func() {
	var (
		provider   Provider
		fileSystem afero.Fs
		console    ui.Console
	)
	BeforeEach(func() {
		fileSystem = afero.NewMemMapFs()
		console = ui.NewMemoryConsole()
		provider = NewMoveProvider(fileSystem, console)
	})
	Describe("Execute", func() {
		It("can move file", func() {
			afero.WriteFile(fileSystem, "/test/file", []byte("test"), 0666)
			task := NewMoveTask("/test/file", "/test/renamed")

			err := provider.Execute(task)
			Expect(err).To(BeNil())

			exists, err := afero.Exists(fileSystem, "/test/renamed")
			Expect(err).To(BeNil())
			Expect(exists).To(BeTrue())

			isDirectory, err := afero.IsDir(fileSystem, "/test/renamed")
			Expect(err).To(BeNil())
			Expect(isDirectory).To(BeFalse())
		})
		It("can move directory", func() {
			afero.WriteFile(fileSystem, "/test/file", []byte("test"), 0666)
			task := NewMoveTask("/test", "/test1")

			err := provider.Execute(task)
			Expect(err).To(BeNil())

			exists, err := afero.Exists(fileSystem, "/test1")
			Expect(err).To(BeNil())
			Expect(exists).To(BeTrue())

			isDirectory, err := afero.IsDir(fileSystem, "/test1")
			Expect(err).To(BeNil())
			Expect(isDirectory).To(BeTrue())

			/*
				exists, err = afero.Exists(fileSystem, "/test1/file")
				Expect(err).To(BeNil())
				Expect(exists).To(BeTrue())
			*/
		})
	})
	Describe("Decode", func() {
		It("should parse task", func() {
			m:= make(map[string]interface{})
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
