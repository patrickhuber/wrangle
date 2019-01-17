package templates_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/patrickhuber/wrangle/templates"
)

var _ = Describe("", func() {
	It("can evaluate string", func() {
		template := templates.NewTemplate("((key))")
		resolver, err := newSimpleResolver("/key", "value")
		Expect(err).To(BeNil())
		document, err := template.Evaluate(resolver)
		Expect(err).To(BeNil())
		Expect(document).To(Equal("value"))
	})

	It("can evaluate int", func() {
		template := templates.NewTemplate("((key))")
		resolver, err := newSimpleResolver("/key", 1)
		Expect(err).To(BeNil())
		document, err := template.Evaluate(resolver)
		Expect(err).To(BeNil())
		Expect(document).To(Equal(1))
	})

	It("can evaluate two keys in string", func() {
		template := templates.NewTemplate("((key)):((other))")
		resolver, err := newSimpleResolver("/key", "value", "/other", "thing")
		Expect(err).To(BeNil())
		document, err := template.Evaluate(resolver)
		Expect(err).To(BeNil())
		Expect(document).To(Equal("value:thing"))
	})

	It("can evaluate map string of string", func() {
		m := map[string]string{"key": "((key))"}
		template := templates.NewTemplate(m)
		resolver, err := newSimpleResolver("/key", "value")
		Expect(err).To(BeNil())
		document, err := template.Evaluate(resolver)
		Expect(err).To(BeNil())
		actual, ok := document.(map[string]interface{})
		Expect(ok).To(BeTrue())
		Expect(len(actual)).To(Equal(1))
		Expect(actual["key"]).To(Equal("value"))
	})

	It("can evaluate map string of interface", func() {
		m := map[string]interface{}{"key": "((key))"}
		template := templates.NewTemplate(m)
		resolver, err := newSimpleResolver("/key", "value")
		Expect(err).To(BeNil())
		document, err := template.Evaluate(resolver)
		Expect(err).To(BeNil())
		actual, ok := document.(map[string]interface{})
		Expect(ok).To(BeTrue())
		Expect(len(actual)).To(Equal(1))
		Expect(actual["key"]).To(Equal("value"))
	})

	It("can evaluate map interface of interface", func() {
		m := map[interface{}]interface{}{"key": "((key))"}
		template := templates.NewTemplate(m)
		resolver, err := newSimpleResolver("/key", "value")
		Expect(err).To(BeNil())
		document, err := template.Evaluate(resolver)
		Expect(err).To(BeNil())
		actual, ok := document.(map[interface{}]interface{})
		Expect(ok).To(BeTrue())
		Expect(len(actual)).To(Equal(1))
		Expect(actual["key"]).To(Equal("value"))
	})

	It("can evaluate nested map", func() {

		m := map[string]interface{}{"key": map[string]string{"nested": "((nested))"}}
		template := templates.NewTemplate(m)
		resolver, err := newSimpleResolver("/nested", "value")
		Expect(err).To(BeNil())
		document, err := template.Evaluate(resolver)
		Expect(err).To(BeNil())
		actual, ok := document.(map[string]interface{})
		Expect(ok).To(BeTrue())
		Expect(len(actual)).To(Equal(1))
		nested := actual["key"]
		nestedMap, ok := nested.(map[string]interface{})
		Expect(ok).To(BeTrue())
		Expect(nestedMap["nested"]).To(Equal("value"))
	})

	It("can evaluate string array", func() {

		a := []string{"one", "((key))", "three"}
		template := templates.NewTemplate(a)
		resolver, err := newSimpleResolver("/key", "value")
		Expect(err).To(BeNil())
		document, err := template.Evaluate(resolver)
		Expect(err).To(BeNil())
		actual, ok := document.([]interface{})
		Expect(ok).To(BeTrue())
		Expect(len(actual)).To(Equal(3))
		Expect(actual[1]).To(Equal("value"))
	})

	It("can evaluate interface array", func() {

		a := []interface{}{"one", "((key))", "three"}
		template := templates.NewTemplate(a)
		resolver, err := newSimpleResolver("/key", "value")
		Expect(err).To(BeNil())
		document, err := template.Evaluate(resolver)
		Expect(err).To(BeNil())
		actual, ok := document.([]interface{})
		Expect(ok).To(BeTrue())
		Expect(len(actual)).To(Equal(3))
		Expect(actual[1]).To(Equal("value"))
	})

	It("can patch in slice for string", func() {
		template := templates.NewTemplate("((key))")
		resolver, err := newSimpleResolver("/key", []string{"one", "two"})
		Expect(err).To(BeNil())
		document, err := template.Evaluate(resolver)
		Expect(err).To(BeNil())
		s, ok := document.([]string)
		Expect(ok).To(BeTrue())
		Expect(len(s)).To(Equal(2))
		Expect(s[0]).To(Equal("one"))
		Expect(s[1]).To(Equal("two"))
	})

	It("can patch in map for string", func() {
		template := templates.NewTemplate("((key))")
		resolver, err := newSimpleResolver("/key", map[string]string{"one": "test", "two": "other"})
		Expect(err).To(BeNil())
		document, err := template.Evaluate(resolver)
		Expect(err).To(BeNil())
		s, ok := document.(map[string]string)
		Expect(ok).To(BeTrue())
		Expect(len(s)).To(Equal(2))
		Expect(s["one"]).To(Equal("test"))
		Expect(s["two"]).To(Equal("other"))
	})

	It("can evaluate resolver pipeline", func() {
		template := templates.NewTemplate("((key1))")
		resolver1, err := newSimpleResolver("/key1", "((key2))")
		Expect(err).To(BeNil())
		resolver2, err := newSimpleResolver("/key2", "value")
		Expect(err).To(BeNil())
		document, err := template.Evaluate(resolver1, resolver2)
		Expect(err).To(BeNil())
		Expect(document).To(Equal("value"))
	})

})
