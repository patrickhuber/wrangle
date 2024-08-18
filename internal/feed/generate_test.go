package feed_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/patrickhuber/wrangle/internal/feed"
)

func TestGenerate(t *testing.T) {
	request := &feed.GenerateRequest{
		Items: []*feed.GenerateItem{
			{
				Package: &feed.GeneratePackage{
					Name:     "test",
					Versions: []string{"1.0.0", "1.0.1"},
				},
				Platforms: []*feed.GeneratePlatform{
					{
						Name: "windows",
						Architectures: []string{
							"amd64",
							"arm64",
						},
					},
				},
			},
		},
	}
	response, err := feed.Generate(request)
	require.NoError(t, err)
	require.NotNil(t, response)
	require.Equal(t, 1, len(response.Packages))

	pkg := response.Packages[0]
	require.Equal(t, 2, len(pkg.Versions))
	for _, v := range pkg.Versions {
		require.NotEqual(t, "", v.Version)
	}
}
