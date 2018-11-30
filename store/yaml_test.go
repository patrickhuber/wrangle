package store_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
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

func TestCanParseYaml(test *testing.T) {
	require := require.New(test)
	var data = `
a: Easy!
b:
  c: 2
  d: [3, 4]
`
	t, err := parseT(data)
	if err != nil {
		test.Errorf(err.Error())
	}

	require.Equal("Easy!", t.A)
	require.Equal(2, t.B.RenamedC)
	require.Equal(2, len(t.B.D))
}

func parseT(data string) (*T, error) {
	t := &T{}
	err := yaml.Unmarshal([]byte(data), &t)
	if err != nil {
		return t, fmt.Errorf("error: %v", err)
	}
	return t, nil
}

func TestCanParseYamlOutOfOrder(test *testing.T) {
	require := require.New(test)
	var data = `
b:
  d: [1,2,3]
  c: 2  
a: test`
	t, err := parseT(data)
	require.Nil(err)
	require.Equal("test", t.A)
	require.Equal(2, t.B.RenamedC)
	require.Equal(3, len(t.B.D))
}

func TestCanParseYamlToMap(test *testing.T) {
	require := require.New(test)

	var data = `
b:
  d: [1,2,3]
  c: 2  
a: test`
	m := make(map[interface{}]interface{})
	err := yaml.Unmarshal([]byte(data), &m)
	require.Nil(err)

	a, ok := m["a"]
	require.Truef(ok, "unable to find key %s", "a")

	actualA, ok := a.(string)
	require.Truef(ok, "unable to cast to string")
	require.Equal("test", actualA)

	b, ok := m["b"]
	require.True(ok, "unable to find key b")

	bMap, ok := b.(map[interface{}]interface{})
	require.Truef(ok, "unable to cast b to type map[interface{}]interface{}")

	d, ok := bMap["d"]
	require.True(ok, "unable to find key b.d")

	dArray, ok := d.([]interface{})
	require.True(ok, "Unable to cast b.d to type []interface")
	require.Equal(3, len(dArray))
}

func TestCanUseVariable(test *testing.T) {
	require := require.New(test)

	var data = "b:\n  c: ((somevar))"
	m := make(map[interface{}]interface{})
	err := yaml.Unmarshal([]byte(data), &m)
	require.Nil(err)
}

func TestCanParseMultiLine(test *testing.T) {
	require := require.New(test)

	var data = "key: value\nid: value"
	m := make(map[interface{}]interface{})
	err := yaml.Unmarshal([]byte(data), &m)
	require.Nil(err)
}

func TestCanParseKeyValue(test *testing.T) {
	r := require.New(test)
	var data = "key: value"
	m := make(map[interface{}]interface{})
	err := yaml.Unmarshal([]byte(data), &m)
	r.Nil(err)
	r.Equal(1, len(m))
	r.Equal("value", m["key"])
}
