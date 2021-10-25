package githubrelease_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestGithubRelease(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "GithubRelease Suite")
}
