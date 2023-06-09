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

var _ = Describe("Export", func() {

	DescribeTable("Execute", func(shell string, expected string) {
		env := env.NewMemory()
		env.Set("TEST", "TEST")
		console := console.NewMemory()
		shells := map[string]shellhook.Shell{
			shellhook.Bash:       shellhook.NewBash(),
			shellhook.Powershell: shellhook.NewPowershell(),
		}
		export := services.NewExport(env, shells, console)
		err := export.Execute(&services.ExportRequest{
			Shell: shell,
		})
		Expect(err).To(BeNil())

		outBuffer := console.Out().(*bytes.Buffer)
		result := outBuffer.String()
		Expect(result).To(Equal(expected))
	},
		Entry(shellhook.Bash, shellhook.Bash, "export TEST=TEST;\n"),
		Entry(shellhook.Powershell, shellhook.Powershell, "$env:TEST=\"TEST\";\n"))
})
