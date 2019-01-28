package meta_test

import (
	"github.com/patrickhuber/wrangle/store"
	"github.com/patrickhuber/wrangle/store/meta"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("MetaStore", func() {
	var (
		metaStore store.Store
	)
	BeforeEach(func() {
		metaStore = meta.NewMetaStore("hello", "/some/path.ext")
	})

	Describe("Name", func() {
		It("returns meta", func() {
			Expect(metaStore.Name()).To(Equal("hello"))
		})
	})

	Describe("Get", func() {

		Context("WhenKeyIsConfigFilePath", func() {
			It("returns path to config file", func() {
				data, err := metaStore.Get(meta.ConfigFilePathKey)
				Expect(err).To(BeNil())
				Expect(data).ToNot(BeNil())
				Expect(data.Value()).To(Equal("/some/path.ext"))
			})
		})

		Context("WhenKeyIsConfigFileFolder", func() {
			It("returns path to config file directory", func() {
				data, err := metaStore.Get(meta.ConfigFileFolderKey)
				Expect(err).To(BeNil())
				Expect(data).ToNot(BeNil())
				Expect(data.Value()).To(Equal("/some"))
			})
		})

		Context("WhenKeyIsNotInStore", func() {
			It("returns error", func() {
				_, err := metaStore.Get("this is not a key")
				Expect(err).ToNot(BeNil())
			})
		})
	})

	Describe("Put", func() {
		It("is not implemented", func() {
			item := store.NewItem("test", store.Value, "value")
			err := metaStore.Set(item)
			Expect(err).ToNot(BeNil())
		})
	})
	Describe("Type", func() {
		It("returns meta", func() {
			value := metaStore.Type()
			Expect(value).To(Equal("meta"))
		})
	})
})
