package githubrelease_test

import (
	"regexp"
	"testing"

	"github.com/patrickhuber/wrangle/internal/resource/githubrelease"
	"github.com/stretchr/testify/require"
)

func TestSource(t *testing.T) {
	source := githubrelease.Source{
		TagFilter: "[0-9]([.][0-9]){2}",
	}
	re := regexp.MustCompile(source.TagFilter)
	findString := re.FindString("v1.2.3")
	require.Equal(t, "1.2.3", findString)

	matchString := re.MatchString("v1.2.3")
	require.True(t, matchString)
}
