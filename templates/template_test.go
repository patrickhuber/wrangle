package templates

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTemplate(t *testing.T) {

	t.Run("CanEvaluateString", func(t *testing.T) {
		r := require.New(t)
		template := NewTemplate("((key))")
		resolver, err := newSimpleResolver("key", "value")
		r.Nil(err)
		document, err := template.Evaluate(resolver)
		r.Nil(err)
		r.Equal("value", document)
	})

	t.Run("CanEvaluateInt", func(t *testing.T) {
		r := require.New(t)
		template := NewTemplate("((key))")
		resolver, err := newSimpleResolver("key", 1)
		r.Nil(err)
		document, err := template.Evaluate(resolver)
		r.Nil(err)
		r.Equal(1, document)
	})

	t.Run("CanEvaluateTwoKeysInString", func(t *testing.T) {
		r := require.New(t)
		template := NewTemplate("((key)):((other))")
		resolver, err := newSimpleResolver("key", "value", "other", "thing")
		r.Nil(err)
		document, err := template.Evaluate(resolver)
		r.Nil(err)
		r.Equal("value:thing", document)
	})

	t.Run("CanEvaluateMapStringOfString", func(t *testing.T) {
		r := require.New(t)
		m := map[string]string{"key": "((key))"}
		template := NewTemplate(m)
		resolver, err := newSimpleResolver("key", "value")
		r.Nil(err)
		document, err := template.Evaluate(resolver)
		r.Nil(err)
		actual, ok := document.(map[string]interface{})
		r.True(ok)
		r.Equal(1, len(actual))
		r.Equal("value", actual["key"])
	})

	t.Run("CanEvaluateMapStringOfInterface", func(t *testing.T) {
		r := require.New(t)
		m := map[string]interface{}{"key": "((key))"}
		template := NewTemplate(m)
		resolver, err := newSimpleResolver("key", "value")
		r.Nil(err)
		document, err := template.Evaluate(resolver)
		r.Nil(err)
		actual, ok := document.(map[string]interface{})
		r.True(ok)
		r.Equal(1, len(actual))
		r.Equal("value", actual["key"])
	})

	t.Run("CanEvaluateNestedMap", func(t *testing.T) {
		r := require.New(t)
		m := map[string]interface{}{"key": map[string]string{"nested": "((nested))"}}
		template := NewTemplate(m)
		resolver, err := newSimpleResolver("nested", "value")
		r.Nil(err)
		document, err := template.Evaluate(resolver)
		r.Nil(err)
		actual, ok := document.(map[string]interface{})
		r.True(ok)
		r.Equal(1, len(actual))
		nested := actual["key"]
		nestedMap, ok := nested.(map[string]interface{})
		r.True(ok)
		r.Equal("value", nestedMap["nested"])
	})

	t.Run("CanEvaluateStringArray", func(t *testing.T) {
		r := require.New(t)
		a := []string{"one", "((key))", "three"}
		template := NewTemplate(a)
		resolver, err := newSimpleResolver("key", "value")
		r.Nil(err)
		document, err := template.Evaluate(resolver)
		r.Nil(err)
		actual, ok := document.([]interface{})
		r.True(ok)
		r.Equal(3, len(actual))
		r.Equal("value", actual[1])
	})

	t.Run("CanEvaluateInterfaceArray", func(t *testing.T) {
		r := require.New(t)
		a := []interface{}{"one", "((key))", "three"}
		template := NewTemplate(a)
		resolver, err := newSimpleResolver("key", "value")
		r.Nil(err)
		document, err := template.Evaluate(resolver)
		r.Nil(err)
		actual, ok := document.([]interface{})
		r.True(ok)
		r.Equal(3, len(actual))
		r.Equal("value", actual[1])
	})

	t.Run("CanPatchInSliceForString", func(t *testing.T) {
		r := require.New(t)
		template := NewTemplate("((key))")
		resolver, err := newSimpleResolver("key", []string{"one", "two"})
		r.Nil(err)
		document, err := template.Evaluate(resolver)
		r.Nil(err)
		s, ok := document.([]string)
		r.True(ok)
		r.Equal(2, len(s))
		r.Equal("one", s[0])
		r.Equal("two", s[1])
	})

	t.Run("CanPatchInMapForString", func(t *testing.T) {
		r := require.New(t)
		template := NewTemplate("((key))")
		resolver, err := newSimpleResolver("key", map[string]string{"one": "test", "two": "other"})
		r.Nil(err)
		document, err := template.Evaluate(resolver)
		r.Nil(err)
		s, ok := document.(map[string]string)
		r.True(ok)
		r.Equal(2, len(s))
		r.Equal("test", s["one"])
		r.Equal("other", s["two"])
	})

	t.Run("CanEvaluateResolverPipeline", func(t *testing.T) {
		r := require.New(t)
		template := NewTemplate("((key1))")
		resolver1, err := newSimpleResolver("key1", "((key2))")
		r.Nil(err)
		resolver2, err := newSimpleResolver("key2", "value")
		r.Nil(err)
		document, err := template.Evaluate(resolver1, resolver2)
		r.Nil(err)
		r.Equal("value", document)
	})
}
