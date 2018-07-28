package commands

import (
	"fmt"
	"net/http/httptest"
	"strings"

	"github.com/patrickhuber/wrangle/filepath"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/spf13/afero"

	"github.com/patrickhuber/wrangle/config"
	"github.com/patrickhuber/wrangle/fakes"
	"github.com/patrickhuber/wrangle/filesystem"
	"github.com/patrickhuber/wrangle/ui"
)

var _ = Describe("Install", func() {
	Describe("NewInstall", func() {
		It("returns value", func() {
			platform := "windows"
			outFolder := ""
			fileSystem := filesystem.NewOsFsWrapper(afero.NewMemMapFs())
			_, err := NewInstall(platform, outFolder, fileSystem, ui.NewMemoryConsole())
			Expect(err).ToNot(BeNil())
		})

	})
	Describe("Execute", func() {
		var (
			platform   string
			outFolder  string
			fileName   string
			extractOut string
			server     *httptest.Server
		)
		BeforeSuite(func() {
			server = fakes.NewHTTPServerWithArchive(
				[]fakes.TestFile{
					{Path: "/test", Data: "this is data"},
					{Path: "/test.exe", Data: "this is data"},
				})
		})
		AfterSuite(func() {
			server.Close()
		})
		AfterEach(func() {
			url := server.URL
			if !strings.HasSuffix(url, "/") {
				url += "/"
			}
			url += fileName
			err := runInstallCommand(platform, outFolder, url, fileName, extractOut)
			Expect(err).To(BeNil())
		})
		Context("WhenWindows", func() {
			BeforeEach(func() {
				platform = "windows"
				outFolder = "c:\\out"
				extractOut = "test.exe"
			})
			Context("WhenTar", func() {
				It("installs", func() {
					fileName = "test.tar"
				})
			})
			Context("WhenTgz", func() {
				It("installs", func() {
					fileName = "test.tgz"
				})
			})
			Context("WhenZip", func() {
				It("installs", func() {
					fileName = "test.zip"
				})
			})
			Context("WhenBinary", func() {
				It("installs", func() {
					fileName = "test.exe"
					extractOut = ""
				})
			})
		})
		Context("WhenDarwin", func() {
			BeforeEach(func() {
				platform = "darwin"
				outFolder = "/out"
				extractOut = "test"
			})
			Context("WhenTar", func() {
				It("installs", func() {
					fileName = "test.tar"
				})
			})
			Context("WhenTgz", func() {
				It("installs", func() {
					fileName = "test.tgz"
				})
			})
			Context("WhenZip", func() {
				It("installs", func() {
					fileName = "test.zip"
				})
			})
			Context("WhenBinary", func() {
				It("installs", func() {
					extractOut = ""
					fileName = "test"
				})
			})
		})
		Context("WhenLinux", func() {
			BeforeEach(func() {
				platform = "darwin"
				outFolder = "/out"
				extractOut = "test"
			})
			Context("WhenTar", func() {
				It("installs", func() {
					fileName = "test.tar"
				})
			})
			Context("WhenTgz", func() {
				It("installs", func() {
					fileName = "test.tgz"
				})
			})
			Context("WhenZip", func() {
				It("installs", func() {
					fileName = "test.zip"
				})
			})
			Context("WhenBinary", func() {
				It("installs", func() {
					fileName = "test"
					extractOut = ""
				})
			})
		})
	})
})

func runInstallCommand(platform, outFolder, downloadURL, downloadOut, extractOut string) error {
	version := "1.0.0"
	name := "test"
	alias := "alias"
	extractFilter := extractOut

	content := getContent(name, version, platform, alias, downloadURL, downloadOut, extractFilter, extractOut)
	fs := filesystem.NewMemMapFs()
	command, err := NewInstall(platform, outFolder, fs, ui.NewMemoryConsole())
	if err != nil {
		return err
	}

	cfg, err := config.SerializeString(content)
	if err != nil {
		return err
	}

	err = command.Execute(cfg, name)
	if err != nil {
		return err
	}

	// verify the downloaded file exists
	expectedFilePath := filepath.Join(outFolder, downloadOut)
	expectedFilePath = filepath.ToSlash(expectedFilePath)
	ok, err := afero.Exists(fs, expectedFilePath)
	if err != nil {
		return err
	}
	if !ok {
		return fmt.Errorf("unable to find file '%s'", expectedFilePath)
	}

	if extractFilter != "" && extractOut != "" {
		// verify the extracted file exists
		expectedFilePath = filepath.Join(outFolder, extractOut)
		expectedFilePath = filepath.ToSlash(expectedFilePath)
		ok, err := afero.Exists(fs, expectedFilePath)
		if err != nil {
			return err
		}
		if !ok {
			return fmt.Errorf("unable to find file '%s'", expectedFilePath)
		}
	}
	return nil
}

func getContent(name, version, platform, alias, downloadURL, downloadOut, extractFilter, extractOut string) string {
	content := `
packages:
- name: ((name))
  version: ((version))
  platforms:
  - name: ((platform))
    alias: ((alias))
    download:
      url: ((download_url))
      out: ((download_out))
`
	if extractFilter != "" && extractOut != "" {
		content += `
    extract:
      filter: ((extract_filter))
      out: ((extract_out))
`
		content = strings.Replace(content, "((extract_filter))", extractFilter, -1)
		content = strings.Replace(content, "((extract_out))", extractOut, -1)
	}
	content = strings.Replace(content, "((name))", name, -1)
	content = strings.Replace(content, "((version))", version, -1)
	content = strings.Replace(content, "((platform))", platform, -1)
	content = strings.Replace(content, "((alias))", alias, -1)
	content = strings.Replace(content, "((download_url))", downloadURL, -1)
	content = strings.Replace(content, "((download_out))", downloadOut, -1)
	return content
}
