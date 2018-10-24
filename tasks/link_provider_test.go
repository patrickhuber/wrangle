package tasks_test

import (
	"github.com/patrickhuber/wrangle/filesystem"
	"github.com/patrickhuber/wrangle/tasks"
	"github.com/patrickhuber/wrangle/ui"
	"github.com/spf13/afero"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("LinkProvider", func() {
	It("can create symlink", func() {
		fs := filesystem.NewMemMapFs()
		afero.WriteFile(fs, "/source", []byte("this is data"), 0600)

		console := ui.NewMemoryConsole()
		provider := tasks.NewLinkProvider(fs, console)

		task := tasks.NewTask("", "", map[string]string{
			"source":      "/source",
			"destination": "/destination",
		})

		err := provider.Execute(task)
		Expect(err).To(BeNil())

		ok, err := afero.Exists(fs, "/destination")
		Expect(err).To(BeNil())
		Expect(ok).To(BeTrue())
	})
})
