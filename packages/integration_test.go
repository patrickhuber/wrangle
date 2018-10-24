//+build integration

package packages_test

import (
	"runtime"
	"strings"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/patrickhuber/wrangle/filesystem"
	"github.com/patrickhuber/wrangle/packages"
)

var _ = Describe("ManagerIntegration", func() {
	Describe("Execute", func() {
		Context("WhenInstallingCredHubCli", func() {
			It("unpacks and installs correctly", func() {

				platform := runtime.GOOS

				packageDir := ""
				extractOut := "credhub-((platform))-((version))"
				alias := "credhub"

				if platform == "windows" {
					packageDir = "c:/tools/wrangle/packages"
					extractOut = extractOut + ".exe"
					alias = alias + ".exe"
				} else {
					packageDir = "/opt/wrangle/packages"
				}
				fs := filesystem.NewOsFs()
				manager := packages.NewManager(fs)

				version := "2.0.0"
				url := "https://github.com/cloudfoundry-incubator/credhub-cli/releases/download/((version))/credhub-((platform))-((version)).tgz"
				url = strings.Replace(url, "((version))", version, -1)
				url = strings.Replace(url, "((platform))", platform, -1)

				downloadOut := "credhub-((platform))-((version)).tgz"
				downloadOut = strings.Replace(downloadOut, "((version))", version, -1)
				downloadOut = strings.Replace(downloadOut, "((platform))", platform, -1)

				extractOut = strings.Replace(extractOut, "((version))", version, -1)
				extractOut = strings.Replace(extractOut, "((platform))", platform, -1)

				p := packages.New("credhub", version, alias,
					packages.NewDownload(
						url,
						packageDir,
						downloadOut),
					packages.NewExtract(
						"credhub",
						packageDir,
						extractOut))

				err := manager.Install(p)
				Expect(err).To(BeNil())
			})
		})
	})
})
