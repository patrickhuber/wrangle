package conformance

import (
	"github.com/patrickhuber/wrangle/pkg/feed"
	"github.com/patrickhuber/wrangle/pkg/packages"
)

func GetItemList() []*feed.Item {
	return []*feed.Item{
		{
			Package: &packages.Package{
				Name: "test",
				Versions: []*packages.Version{
					{
						Version: "1.0.1",
					},
					{
						Version: "1.0.0",
					},
					{
						Version: "1.1.0",
					},
				},
			},
			State: &feed.State{
				LatestVersion: "1.0.0",
			},
			Template: "",
			Platforms: []*feed.Platform{
				{
					Name:          "windows",
					Architectures: []string{"amd64", "386"},
				},
			},
		},
		{
			Package: &packages.Package{
				Name: "ffa",
				Versions: []*packages.Version{
					{
						Version: "1.0.0",
					},
				},
			},
			State: &feed.State{
				LatestVersion: "1.0.0",
			},
			Template: "",
			Platforms: []*feed.Platform{
				{
					Name:          "windows",
					Architectures: []string{"amd64", "386"},
				},
			},
		},
		{
			Package: &packages.Package{
				Name: "tsa",
				Versions: []*packages.Version{
					{
						Version: "1.0.0",
					},
				},
			},
			State: &feed.State{
				LatestVersion: "1.0.0",
			},
			Template: "",
			Platforms: []*feed.Platform{
				{
					Name:          "windows",
					Architectures: []string{"amd64", "386"},
				},
			},
		},
		{
			Package: &packages.Package{
				Name: "other",
				Versions: []*packages.Version{
					{
						Version: "1.0.0",
					},
					{
						Version: "2.0.0",
					},
					{
						Version: "3.0.0",
					},
				},
			},
			State: &feed.State{
				LatestVersion: "3.0.0",
			},
			Template: "",
			Platforms: []*feed.Platform{
				{
					Name:          "windows",
					Architectures: []string{"amd64", "386"},
				},
			},
		},
	}
}
