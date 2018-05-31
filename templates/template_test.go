package templates

import (
	"testing"

	"github.com/stretchr/testify/require"
)

type SimpleResolver struct {
	Map map[string]interface{}
}

func (resolver *SimpleResolver) Get(key string) interface{} {
	return resolver.Map[key]
}

func TestTemplate(t *testing.T) {

	t.Run("CanEvaluateString", func(t *testing.T) {
		r := require.New(t)
		template := NewTemplate("((key))")
		resolver := &SimpleResolver{Map: make(map[string]interface{})}
		resolver.Map["key"] = "value"
		document, err := template.Evaluate(resolver)
		r.Nil(err)
		r.Equal("value", document)
	})

	t.Run("CanEvaluateTwoKeysInString", func(t *testing.T) {
		r := require.New(t)
		template := NewTemplate("((key)):((other))")
		resolver := &SimpleResolver{
			Map: map[string]interface{}{"key": "value", "other": "thing"}}
		document, err := template.Evaluate(resolver)
		r.Nil(err)
		r.Equal("value:thing", document)
	})

	t.Run("CanEvaluateMapStringOfString", func(t *testing.T) {
		r := require.New(t)
		m := map[string]string{"key": "((key))"}
		template := NewTemplate(m)
		resolver := &SimpleResolver{
			Map: map[string]interface{}{"key": "value"},
		}
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
		resolver := &SimpleResolver{
			Map: map[string]interface{}{"key": "value"},
		}
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
		resolver := &SimpleResolver{
			Map: map[string]interface{}{"nested": "value"},
		}
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
		resolver := &SimpleResolver{
			Map: map[string]interface{}{"key": "value"},
		}
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
		resolver := &SimpleResolver{
			Map: map[string]interface{}{"key": "value"},
		}
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
		m := make(map[string]interface{})
		m["key"] = []string{"one", "two"}
		resolver := &SimpleResolver{Map: m}
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
		m := make(map[string]interface{})
		m["key"] = map[string]string{"one": "test", "two": "other"}
		resolver := &SimpleResolver{Map: m}
		document, err := template.Evaluate(resolver)
		r.Nil(err)
		s, ok := document.(map[string]string)
		r.True(ok)
		r.Equal(2, len(s))
		r.Equal("test", s["one"])
		r.Equal("other", s["two"])
	})
}
