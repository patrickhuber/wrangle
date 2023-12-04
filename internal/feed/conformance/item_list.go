package conformance

import (
	"github.com/patrickhuber/wrangle/internal/feed"
	"github.com/patrickhuber/wrangle/internal/packages"
)

func createItem(name string, versions ...string) *feed.Item {

	return &feed.Item{
		Package: &packages.Package{
			Name:     name,
			Versions: createVersions(name, versions...),
		},
		State: &feed.State{
			LatestVersion: "1.0.0",
		},
		Template: "",
		Platforms: []*feed.Platform{
			{
				Name:          "linux",
				Architectures: []string{"amd64", "386"},
			},
		},
	}
}

func createVersions(name string, versions ...string) []*packages.Version {
	result := []*packages.Version{}
	for _, v := range versions {
		result = append(result, createVersion(name, v))
	}
	return result
}

func createVersion(name string, version string) *packages.Version {
	return &packages.Version{
		Version: version,
		Manifest: &packages.Manifest{
			Package: &packages.ManifestPackage{
				Name:    name,
				Version: version,
				Targets: []*packages.ManifestTarget{
					{
						Platform:     "linux",
						Architecture: "amd64",
						Steps: []*packages.ManifestStep{
							{
								Action: "download",
								With: map[string]any{
									"url": "https://www.google.com",
									"out": "test",
								},
							},
						},
					},
					{
						Platform:     "darwin",
						Architecture: "amd64",
						Steps: []*packages.ManifestStep{
							{
								Action: "download",
								With: map[string]any{
									"url": "https://www.google.com",
									"out": "test",
								},
							},
						},
					},
					{
						Platform:     "windows",
						Architecture: "amd64",
						Steps: []*packages.ManifestStep{
							{
								Action: "download",
								With: map[string]any{
									"url": "https://www.google.com",
									"out": "test",
								},
							},
						},
					},
				},
			},
		},
	}
}

func GetItemList() []*feed.Item {
	return []*feed.Item{
		createItem("test", "1.0.1", "1.0.0", "1.1.0"),
		createItem("ffa", "1.0.0"),
		createItem("tsa", "1.0.0"),
		createItem("other", "1.0.0", "2.0.0", "3.0.0"),
	}
}
