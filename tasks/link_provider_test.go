package tasks_test

import (
	"github.com/patrickhuber/wrangle/filepath"
	"github.com/patrickhuber/wrangle/filesystem"
	"github.com/patrickhuber/wrangle/tasks"
	"github.com/patrickhuber/wrangle/ui"	

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("LinkProvider", func() {
	It("can create symlink", func() {

		tc := newTaskContext("/wrangle", "test", "1.0.0")
		sourceFile := filepath.Join(tc.PackageVersionPath(), "source")
		fs := filesystem.NewMemory()
		fs.Write(sourceFile, []byte("this is data"), 0600)

		console := ui.NewMemoryConsole()
		provider := tasks.NewLinkProvider(fs, console)

		task := tasks.NewLinkTask("source", "destination")

		err := provider.Execute(task, tc)
		Expect(err).To(BeNil())

		expected := filepath.Join(tc.Bin(), "destination")
		ok, err := fs.Exists(expected)
		Expect(err).To(BeNil())
		Expect(ok).To(BeTrue())
	})
	Describe("NewLinkTask", func() {
		It("should map parameters", func() {
			task := tasks.NewLinkTask("source", "alias")

			s, ok := task.Params()["source"]
			Expect(ok).To(BeTrue())
			Expect(s).To(Equal("source"))

			d, ok := task.Params()["alias"]
			Expect(ok).To(BeTrue())
			Expect(d).To(Equal("alias"))
		})
	})
})
