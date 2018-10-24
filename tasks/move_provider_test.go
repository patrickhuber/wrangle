package tasks_test

import (
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
			task := NewTask("", "", map[string]string{"source": "/test/file", "destination": "/test/renamed"})

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
			task := NewTask("", "", map[string]string{"source": "/test", "destination": "/test1"})

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
})
