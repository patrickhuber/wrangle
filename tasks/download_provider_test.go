package tasks_test

import (
	"github.com/patrickhuber/wrangle/fakes"
	"github.com/patrickhuber/wrangle/tasks"
	"github.com/patrickhuber/wrangle/ui"
	"github.com/spf13/afero"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("DownloadProvider", func() {
	It("downloads file", func() {
		server := fakes.NewHTTPServerWithArchive([]fakes.TestFile{{Path: "/data", Data: "this is data"}})
		defer server.Close()

		task := tasks.NewDownloadTask(
			server.URL,
			"/some/path")
		Expect(task).ToNot(BeNil())

		fileSystem := afero.NewMemMapFs()
		console := ui.NewMemoryConsole()
		provider := tasks.NewDownloadProvider(fileSystem, console)

		err := provider.Execute(task)
		Expect(err).To(BeNil())

		ok, err := afero.Exists(fileSystem, "/some/path")
		Expect(err).To(BeNil())
		Expect(ok).To(BeTrue())
	})
})
