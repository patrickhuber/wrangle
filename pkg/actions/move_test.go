package actions_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/patrickhuber/go-log"
	filesystem "github.com/patrickhuber/go-xplat/fs"
	"github.com/patrickhuber/wrangle/pkg/actions"
)

var _ = Describe("Move", func() {
	It("can move file", func() {
		fs := filesystem.NewMemory()
		logger := log.Memory()
		provider := actions.NewMoveProvider(fs, logger)

		err := fs.WriteFile("/folder/file.txt", []byte("this is a test"), 0644)
		Expect(err).To(BeNil())

		t := &actions.Action{
			Type: "move",
			Parameters: map[string]interface{}{
				"source":      "file.txt",
				"destination": "moved.txt",
			},
		}
		ctx := &actions.Metadata{
			PackageVersionPath: "/folder",
		}
		err = provider.Execute(t, ctx)
		Expect(err).To(BeNil())
		ok, err := fs.Exists("/folder/moved.txt")
		Expect(err).To(BeNil())
		Expect(ok).To(BeTrue())
	})
})
