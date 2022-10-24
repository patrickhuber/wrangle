package ilog_test

import (
	"bytes"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/patrickhuber/wrangle/pkg/ilog"
)

var _ = Describe("Memory", func() {
	var (
		out bytes.Buffer
	)
	Describe("Debug", func() {
		It("can log when debug set", func() {
			logger := ilog.MemoryWith(&out, ilog.SetLevel(ilog.DebugLevel))
			logger.Debug("test")
			Expect(out.Len()).ToNot(BeZero())
		})
		It("does not log when error set", func() {
			logger := ilog.MemoryWith(&out, ilog.SetLevel(ilog.ErrorLevel))
			logger.Debug("test")
			Expect(out.Len()).To(BeZero())
		})
	})
	Describe("Level", func() {
		It("can set", func() {
			logger := ilog.Memory(ilog.SetLevel(ilog.FatalLevel))
			Expect(logger.Level()).To(Equal(ilog.FatalLevel))
		})
	})
	AfterEach(func() {
		out.Reset()
	})
})
