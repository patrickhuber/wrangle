package tasks_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/patrickhuber/wrangle/pkg/filesystem"
	"github.com/patrickhuber/wrangle/pkg/ilog"
	"github.com/patrickhuber/wrangle/pkg/tasks"
)

var _ = Describe("Move", func() {
	It("can move file", func() {
		fs := filesystem.NewMemory()
		logger := ilog.Memory()
		provider := tasks.NewMoveProvider(fs, logger)

		err := fs.Write("/folder/file.txt", []byte("this is a test"), 0644)
		Expect(err).To(BeNil())

		t := &tasks.Task{
			Type: "move",
			Parameters: map[string]interface{}{
				"source":      "file.txt",
				"destination": "moved.txt",
			},
		}
		ctx := &tasks.Metadata{
			PackageVersionPath: "/folder",
		}
		err = provider.Execute(t, ctx)
		Expect(err).To(BeNil())
		ok, err := fs.Exists("/folder/moved.txt")
		Expect(err).To(BeNil())
		Expect(ok).To(BeTrue())
	})
})
