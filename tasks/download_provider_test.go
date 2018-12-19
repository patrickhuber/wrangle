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
	var(
		fileSystem afero.Fs
		console ui.Console
		provider tasks.Provider
	)
	BeforeEach(func(){
		fileSystem = afero.NewMemMapFs()
		console = ui.NewMemoryConsole()
		provider = tasks.NewDownloadProvider(fileSystem, console)
	})
	Describe("Execute", func(){
		It("downloads file", func() {
			server := fakes.NewHTTPServerWithArchive([]fakes.TestFile{{Path: "/data", Data: "this is data"}})
			defer server.Close()
	
			task := tasks.NewDownloadTask(
				server.URL,
				"/some/path")
			Expect(task).ToNot(BeNil())
	
			err := provider.Execute(task)
			Expect(err).To(BeNil())
	
			ok, err := afero.Exists(fileSystem, "/some/path")
			Expect(err).To(BeNil())
			Expect(ok).To(BeTrue())
		})
	})
	Describe("Unmarshal", func(){
		It("should parse task", func(){
			task, err := provider.Unmarshal("download:\n  url: https://www.google.com\n  out: /some/file\n")
			Expect(err).To(BeNil())
			Expect(task).ToNot(BeNil())
			downloadTask, ok := task.(*tasks.DownloadTask)
			Expect(ok).To(BeTrue())
			Expect(downloadTask.Details.Out).To(Equal("/some/file"))
			Expect(downloadTask.Details.URL).To(Equal("https://www.google.com"))
		})
	})		
})
