package collections_test

import (
	. "github.com/patrickhuber/wrangle/collections"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Dictionary", func() {
	var (
		dictionary Dictionary
	)
	BeforeEach(func() {
		dictionary = NewDictionary()
	})
	Describe("Lookup", func() {
		key := "key"
		value := "value"

		Context("WhenKeyIsPresent", func() {
			BeforeEach(func() {
				dictionary.Set(key, value)
			})
			It("should return true and the value", func() {
				v, ok := dictionary.Lookup(key)
				Expect(ok).To(BeTrue())
				Expect(v).To(Equal(value))
			})
		})
		Context("WhenKeyIsMissing", func() {
			It("should return false", func() {
				_, ok := dictionary.Lookup(key)
				Expect(ok).To(BeFalse())
			})
		})
	})
	Describe("Get", func() {
		key := "key"
		value := "value"

		Context("WhenKeyIsPresent", func() {
			BeforeEach(func() {
				dictionary.Set(key, value)
			})
			It("should return value", func() {
				v, err := dictionary.Get(key)
				Expect(err).To(BeNil())
				Expect(v).To(Equal(value))
			})
		})
		Context("WhenKeyIsMissing", func() {
			It("should throw an error", func() {
				_, err := dictionary.Get(key)
				Expect(err).ToNot(BeNil())
			})
		})
	})
	Describe("Set", func() {
		key := "key"
		value := "value"
		Context("WhenKeyIsPresent", func() {
			BeforeEach(func() {
				dictionary.Set(key, value)
			})
			It("should overwrite existing", func() {
				err := dictionary.Set(key, "other")
				Expect(err).To(BeNil())
				v, err := dictionary.Get(key)
				Expect(err).To(BeNil())
				Expect(v).To(Equal("other"))
			})
		})
		Context("WhenKeyIsMissing", func() {
			It("should set value", func() {
				err := dictionary.Set(key, value)
				Expect(err).To(BeNil())
				v, err := dictionary.Get(key)
				Expect(err).To(BeNil())
				Expect(v).To(Equal(value))
			})
		})
	})
	Describe("Unset", func() {
		key := "key"
		value := "value"
		Context("WhenKeyIsPresent", func() {
			BeforeEach(func() {
				dictionary.Set(key, value)
			})
			It("should remove key", func() {
				err := dictionary.Unset(key)
				Expect(err).To(BeNil())
			})
		})
		Context("WhenKeyIsNotPresent", func() {
			It("should error", func() {
				err := dictionary.Unset(key)
				Expect(err).NotTo(BeNil())
			})
		})
	})
})
