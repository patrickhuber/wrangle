package commands_test

import (
	"bytes"
	"fmt"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/patrickhuber/wrangle/collections"
	"github.com/patrickhuber/wrangle/commands"
	"github.com/patrickhuber/wrangle/global"
	"github.com/patrickhuber/wrangle/ui"
)

var _ = Describe("Env", func() {
	var (
		console    ui.Console
		dictionary collections.Dictionary
		cmd        commands.Env
	)
	BeforeEach(func() {
		console = ui.NewMemoryConsole()
		dictionary = collections.NewDictionary()
		cmd = commands.NewEnv(console, dictionary)
	})
	Describe("NewEnv", func() {
		It("creates new env", func() {
			Expect(cmd).ToNot(BeNil())
		})
	})
	Describe("Execute", func() {
		const (
			packagesPath   = "/packages"
			configFilePath = "/config/config.yml"
		)
		Context("WhenAllEnvVarsSet", func() {
			BeforeEach(func() {
				dictionary.Set(global.PackagePathKey, packagesPath)
				dictionary.Set(global.ConfigFileKey, configFilePath)
			})
			It("should render both env vars", func() {
				err := cmd.Execute()
				Expect(err).To(BeNil())
				buffer := console.Out().(*bytes.Buffer)
				expected := fmt.Sprintf("%s=%s\n%s=%s\n", global.PackagePathKey, packagesPath, global.ConfigFileKey, configFilePath)
				Expect(buffer.String()).To(Equal(expected))
			})
		})
		Context("WhenOnlyPackagePathSet", func() {
			BeforeEach(func() {
				dictionary.Set(global.PackagePathKey, packagesPath)
			})
			It("should render package path and not config file", func() {
				err := cmd.Execute()
				Expect(err).To(BeNil())
				buffer := console.Out().(*bytes.Buffer)
				expected := fmt.Sprintf("%s=%s\n%s=\n", global.PackagePathKey, packagesPath, global.ConfigFileKey)
				Expect(buffer.String()).To(Equal(expected))
			})
		})
		Context("WhenOnlyConfigFileSet", func() {
			BeforeEach(func() {
				dictionary.Set(global.ConfigFileKey, configFilePath)
			})
			It("should render package path and not config file", func() {
				err := cmd.Execute()
				Expect(err).To(BeNil())
				buffer := console.Out().(*bytes.Buffer)
				expected := fmt.Sprintf("%s=\n%s=%s\n", global.PackagePathKey, global.ConfigFileKey, configFilePath)
				Expect(buffer.String()).To(Equal(expected))
			})
		})
		Context("WhenNoEnvVarsSet", func() {
			It("should render neither", func() {
				err := cmd.Execute()
				Expect(err).To(BeNil())
				buffer := console.Out().(*bytes.Buffer)
				expected := fmt.Sprintf("%s=\n%s=\n", global.PackagePathKey, global.ConfigFileKey)
				Expect(buffer.String()).To(Equal(expected))
			})
		})
	})
})
