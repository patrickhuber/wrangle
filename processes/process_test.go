package processes

import (
	"bytes"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Dispatch", func() {

	It("can run go version", func() {
		stdOut := bytes.Buffer{}
		stdErr := bytes.Buffer{}
		stdIn := bytes.Buffer{}
		command := NewProcess("go", []string{"version"}, make(map[string]string), &stdOut, &stdErr, &stdIn)
		err := command.Dispatch()
		Expect(err).To(BeNil())
	})

	It("writes to standard error", func() {
		stdOut := bytes.Buffer{}
		stdErr := bytes.Buffer{}
		stdIn := bytes.Buffer{}
		command := NewProcess("go", []string{}, make(map[string]string), &stdOut, &stdErr, &stdIn)
		err := command.Dispatch()
		Expect(err).ToNot(BeNil())
		Expect(stdErr.String()).ToNot(BeEmpty())
	})

	It("writes to standard output", func() {
		stdOut := bytes.Buffer{}
		stdErr := bytes.Buffer{}
		stdIn := bytes.Buffer{}
		command := NewProcess("go", []string{"version"}, make(map[string]string), &stdOut, &stdErr, &stdIn)
		err := command.Dispatch()
		Expect(err).To(BeNil())
		Expect(stdOut.String()).ToNot(BeEmpty())
		Expect(stdErr.String()).To(BeEmpty())
	})

	It("can run with nil environment variables", func() {
		stdOut := bytes.Buffer{}
		stdErr := bytes.Buffer{}
		stdIn := bytes.Buffer{}
		command := NewProcess("go", []string{"version"}, make(map[string]string), &stdOut, &stdErr, &stdIn)
		err := command.Dispatch()
		Expect(err).To(BeNil())
	})
})
