package env_test

import (
	"os"

	"github.com/patrickhuber/wrangle/collections"
	. "github.com/patrickhuber/wrangle/env"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Dictionary", func() {
	var (
		dictionary collections.Dictionary
	)
	key := "WRANGLE_ENV_TEST"
	value := "value"
	BeforeEach(func() {
		dictionary = NewDictionary()
	})
	Describe("Get", func() {
		Context("WhenVariableSet", func() {
			BeforeEach(func() {
				os.Setenv(key, value)
			})
			It("should find variable", func() {
				v, err := dictionary.Get(key)
				Expect(err).To(BeNil())
				Expect(v).To(Equal(value))
			})
			AfterEach(func() {
				os.Unsetenv(key)
			})
		})
		Context("WhenVariableNotSet", func() {
			It("should fail", func() {
				_, err := dictionary.Get(key)
				Expect(err).NotTo(BeNil())
			})
		})
	})
	Describe("Set", func() {
		Context("WhenVariableSet", func() {
			BeforeEach(func() {
				os.Setenv(key, value)
			})
			It("should overwrite value", func() {
				err := dictionary.Set(key, "other")
				Expect(err).To(BeNil())
				v, ok := os.LookupEnv(key)
				Expect(ok).To(BeTrue())
				Expect(v).To(Equal("other"))
			})
			AfterEach(func() {
				os.Unsetenv(key)
			})
		})
		Context("WhenVariableNotSet", func() {
			It("should set value", func() {
				err := dictionary.Set(key, value)
				Expect(err).To(BeNil())
				v, ok := os.LookupEnv(key)
				Expect(ok).To(BeTrue())
				Expect(v).To(Equal(value))
			})
		})
	})
	Describe("Unset", func() {
		Context("WhenVariableSet", func() {
			BeforeEach(func() {
				os.Setenv(key, value)
			})
			It("should unset variable", func() {
				err := dictionary.Unset(key)
				Expect(err).To(BeNil())
				_, ok := os.LookupEnv(key)
				Expect(ok).To(BeFalse())
			})
			AfterEach(func() {
				os.Unsetenv(key)
			})
		})
		Context("WhenVariableNotSet", func() {
			It("should fail", func() {
				_, err := dictionary.Get(key)
				Expect(err).NotTo(BeNil())
			})
		})
	})
	Describe("Lookup", func() {
		Context("WhenVariableSet", func() {
			BeforeEach(func() {
				os.Setenv(key, value)
			})
			It("should find variable", func() {
				v, ok := dictionary.Lookup(key)
				Expect(ok).To(BeTrue())
				Expect(v).To(Equal(value))
			})
			AfterEach(func() {
				os.Unsetenv(key)
			})
		})
		Context("WhenVariableNotSet", func() {
			It("should return false", func() {
				_ = dictionary.Unset(key)
				_, ok := dictionary.Lookup(key)
				Expect(ok).To(BeFalse())
			})
		})
	})
})
