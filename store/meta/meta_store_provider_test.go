package meta_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/patrickhuber/wrangle/config"
	"github.com/patrickhuber/wrangle/store"
	"github.com/patrickhuber/wrangle/store/meta"
)

var _ = Describe("MetaStoreProvider", func() {
	var (
		provider store.Provider
	)

	BeforeEach(func() {
		provider = meta.NewMetaStoreProvider("/this/is/the/config.yml")
	})

	Describe("Name", func() {
		It("should equal meta", func() {
			Expect(provider.Name()).To(Equal("meta"))
		})
	})

	Describe("Create", func() {
		var (
			s   store.Store
			err error
		)
		BeforeEach(func() {
			s, err = provider.Create(&config.Store{Name: "test"})
		})
		It("should assign source name", func() {
			Expect(err).To(BeNil())
			Expect(s).ToNot(BeNil())
			Expect(s.Name()).To(Equal("test"))
		})
	})
})
