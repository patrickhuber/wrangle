package env

import (
	"os"

	"github.com/patrickhuber/wrangle/collections"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("", func() {
	var (
		variables collections.Dictionary
		lookup    map[string]string
	)
	BeforeEach(func() {
		variables = collections.NewDictionaryFromMap(map[string]string{})
		lookup = map[string]string{}
	})
	It("can read environment variable", func() {
		err := variables.Set("TEST123", "abc123")
		Expect(err).To(BeNil())

		lookup["somevalue"] = "TEST123"

		store := NewEnvStore("", lookup, variables)
		Expect(store).ToNot(BeNil())

		data, err := store.Get("somevalue")
		Expect(err).To(BeNil())
		Expect(data).ToNot(BeNil())
		Expect(data.Value()).To(Equal("abc123"))
	})

	It("can read environment variable with prefixed name", func() {
		err := variables.Set("TEST123", "abc123")
		Expect(err).To(BeNil())

		lookup["somevalue"] = "TEST123"

		store := NewEnvStore("", lookup, variables)
		Expect(store).ToNot(BeNil())

		data, err := store.Get("/somevalue")
		Expect(err).To(BeNil())
		Expect(data).ToNot(BeNil())
		Expect(data.Value()).To(Equal("abc123"))
	})

	Context("WhenEnvironmentVariableNotSet", func() {
		It("errors", func() {
			lookup["somevalue"] = "TEST123"
			store := NewEnvStore("", lookup, variables)
			Expect(store).ToNot(BeNil())

			os.Unsetenv("TEST123")
			_, err := store.Get("somevalue")
			Expect(err).ToNot(BeNil())
		})
	})

	Describe("Name", func() {
		It("returns name", func() {
			store := NewEnvStore("env", nil, nil)
			Expect(store).ToNot(BeNil())
			Expect(store.Name()).To(Equal("env"))
		})
	})

	Describe("Type", func() {
		It("returns type", func() {
			store := NewEnvStore("", nil, nil)
			Expect(store).ToNot(BeNil())
			Expect(store.Type()).To(Equal("env"))
		})
	})
})
