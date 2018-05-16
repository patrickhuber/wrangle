package store

import (
	"fmt"
	"testing"

	"gopkg.in/yaml.v2"
)

type T struct {
	A string
	B struct {
		RenamedC int   `yaml:"c"`
		D        []int `yaml:",flow"`
	}
}

func TestCanParseYaml(test *testing.T) {
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

	assertStringsAreEqual(test, t.A, "Easy!")
	assertIntsAreEqual(test, t.B.RenamedC, 2)
	assertIntsAreEqual(test, len(t.B.D), 2)
}

func parseT(data string) (*T, error) {
	t := &T{}
	err := yaml.Unmarshal([]byte(data), &t)
	if err != nil {
		return t, fmt.Errorf("error: %v", err)
	}
	return t, nil
}

func TestCanParseYamlOutOfOrer(test *testing.T) {
	var data = `
b:
  d: [1,2,3]
  c: 2  
a: test`
	t, err := parseT(data)
	if err != nil {
		test.Errorf(err.Error())
	}
	assertStringsAreEqual(test, t.A, "test")
	assertIntsAreEqual(test, t.B.RenamedC, 2)
	assertIntsAreEqual(test, len(t.B.D), 3)
}

func assertStringsAreEqual(t *testing.T, expected string, actual string) {
	if expected != actual {
		t.Errorf("expected %s, found %s", expected, actual)
	}
}

func assertIntsAreEqual(t *testing.T, expected int, actual int) {
	if expected != actual {
		t.Errorf("expected %d, found %d", expected, actual)
	}
}

func TestCanParseYamlToMap(test *testing.T) {
	var data = `
b:
  d: [1,2,3]
  c: 2  
a: test`
	m := make(map[interface{}]interface{})
	err := yaml.Unmarshal([]byte(data), &m)
	if err != nil {
		test.Errorf(err.Error())
	}
	a, ok := m["a"]
	if !ok {
		test.Fatalf("unable to find key %s", "a")
	}
	expectedA := "test"
	actualA, ok := a.(string)
	if actualA != expectedA {
		test.Fatalf("expected %s actual %s", expectedA, actualA)
	}
	b, ok := m["b"]
	if !ok {
		test.Fatalf("unable to find key %s", "b")
	}
	bMap, ok := b.(map[interface{}]interface{})
	if !ok {
		test.Fatalf("Unable to cast b to type map[interface{}]interface{}")
	}
	d, ok := bMap["d"]
	if !ok {
		test.Fatalf("Unable to find key b.d")
	}
	dArray, ok := d.([]interface{})
	if !ok {
		test.Fatalf("Unable to cast b.d to type []interface")
	}
	if len(dArray) != 3 {
		test.Fatalf("expected to find array of length 3, found %d", len(dArray))
	}
}

func TestCanUseVariable(test *testing.T) {
	var data = `b:\n  c: ((somevar))`
	m := make(map[interface{}]interface{})
	err := yaml.Unmarshal([]byte(data), &m)
	if err != nil {
		test.Errorf(err.Error())
	}
}
