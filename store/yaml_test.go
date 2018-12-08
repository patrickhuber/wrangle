package store_test

import (
	"fmt"

	"gopkg.in/yaml.v2"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

type T struct {
	A string
	B struct {
		RenamedC int   `yaml:"c"`
		D        []int `yaml:",flow"`
	}
}

var _ = Describe("", func() {
	It("", func() {})
})

var _ = Describe("yaml", func() {
	It("can parse", func() {
		var data = `
a: Easy!
b:
  c: 2
  d: [3, 4]
`
		t, err := parseT(data)
		Expect(err).To(BeNil())
		Expect(t.A).To(Equal("Easy!"))
		Expect(t.B.RenamedC).To(Equal(2))
		Expect(len(t.B.D)).To(Equal(2))
	})

	It("can parse yaml out of order", func() {
		var data = `
b:
  d: [1,2,3]
  c: 2  
a: test`
		t, err := parseT(data)
		Expect(err).To(BeNil())
		Expect(t.A).To(Equal("test"))
		Expect(t.B.RenamedC).To(Equal(2))
		Expect(len(t.B.D)).To(Equal(3))
	})

	It("can parse yaml to map", func() {
		var data = `
b:
  d: [1,2,3]
  c: 2  
a: test`
		m := make(map[interface{}]interface{})
		err := yaml.Unmarshal([]byte(data), &m)
		Expect(err).To(BeNil())

		a, ok := m["a"]
		Expect(ok).To(BeTrue(), "unable to find key %s", "a")

		actualA, ok := a.(string)
		Expect(ok).To(BeTrue(), "unable to cast to string")
		Expect(actualA).To(Equal("test"))

		b, ok := m["b"]
		Expect(ok).To(BeTrue(), "unable to find key b")

		bMap, ok := b.(map[interface{}]interface{})
		Expect(ok).To(BeTrue(), "unable to cast b to type map[interface{}]interface{}")

		d, ok := bMap["d"]
		Expect(ok).To(BeTrue(), "unable to find key b.d")

		dArray, ok := d.([]interface{})
		Expect(ok).To(BeTrue(), "unable to cast b.d to type []interface")
		Expect(len(dArray)).To(Equal(3))
	})

	It("can use variable", func() {
		var data = "b:\n  c: ((somevar))"
		m := make(map[interface{}]interface{})
		err := yaml.Unmarshal([]byte(data), &m)
		Expect(err).To(BeNil())
	})

	It("can parse multiline", func() {
		var data = "key: value\nid: value"
		m := make(map[interface{}]interface{})
		err := yaml.Unmarshal([]byte(data), &m)
		Expect(err).To(BeNil())
	})

	It("can parse key value", func() {
		var data = "key: value"
		m := make(map[interface{}]interface{})
		err := yaml.Unmarshal([]byte(data), &m)
		Expect(err).To(BeNil())
		Expect(len(m)).To(Equal(1))
		Expect(m["key"]).To(Equal("value"))
	})
})

func parseT(data string) (*T, error) {
	t := &T{}
	err := yaml.Unmarshal([]byte(data), &t)
	if err != nil {
		return t, fmt.Errorf("error: %v", err)
	}
	return t, nil
}
