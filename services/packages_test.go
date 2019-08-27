package services_test

import (
	"github.com/patrickhuber/wrangle/filesystem"
	"github.com/patrickhuber/wrangle/feed"
	"github.com/patrickhuber/wrangle/services"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"bytes"
	"github.com/patrickhuber/wrangle/ui"	
)

var _ = Describe("Packages", func() {
	Describe("Execute", func() {
		It("lists packages from package path", func() {
			console := ui.NewMemoryConsole()
			packagePath := "/opt/wrangle/packages"

			fileSystem := filesystem.NewMemory()
			fileSystem.Write("/opt/wrangle/packages/test/0.1.1/test.0.1.1.yml", []byte("this is a package"), 0600)

			feedService := feed.NewFsFeedService(fileSystem, packagePath)

			service := services.NewPackagesService(feedService, console)
			Expect(service).ToNot(BeNil())			
			Expect(service.List()).To(BeNil())

			output := console.OutAsString()

			var lines = make([]bytes.Buffer,3,3)
			linecount := 0

			for i:=0;i<len(output);i++{
				if output[i] == '\n'{
					linecount ++
				}else if output[i] == '\r'{

				}else{
					lines[linecount].WriteByte(output[i])
				}
			}
			Expect(lines[0].String()).To(Equal("name version"))
			Expect(lines[1].String()).To(Equal("---- -------"))
			Expect(lines[2].String()).To(Equal("test 0.1.1"))
		})
	})
})
