package env

import (
	"os"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("", func() {
	It("can read environment variable", func() {
		err := os.Setenv("TEST123", "abc123")
		Expect(err).To(BeNil())

		lookup := map[string]string{
			"somevalue": "TEST123",
		}

		store := NewEnvStore("", lookup)
		Expect(store).ToNot(BeNil())

		data, err := store.GetByName("somevalue")
		Expect(err).To(BeNil())
		Expect(data).ToNot(BeNil())
		Expect(data.Value()).To(Equal("abc123"))
	})

	It("can read environment variable with prefixed name", func() {
		err := os.Setenv("TEST123", "abc123")
		Expect(err).To(BeNil())

		lookup := map[string]string{
			"somevalue": "TEST123",
		}

		store := NewEnvStore("", lookup)
		Expect(store).ToNot(BeNil())

		data, err := store.GetByName("/somevalue")
		Expect(err).To(BeNil())
		Expect(data).ToNot(BeNil())
		Expect(data.Value()).To(Equal("abc123"))
	})

	Context("WhenEnvironmentVariableNotSet", func() {
		It("errors", func() {
			lookup := map[string]string{
				"somevalue": "TEST123",
			}
			store := NewEnvStore("", lookup)
			Expect(store).ToNot(BeNil())

			os.Unsetenv("TEST123")
			_, err := store.GetByName("somevalue")
			Expect(err).ToNot(BeNil())
		})
	})

	Describe("Name", func() {
		It("returns name", func() {
			store := NewEnvStore("env", nil)
			Expect(store).ToNot(BeNil())
			Expect(store.Name()).To(Equal("env"))
		})
	})

	Describe("Type", func() {
		it("returns type", func() {
			store := NewEnvStore("", nil)
			Expect(store).ToNot(BeNil())
			Expect(store.Type()).To(Equal("env"))
		})
	})
})
