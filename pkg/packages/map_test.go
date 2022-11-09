package packages_test

import (
	"bytes"
	"os"

	"github.com/google/go-cmp/cmp"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/patrickhuber/wrangle/pkg/packages"
	"gopkg.in/yaml.v3"
)

var manifest = &packages.Manifest{
	Package: &packages.ManifestPackage{
		Name:    "test",
		Version: "1.0.0",
		Targets: []*packages.ManifestTarget{
			{
				Platform:     "linux",
				Architecture: "amd64",
				Steps: []packages.ManifestStep{
					{
						Action: "download",
						With: map[string]any{
							"url": "https://www.google.com",
							"out": "index.html",
						},
					},
				},
			},
		},
	},
}
var version = &packages.Version{
	Version: "1.0.0",
	Targets: []*packages.Target{
		{
			Platform:     "linux",
			Architecture: "amd64",
			Tasks: []*packages.Task{
				{
					Name: "download",
					Properties: map[string]any{
						"url": "https://www.google.com",
						"out": "index.html",
					},
				},
			},
		},
	},
}

var _ = Describe("Map", func() {
	Describe("ManifestToPackageVersion", func() {
		It("can map", func() {
			actual := packages.ManifestToPackageVersion(manifest)
			Expect(actual).To(Equal(version))
		})
	})
	Describe("PackageVersionToManifest", func() {
		It("can map", func() {
			actual := packages.PackageVersionToManifest("test", version)
			Expect(actual).To(Equal(manifest))
		})
	})
	Describe("ManifestToPackageVersionYaml", func() {
		It("can map", func() {
			m := &packages.Manifest{}
			bytes, err := os.ReadFile("./fakes/simple.yml")
			Expect(err).To(BeNil())
			err = yaml.Unmarshal(bytes, m)
			Expect(err).To(BeNil())
			actual := packages.ManifestToPackageVersion(m)
			diff := cmp.Diff(actual, version)
			Expect(diff).To(Equal(""), diff)
		})
	})
	Describe("PackageVersionToManifestYaml", func() {
		It("can map", func() {
			actual := packages.PackageVersionToManifest("test", version)
			var b bytes.Buffer
			encoder := yaml.NewEncoder(&b)
			encoder.SetIndent(2)
			err := encoder.Encode(actual)
			Expect(err).To(BeNil())
			versionYaml, err := os.ReadFile("./fakes/simple.yml")
			Expect(err).To(BeNil())
			diff := cmp.Diff(b.String(), string(versionYaml))
			Expect(diff).To(Equal(""), diff)
		})
	})
})
