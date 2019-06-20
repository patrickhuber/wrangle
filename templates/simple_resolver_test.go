package templates_test

import (
	"fmt"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/patrickhuber/wrangle/templates"
)

type simpleResolver struct {
	innerMap map[string]interface{}
}

func newSimpleResolver(values ...interface{}) (templates.VariableResolver, error) {
	innerMap := make(map[string]interface{})
	if len(values)%2 == 1 {
		return nil, fmt.Errorf("values must be a list of key value pairs. ex: key1, value1, key2, value2")
	}
	for i := 0; i < len(values); i += 2 {
		key, ok := values[i].(string)
		if !ok {
			return nil, fmt.Errorf("key '%v' is not a string", key)
		}
		value := values[i+1]
		innerMap[key] = value
	}
	return &simpleResolver{innerMap: innerMap}, nil
}

func (resolver *simpleResolver) Get(key string) (interface{}, error) {
	value, ok := resolver.innerMap[key]
	if !ok {
		return nil, fmt.Errorf("unable to find key '%s'", key)
	}
	return value, nil
}

func (resolver *simpleResolver) Lookup(key string)(interface{}, bool, error){
	value, ok := resolver.innerMap[key]
	if !ok{
		return nil, ok, nil
	}
	return value, true, nil
}

var _ = Describe("SimpleResolver", func() {
	It("can create", func() {

		resolver, err := newSimpleResolver("key", "value", "key1", "value1")
		Expect(err).To(BeNil())
		Expect(resolver).ToNot(BeNil())

		value, err := resolver.Get("key")
		Expect(err).To(BeNil())
		Expect(value).ToNot(BeNil())
		Expect(value).To(Equal("value"))

		value, err = resolver.Get("key1")
		Expect(err).To(BeNil())
		Expect(value).ToNot(BeNil())
		Expect(value).To(Equal("value1"))
	})
})
