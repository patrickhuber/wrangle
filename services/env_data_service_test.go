package services_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/patrickhuber/wrangle/collections"
	"github.com/patrickhuber/wrangle/global"
	"github.com/patrickhuber/wrangle/services"
)

var _ = Describe("EnvDataService", func() {
	var (
		dictionary collections.Dictionary
		cmd        services.EnvDataService
	)
	BeforeEach(func() {
		dictionary = collections.NewDictionary()
		cmd = services.NewEnvDataService(dictionary)
	})
	Describe("NewEnvDataService", func() {
		It("creates new env", func() {
			Expect(cmd).ToNot(BeNil())
		})
	})
	Describe("List", func() {
		const (
			rootPath       = "/root"
			binPath        = "/root/bin"
			packagesPath   = "/packages"
			configFilePath = "/config/config.yml"
		)
		Context("WhenAllEnvVarsSet", func() {
			BeforeEach(func() {
				dictionary.Set(global.PackagePathKey, packagesPath)
				dictionary.Set(global.ConfigFileKey, configFilePath)
				dictionary.Set(global.BinPathKey, binPath)
				dictionary.Set(global.RootPathKey, rootPath)
			})
			It("should render all env vars", func() {
				variables := cmd.List()
				Expect(len(variables)).To(Equal(4))
				Expect(variables[global.PackagePathKey]).To(Equal(packagesPath))
				Expect(variables[global.ConfigFileKey]).To(Equal(configFilePath))
				Expect(variables[global.BinPathKey]).To(Equal(binPath))
				Expect(variables[global.RootPathKey]).To(Equal(rootPath))
			})
		})
		Context("WhenOnlyPackagePathSet", func() {
			BeforeEach(func() {
				dictionary.Set(global.PackagePathKey, packagesPath)
			})
			It("should render package path and not config file", func() {
				variables := cmd.List()
				Expect(len(variables)).To(Equal(4))
				Expect(variables[global.PackagePathKey]).To(Equal(packagesPath))
				Expect(variables[global.ConfigFileKey]).To(Equal(""))
				Expect(variables[global.BinPathKey]).To(Equal(""))
				Expect(variables[global.RootPathKey]).To(Equal(""))
			})
		})
		Context("WhenOnlyConfigFileSet", func() {
			BeforeEach(func() {
				dictionary.Set(global.ConfigFileKey, configFilePath)
			})
			It("should render package path and not config file", func() {

				variables := cmd.List()
				Expect(len(variables)).To(Equal(4))

				Expect(variables[global.PackagePathKey]).To(Equal(""))
				Expect(variables[global.ConfigFileKey]).To(Equal(configFilePath))
				Expect(variables[global.BinPathKey]).To(Equal(""))
				Expect(variables[global.RootPathKey]).To(Equal(""))
			})
		})
		Context("WhenNoEnvVarsSet", func() {
			It("return all empty", func() {

				variables := cmd.List()
				Expect(len(variables)).To(Equal(4))

				Expect(variables[global.PackagePathKey]).To(Equal(""))
				Expect(variables[global.ConfigFileKey]).To(Equal(""))
				Expect(variables[global.BinPathKey]).To(Equal(""))
				Expect(variables[global.RootPathKey]).To(Equal(""))
			})
		})
	})
})
