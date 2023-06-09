package services_test

import (
	"bytes"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/patrickhuber/go-shellhook"
	"github.com/patrickhuber/go-xplat/console"
	"github.com/patrickhuber/go-xplat/env"
	"github.com/patrickhuber/wrangle/internal/services"
)

var _ = Describe("Hook", func() {

	DescribeTable("Execute", func(shell string) {
		env := env.NewMemory()
		env.Set("TEST", "TEST")
		console := console.NewMemory()
		shells := map[string]shellhook.Shell{
			shellhook.Bash:       shellhook.NewBash(),
			shellhook.Powershell: shellhook.NewPowershell(),
		}
		export := services.NewHook(env, shells, console)
		err := export.Execute(&services.HookRequest{
			Executable: "/path/to/executable",
			Shell:      shell,
		})
		Expect(err).To(BeNil())

		outBuffer := console.Out().(*bytes.Buffer)
		result := outBuffer.String()
		Expect(result).ToNot(BeEmpty())
	},
		Entry(shellhook.Bash, shellhook.Bash),
		Entry(shellhook.Powershell, shellhook.Powershell))
})
