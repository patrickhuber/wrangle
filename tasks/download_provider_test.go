package tasks_test

import (
	"github.com/patrickhuber/wrangle/filepath"
	"gopkg.in/yaml.v2"
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
				"file")
			Expect(task).ToNot(BeNil())
	
			taskContext:= newTaskContext("/opt/wrangle", "test", "1.0.0")
			err := provider.Execute(task, taskContext)
			Expect(err).To(BeNil())
	
			expected := filepath.Join(taskContext.PackageVersionPath(), "file")
			ok, err := afero.Exists(fileSystem, expected)
			Expect(err).To(BeNil())
			Expect(ok).To(BeTrue())
		})
	})
	Describe("Decode", func(){
		It("should parse task", func(){

			m:= make(map[string]interface{})
			err := yaml.Unmarshal([]byte("download:\n  url: https://www.google.com\n  out: /some/file\n"), m)
			Expect(err).To(BeNil())

			task, err := provider.Decode(m)			
			Expect(err).To(BeNil())
			Expect(task).ToNot(BeNil())
			
			downloadTask, ok := task.(*tasks.DownloadTask)
			Expect(ok).To(BeTrue())
			Expect(downloadTask.Details.Out).To(Equal("/some/file"))
			Expect(downloadTask.Details.URL).To(Equal("https://www.google.com"))
		})
	})		
})
