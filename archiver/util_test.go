package archiver

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Util", func() {
	Describe("commonDirectory", func() {
		Context("WhenSingleFile", func() {
			It("returns directory", func() {
				files := []string{"/test"}
				directory := commonDirectory(files...)
				Expect(directory).To(Equal("/"))
			})
		})
		Context("WhenTwoFilesRootDiretoryFile", func() {
			It("returns directory", func() {
				files := []string{"/test", "/test1"}
				directory := commonDirectory(files...)
				Expect(directory).To(Equal("/"))
			})
		})
		Context("WhenSeveralFilesNestedDirectories", func() {
			It("returns common root", func() {
				files := []string{"/a/test", "/a/b/test1", "/a/b/c/test2"}
				directory := commonDirectory(files...)
				Expect(directory).To(Equal("/a"))
			})
		})
		Context("WhenWindowsDirectories", func() {
			It("returns common root", func() {
				files := []string{"c:\\a\\test", "c:\\a\\b\\test", "c:\\a\\b"}
				directory := commonDirectory(files...)
				Expect(directory).To(Equal("c:/a"))
			})
		})
		Context("WhenFileSharesNameWithDirectory", func() {
			It("returns file directory", func() {
				files := []string{"/a/test", "/a"}
				directory := commonDirectory(files...)
				Expect(directory).To(Equal("/"))
			})
		})
	})
})
