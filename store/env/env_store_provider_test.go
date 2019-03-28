package env

import (
	"github.com/patrickhuber/wrangle/collections"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/patrickhuber/wrangle/config"
)

var _ = Describe("", func() {
	It("can crate env store provider", func() {
		variables := collections.NewDictionary()
		provider := NewEnvStoreProvider(variables)
		name := provider.Name()

		Expect(name).To(Equal("env"))
	})
	It("can create env store", func() {
		variables := collections.NewDictionary()
		provider := NewEnvStoreProvider(variables)

		prop1 := "test1"
		prop2 := "test2"
		env1 := "TEST_WRANGLE_ENV_VAL1"
		env2 := "TEST_WRANGLE_ENV_VAL2"
		val1 := "abc123"
		val2 := "zyx987"

		source := &config.Store{
			Name:      "env",
			StoreType: "env",
			Params: map[string]string{
				prop1: env1,
				prop2: env2,
			},
		}

		// create the store
		s, err := provider.Create(source)
		Expect(err).To(BeNil())
		Expect(s).ToNot(BeNil())

		// make sure the name is set
		Expect(s.Name()).To(Equal(source.Name))

		// set the env vars
		variables.Set(env1, val1)
		variables.Set(env2, val2)

		// verify
		d1, err := s.Get(prop1)
		Expect(err).To(BeNil())
		Expect(d1.Value()).To(Equal(val1))

		d2, err := s.Get(prop2)
		Expect(err).To(BeNil())
		Expect(d2.Value()).To(Equal(val2))
	})
})
