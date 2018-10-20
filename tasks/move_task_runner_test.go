package tasks_test

import (
	. "github.com/patrickhuber/wrangle/tasks"
	"github.com/spf13/afero"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("MoveTaskRunner", func() {
	var (
		taskRunner TaskRunner
		fileSystem afero.Fs
	)
	BeforeEach(func() {
		fileSystem = afero.NewMemMapFs()
		taskRunner = NewMoveTaskRunner(fileSystem)
	})
	Describe("Execute", func() {
		It("can move file", func() {
			afero.WriteFile(fileSystem, "/test/file", []byte("test"), 0666)
			task := NewTask("", "", map[string]string{"source": "/test/file", "destination": "/test/renamed"})

			err := taskRunner.Execute(task)
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

			err := taskRunner.Execute(task)
			Expect(err).To(BeNil())

			exists, err := afero.Exists(fileSystem, "/test1")
			Expect(err).To(BeNil())
			Expect(exists).To(BeTrue())

			isDirectory, err := afero.IsDir(fileSystem, "/test1")
			Expect(err).To(BeNil())
			Expect(isDirectory).To(BeTrue())

			exists, err = afero.Exists(fileSystem, "/test1/file")
			Expect(err).To(BeNil())
			Expect(exists).To(BeTrue())
		})
	})
})
