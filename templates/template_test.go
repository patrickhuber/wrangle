package templates

import (
	"testing"

	"github.com/stretchr/testify/require"
)

type SimpleResolver struct {
	Map map[string]string
}

func (resolver *SimpleResolver) Get(key string) interface{} {
	return resolver.Map[key]
}

func TestTemplate(t *testing.T) {
	t.Run("CanEvaluateString", func(t *testing.T) {
		r := require.New(t)
		template := NewTemplate("((key))")
		resolver := &SimpleResolver{Map: make(map[string]string)}
		resolver.Map["key"] = "value"
		document := template.Evaluate(resolver)
		r.Equal("value", document)
	})

	t.Run("CanEvaluateTwoKeysInString", func(t *testing.T) {
		r := require.New(t)
		template := NewTemplate("((key)):((other))")
		resolver := &SimpleResolver{
			Map: map[string]string{"key": "value", "other": "thing"}}
		document := template.Evaluate(resolver)
		r.Equal("value:thing", document)
	})

	t.Run("CanEvaluateMap", func(t *testing.T) {
		r := require.New(t)
		m := map[string]string{"key": "((key))"}
		template := NewTemplate(m)
		resolver := &SimpleResolver{
			Map: map[string]string{"key": "value"},
		}
		document := template.Evaluate(resolver)
		actual, ok := document.(map[string]string)
		r.True(ok)
		r.Equal(1, len(actual))
		r.Equal("value", actual["key"])
	})
}
